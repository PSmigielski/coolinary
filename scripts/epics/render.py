#!/usr/bin/env python3
"""Render docs/EPICs/ROADMAP.md from the JSON written by order.py. Usage: render.py <order.json> <docs/EPICs>."""
import json
import sys
from pathlib import Path

data = json.load(open(sys.argv[1]))
root = Path(sys.argv[2])
T, D, M = data["tasks"], data["deps"], data["M"]
DAYS = {"S": 1.0, "M": 2.5, "L": 4.5}
SVC = {"GW": "gateway", "IDN": "identity", "NTF": "notifier", "ADM": "admin", "KIT": "kitchen",
       "PLN": "planner", "OFF": "offers", "AI": "ai", "MOB": "mobile"}
REL = {"R1": "Release 1 — Launch", "R2": "Release 2 — Flyers and family", "R3": "Release 3 — Depth"}

text = {k: (root / v["path"]).read_text() for k, v in T.items()}
flags = {}
for k in T:
    f = []
    if "DECISION NEEDED" in text[k]:
        f.append("decision")
    if "GATE:" in text[k]:
        f.append("gate")
    flags[k] = ", ".join(f)

num = {}
n = 0
for seq in data["order"]:
    for t in seq:
        n += 1
        num[t] = n

lines = []
w = lines.append
w("# Coolinary — roadmap")
w("")
w("The order in which to build every task in this folder: all three releases, every service, the mobile")
w("communication layer included. Work top to bottom. Each task comes after every task it depends on, so the")
w("numbers are a valid build order for one developer. Tasks inside a milestone that don't depend on each other can")
w("also be done in parallel.")
w("")
w("- **Milestones** end with something you can try in the app or the back office (\"Done when\"). Release 1 ends at")
w("  M11, release 2 at M17, release 3 at M21.")
w("- **Estimate** uses S = 1, M = 2.5, L = 4.5 developer-days. Weeks assume 5 working days and one developer.")
w("- **After** lists the task's dependencies (task numbers in this file).")
w("- **Status** comes from each task file (`todo`, `in progress`, `done`, `blocked`). Work through it with")
w("  `/next-task` in Claude Code or `scripts/epics/run-tasks.sh` (see README.md, \"Working through the tasks\").")
w("- **Flags:** `decision` = the task still has an open `DECISION NEEDED:` (none since the 1 November 2026 decisions,")
w("  see [DECISIONS.md](DECISIONS.md));")
w("  `gate` = it must not go live until an outside condition is met (licence, retailer terms, contract).")
w("- Built from the \"Depends on\" field of every task file. When a dependency names another service's epic, it points")
w("  to the task in that epic that defines the RPC, event, route or table it names. Rebuild this file whenever")
w("  dependencies change.")
w("")

# Summary
w("## Milestones")
w("")
w("| Milestone | Release | Done when | Tasks | Done | Days | Weeks (cumulative) |")
w("|---|---|---|---|---|---|---|")
cum = 0.0
for i, seq in enumerate(data["order"]):
    mid, rel, title, done = M[i]
    days = sum(DAYS[T[t]["size"]] for t in seq)
    cum += days
    first, last = (num[seq[0]], num[seq[-1]]) if seq else ("", "")
    ndone = sum(1 for t in seq if T[t].get("status") == "done")
    w(f"| [{mid} {title}](#{mid.lower()}-{title.lower().replace(' ', '-').replace(',', '').replace('/', '')}) | {rel} | {done} | {len(seq)} (#{first}–{last}) | {ndone}/{len(seq)} | {days:g} | {cum / 5:.0f} |")
w("")
tot = sum(DAYS[T[t]["size"]] for t in T)
by_rel = {}
for t in T:
    by_rel[T[t]["rel"]] = by_rel.get(T[t]["rel"], 0) + DAYS[T[t]["size"]]
w(f"Total: {len(T)} tasks, about {tot:g} developer-days ({tot / 5:.0f} weeks for one developer): "
  + ", ".join(f"{r} {by_rel[r]:g} days" for r in sorted(by_rel)) + ". Frontend screens are not included yet.")
w("")

cur_rel = None
for i, seq in enumerate(data["order"]):
    mid, rel, title, done = M[i]
    if rel != cur_rel:
        cur_rel = rel
        w(f"## {REL[rel]}")
        w("")
    days = sum(DAYS[T[t]["size"]] for t in seq)
    w(f"### {mid} {title}")
    w("")
    w(f"**Done when:** {done}  ")
    w(f"{len(seq)} tasks, about {days:g} days.")
    w("")
    if not seq:
        w("No tasks.")
        w("")
        continue
    w("| # | Task | Title | Service | Size | After | Flags | Status |")
    w("|---|---|---|---|---|---|---|---|")
    for t in seq:
        after = sorted((num[d] for d in D.get(t, []) if d in num), reverse=True)
        a = ", ".join(str(x) for x in sorted(after[:6]))
        if len(after) > 6:
            a += f" +{len(after) - 6}"
        title_ = T[t]["title"].replace("|", "\\|")
        w(f"| {num[t]} | [{t}]({T[t]['path']}) | {title_} | {SVC[T[t]['pre']]} | {T[t]['size']} | {a} | {flags[t]} | {T[t].get('status', 'todo')} |")
    w("")

w("## Later releases that extend earlier tasks")
w("")
w("These earlier tasks name work from a later release (usually a table or hook added later). Build the earlier task")
w("in its own release; when you build the later task, also extend the earlier one as described.")
w("")
w("| Earlier task (#) | Release | Extended by (#) | Release |")
w("|---|---|---|---|")
for a, b in data["forward"]:
    w(f"| [{a}]({T[a]['path']}) ({num[a]}) | {T[a]['rel']} | [{b}]({T[b]['path']}) ({num[b]}) | {T[b]['rel']} |")
w("")
if data["cycles"]:
    w("## Dependency loops resolved")
    w("")
    w("Two task files depend on each other in a loop. The roadmap breaks each loop at the link shown last, which is")
    w("safe to build in either order; fix the \"Depends on\" fields when you start these tasks.")
    w("")
    for c in data["cycles"]:
        w(f"- {c}")
    w("")

(root / "ROADMAP.md").write_text("\n".join(lines) + "\n")
print("written", len(lines), "lines;", n, "tasks numbered")
