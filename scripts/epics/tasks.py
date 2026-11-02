#!/usr/bin/env python3
"""Task state for working through docs/EPICs/ROADMAP.md in order.

Usage (from the repo root):
  tasks.py next                 print the next task to work on (JSON); resumes a task already "in progress"
  tasks.py ready [N]            print a batch of up to N tasks (default 4) to build in parallel (JSON): tasks
                                already "in progress" first, then ready ones, all from one milestone
  tasks.py show <ID>            print one task (JSON), with its status and whether its dependencies are done
  tasks.py start <ID>...        mark them "in progress" and add them to the marker the Stop hook reads
  tasks.py done <ID>            mark it "done" (the Stop hook clears the marker after lint and tests pass)
  tasks.py block <ID> <reason>  mark it "blocked", add a dated note to the task file, take it off the marker
  tasks.py set <ID> <status>    set any status by hand (todo, in progress, done, blocked); todo and blocked
                                take it off the marker
  tasks.py status               counts per milestone
  tasks.py changed              services touched since the current task or batch started (for the Stop hook)

Order and dependencies come from ROADMAP.md (its "After" column); status comes from each task file's
"| Status | ... |" row. A task is ready when it isn't done or blocked and every task it comes after is done.
"""
import datetime
import json
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
EPICS = ROOT / "docs" / "EPICs"
MARKER = ROOT / ".claude" / "state" / "current-task.json"
STATUSES = ("todo", "in progress", "done", "blocked")
LINTABLE = ("gateway", "identity", "notifier", "kitchen", "planner", "offers", "ai", "mobile", "contracts", "admin")
ROW = re.compile(r"^\| (\d+) \| \[([A-Z]+-\d\d-\d\d)\]\(([^)]+)\) \| (.*?) \| (\w+) \| ([SML]) \| ([^|]*) \|")
MS = re.compile(r"^### (M\d+) (.+)$")


def roadmap():
    tasks, order, cur = {}, [], None
    for line in (EPICS / "ROADMAP.md").read_text().splitlines():
        if m := MS.match(line):
            cur = (m.group(1), m.group(2).strip())
            continue
        if (m := ROW.match(line)) and cur:
            num, tid, path, title, svc, size, after = m.groups()
            nums = [int(x) for x in re.findall(r"\d+", after.split("+")[0])]
            tasks[tid] = dict(num=int(num), id=tid, path=path, title=title.replace("\\|", "|"), service=svc,
                              size=size, after=nums, milestone=cur[0], milestone_title=cur[1])
            order.append(tid)
    by_num = {t["num"]: t["id"] for t in tasks.values()}
    for t in tasks.values():
        t["deps"] = [by_num[n] for n in t["after"] if n in by_num]
    return tasks, order


def status_of(t):
    m = re.search(r"^\| Status \| ([^|]+?) \|$", (EPICS / t["path"]).read_text(), re.M)
    return m.group(1).strip() if m else "todo"


def set_status(t, status, note=None):
    if status not in STATUSES:
        sys.exit(f"unknown status {status!r}; use one of {', '.join(STATUSES)}")
    f = EPICS / t["path"]
    text = f.read_text()
    text, n = re.subn(r"^\| Status \| [^|]+? \|$", f"| Status | {status} |", text, count=1, flags=re.M)
    if not n:
        sys.exit(f"{t['id']}: no Status row in {f}")
    if note:
        stamp = datetime.date.today().isoformat()
        if "\n## Blocked\n" not in text:
            text = text.rstrip("\n") + "\n\n## Blocked\n"
        text = text.rstrip("\n") + f"\n\n- {stamp}: {note}\n"
    f.write_text(text)


def branch_for(t):
    slug = re.sub(r"[^a-z0-9]+", "-", t["milestone_title"].lower()).strip("-")
    return f"m{int(t['milestone'][1:]):02d}-{slug}"


def git(*args):
    return subprocess.run(["git", *args], cwd=ROOT, capture_output=True, text=True).stdout.strip()


def info(t, tasks, order, status):
    same = [x for x in order if tasks[x]["milestone"] == t["milestone"]]
    rest = [x for x in same if x != t["id"] and status_of(tasks[x]) != "done"]
    return {**{k: t[k] for k in ("id", "num", "title", "path", "service", "size", "milestone", "milestone_title")},
            "file": f"docs/EPICs/{t['path']}", "status": status, "branch": branch_for(t),
            "deps": t["deps"], "last_in_milestone": not rest}


def none_left(order, st):
    left = [x for x in order if st[x] != "done"]
    if not left:
        return {"next": None, "message": "NO TASK LEFT"}
    blocked = [x for x in order if st[x] == "blocked"]
    return {"next": None, "message": "NOTHING READY", "blocked": blocked, "remaining": len(left)}


