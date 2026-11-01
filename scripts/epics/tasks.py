#!/usr/bin/env python3
"""Task state for working through docs/EPICs/ROADMAP.md in order.

Usage (from the repo root):
  tasks.py next                 print the next task to work on (JSON); resumes a task already "in progress"
  tasks.py show <ID>            print one task (JSON), with its status and whether its dependencies are done
  tasks.py start <ID>           mark it "in progress" and write the marker the Stop hook reads
  tasks.py done <ID>            mark it "done" (the Stop hook clears the marker after lint and tests pass)
  tasks.py block <ID> <reason>  mark it "blocked", add a dated note to the task file, clear the marker
  tasks.py set <ID> <status>    set any status by hand (todo, in progress, done, blocked)
  tasks.py status               counts per milestone
  tasks.py changed              services touched since the current task started (for the Stop hook)

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
    left = [x for x in order if st[x] != "done"]
    if not left:
        print(json.dumps({"next": None, "message": "NO TASK LEFT"}))
    else:
        blocked = [x for x in order if st[x] == "blocked"]
        print(json.dumps({"next": None, "message": "NOTHING READY", "blocked": blocked, "remaining": len(left)}))


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
        set_status(t, "in progress")
        MARKER.parent.mkdir(parents=True, exist_ok=True)
        MARKER.write_text(json.dumps({"id": t["id"], "started_sha": git("rev-parse", "HEAD"),
                                      "started_at": datetime.datetime.now().isoformat(timespec="seconds"),
                                      "attempts": 0}, indent=1))
        print(f"started {t['id']} on branch {branch_for(t)}")
    elif cmd == "done":
        set_status(t, "done")
        print(f"done {t['id']}")
    elif cmd == "block":
        reason = " ".join(args[1:]).strip()
        if not reason:
            sys.exit("block needs a reason")
        set_status(t, "blocked", note=f"blocked: {reason}")
        MARKER.unlink(missing_ok=True)
        print(f"blocked {t['id']}: {reason}")
    elif cmd == "set":
        set_status(t, " ".join(args[1:]))
        print(f"{t['id']}: {' '.join(args[1:])}")
    else:
        sys.exit(__doc__)


if __name__ == "__main__":
    main()