def cmd_ready(tasks, order, limit):
    st = {tid: status_of(tasks[tid]) for tid in order}
    started = [x for x in order if st[x] == "in progress"]
    ready = [x for x in order if st[x] == "todo" and all(st[d] == "done" for d in tasks[x]["deps"])]
    pool = started + ready
    if not pool:
        print(json.dumps(none_left(order, st)))
        return
    first = tasks[pool[0]]
    batch = [x for x in pool if tasks[x]["milestone"] == first["milestone"]][:limit]
    left = [x for x in order if tasks[x]["milestone"] == first["milestone"] and st[x] != "done" and x not in batch]
    print(json.dumps({"milestone": first["milestone"], "milestone_title": first["milestone_title"],
                      "branch": branch_for(first), "last_in_milestone": not left,
                      "tasks": [info(tasks[x], tasks, order, st[x]) for x in batch]}, indent=1))


def cmd_next(tasks, order):
    st = {tid: status_of(tasks[tid]) for tid in order}
    for tid in order:
        if st[tid] == "in progress":
            print(json.dumps(info(tasks[tid], tasks, order, "in progress"), indent=1))
            return
    for tid in order:
        if st[tid] in ("done", "blocked"):
            continue
        if all(st[d] == "done" for d in tasks[tid]["deps"]):
            print(json.dumps(info(tasks[tid], tasks, order, st[tid]), indent=1))
            return
    print(json.dumps(none_left(order, st)))


def find(tasks, tid):
    if tid not in tasks:
        sys.exit(f"unknown task {tid}")
    return tasks[tid]


def cmd_status(tasks, order):
    rows = {}
    for tid in order:
        t = tasks[tid]
        r = rows.setdefault(t["milestone"], {"title": t["milestone_title"], **{s: 0 for s in STATUSES}})
        r[status_of(t)] += 1
    for mid, r in rows.items():
        total = sum(r[s] for s in STATUSES)
        print(f"{mid:<4} {r['title']:<45} done {r['done']:>3}/{total:<3} in progress {r['in progress']} blocked {r['blocked']}")


def load_marker():
    if not MARKER.exists():
        return None
    m = json.loads(MARKER.read_text())
    if "ids" not in m:  # written before batches: one "id"
        m["ids"] = [m.pop("id")] if m.get("id") else []
    return m


def save_marker(m):
    MARKER.parent.mkdir(parents=True, exist_ok=True)
    MARKER.write_text(json.dumps(m, indent=1))


def drop_from_marker(tid):
    """The task is no longer being built: take it off the marker, and remove the marker when nothing is left."""
    m = load_marker()
    if m is None or tid not in m["ids"]:
        return
    m["ids"].remove(tid)
    if m["ids"]:
        save_marker(m)
    else:
        MARKER.unlink(missing_ok=True)


def cmd_changed():
    if not MARKER.exists():
        return
    start = json.loads(MARKER.read_text()).get("started_sha", "")
    names = set()
    if start:
        names |= set(git("diff", "--name-only", start).splitlines())
    names |= set(git("diff", "--name-only").splitlines())
    names |= set(git("ls-files", "--others", "--exclude-standard").splitlines())
    print(" ".join(sorted({n.split("/")[0] for n in names if n.split("/")[0] in LINTABLE})))


def main():
    if len(sys.argv) < 2:
        sys.exit(__doc__)
    cmd, args = sys.argv[1], sys.argv[2:]
    if cmd == "changed":
        return cmd_changed()
    tasks, order = roadmap()
    if cmd == "next":
        return cmd_next(tasks, order)
    if cmd == "status":
        return cmd_status(tasks, order)
    if cmd == "ready":
        if args and not args[0].isdigit():
            sys.exit("ready takes a number: tasks.py ready [N]")
        return cmd_ready(tasks, order, max(1, int(args[0])) if args else 4)
    if not args:
        sys.exit(__doc__)
    t = find(tasks, args[0])
    if cmd == "show":
        st = status_of(t)
        out = info(t, tasks, order, st)
        out["ready"] = all(status_of(tasks[d]) == "done" for d in t["deps"])
        print(json.dumps(out, indent=1))
        return
    if cmd == "start":
        batch = [find(tasks, x) for x in args]
        if len({x["milestone"] for x in batch}) > 1:
            sys.exit("start: all tasks of a batch must be in one milestone (one branch)")
        for x in batch:
            set_status(x, "in progress")
        m = load_marker() or {"ids": [], "started_sha": git("rev-parse", "HEAD"),
                              "started_at": datetime.datetime.now().isoformat(timespec="seconds"), "attempts": 0}
        m["ids"] += [x["id"] for x in batch if x["id"] not in m["ids"]]
        save_marker(m)
        print(f"started {' '.join(x['id'] for x in batch)} on branch {branch_for(t)}")
    elif cmd == "done":
        set_status(t, "done")
        print(f"done {t['id']}")
    elif cmd == "block":
        reason = " ".join(args[1:]).strip()
        if not reason:
            sys.exit("block needs a reason")
        set_status(t, "blocked", note=f"blocked: {reason}")
        drop_from_marker(t["id"])
        print(f"blocked {t['id']}: {reason}")
    elif cmd == "set":
        status = " ".join(args[1:])
        set_status(t, status)
        if status in ("todo", "blocked"):
            drop_from_marker(t["id"])
        print(f"{t['id']}: {status}")
    else:
        sys.exit(__doc__)


if __name__ == "__main__":
    main()
