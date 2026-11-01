#!/usr/bin/env python3
"""Rewrite docs/EPICs/{kitchen,planner}/README.md (generated). Usage: readmes.py <order.json> <docs/EPICs>."""
import json
import re
import sys
from pathlib import Path

data = json.load(open(sys.argv[1]))
root = Path(sys.argv[2])
T, M = data["tasks"], data["M"]
num, n = {}, 0
for seq in data["order"]:
    for t in seq:
        n += 1
        num[t] = n
ms = data["ms"]
ADMIN_RPCS = (root / "admin" / "README.md").read_text()

INTRO = {
    "kitchen": (
        "Go (`kitchen/`). Owns the ingredient catalogue and its translations, name matching (pg_trgm + Haiku), the "
        "barcode lookup chain, unknown products, nutrition tables and the nutrition calculation, the household's fridge "
        "stock and AI intake, expiry reminders, the kitchen profile (decision 3) with per-person targets, and the savings "
        "ledger. It produces `events.kitchen` and publishes `stock.changed` live notices.",
        "What exists today: `KitchenService.ResolveBarcode` over a `products` table with one example row (pg_trgm "
        "installed), the Go platform packages (config, db, cache, log, rpc, admin) and the Go → Bun smoke test.",
        "KIT-01 to KIT-06 and KIT-07 to KIT-13 were written by two writers; dependencies between the halves use epic ids.",
    ),
    "planner": (
        "Go (`planner/`). Owns the recipe library (decision 2): dishes, import and AI drafting, review, calories per "
        "dish, browse and \"what can I make now\", the rule checker on every AI dish, day, family and week plans, swaps "
        "and rewrites, cooking mode and \"cooked it\", per-person portions, the shopping list with live notices, base "
        "pantries, offer badges and reminders, plan from deals, the trip price matrix and the cheapest top-up. It "
        "consumes events and publishes none.",
        "What exists today: `PlannerService.SuggestDishes` over a skeleton `dishes` table with five example rows, "
        "replaced by PLN-02.",
        "PLN-01 to PLN-11 and PLN-12 to PLN-21 were written by two writers; dependencies between the halves use epic ids.",
    ),
}

for svc, pre in (("kitchen", "KIT"), ("planner", "PLN")):
    rp = lambda t: T[t]['path'].removeprefix(svc + '/')
    sdir = root / svc
    epics = sorted(d for d in sdir.iterdir() if d.is_dir() and re.match(r"E\d\d-", d.name))
    mine = sorted(t for t in T if T[t]["pre"] == pre)
    lines = []
    w = lines.append
    a, b, c = INTRO[svc]
    w(f"# {svc} — epics and tasks")
    w("")
    w(a)
    w("")
    w(b)
    w("")
    w(c)
    w("")
    w("## Epics")
    w("")
    w("| Epic | Title | Rel. | Tasks |")
    w("|---|---|---|---|")
    for d in epics:
        ep = f"{pre}-{d.name[1:3]}"
        head = (d / "README.md").read_text().splitlines()[0] if (d / "README.md").exists() else ep
        title = re.sub(rf"^#\s*{ep}\s*[—–-]?\s*", "", head).strip()
        ts = [t for t in mine if t.startswith(ep + "-")]
        rels = sorted({T[t]["rel"] for t in ts})
        rel = rels[0] if len(rels) == 1 else f"{rels[0]}–{rels[-1]}"
        w(f"| [{ep}]({d.name}/README.md) | {title} | {rel} | {len(ts)} |")
    w("")
    counts = {r: sum(1 for t in mine if T[t]["rel"] == r) for r in ("R1", "R2", "R3")}
    w(f"{len(mine)} tasks: " + ", ".join(f"{v} in {k}" for k, v in counts.items() if v) + ".")
    w("")

    # RPCs defined
    w("## What other services use")
    w("")
    w("RPCs defined in this service's tasks (first task that defines each). Admin RPCs are called only by admin.")
    w("")
    w("| RPC | Defined in | Admin |")
    w("|---|---|---|")
    seen = {}
    for t in mine:
        txt = (root / T[t]["path"]).read_text()
        for m in re.finditer(r"```proto\n(.*?)```", txt, re.S):
            before = txt[max(0, m.start() - 400):m.start()]
            block = None
            for line in m.group(1).splitlines():
                sm = re.match(r"^\s*service\s+(\w+)", line)
                if sm:
                    block = sm.group(1)
                rm = re.match(r"^\s*rpc\s+(\w+)\s*\(", line)
                if rm and rm.group(1) not in seen:
                    admin = (block and "Admin" in block) or (block is None and "AdminService" in before)
                    seen[rm.group(1)] = (t, "yes" if (admin or re.search(rf"`{rm.group(1)}`|\\b{rm.group(1)}\\b", ADMIN_RPCS)) else "")
    for name, (t, admin) in sorted(seen.items(), key=lambda x: x[1][0]):
        w(f"| `{name}` | [{t}]({rp(t)}) | {admin} |")
    w("")

    # Build order
    w("## Build order")
    w("")
    w("Task order follows [ROADMAP.md](../ROADMAP.md); `#` is the task's number there and `Milestone` its milestone.")
    w("")
    for rel in ("R1", "R2", "R3"):
        ts = sorted((t for t in mine if T[t]["rel"] == rel), key=lambda t: num[t])
        if not ts:
            continue
        w(f"### Release {rel[1]}")
        w("")
        w("| # | Task | Title | Size | Milestone |")
        w("|---|---|---|---|---|")
        for t in ts:
            title = T[t]["title"].replace("|", "\\|")
            w(f"| {num[t]} | [{t}]({rp(t)}) | {title} | {T[t]['size']} | {ms[t]} |")
        w("")

    def collect(marker):
        out = []
        for t in mine:
            lines_ = (root / T[t]["path"]).read_text().splitlines()
            i = 0
            while i < len(lines_):
                if marker in lines_[i]:
                    chunk = [lines_[i][lines_[i].index(marker):]]
                    j = i + 1
                    while j < len(lines_) and lines_[j].strip() and not re.match(r"^\s*(\d+\.|[-*]|\||#)", lines_[j]):
                        chunk.append(lines_[j].strip())
                        j += 1
                    s_ = re.sub(r"^:?`?\s+", "", " ".join(chunk).replace(marker, "", 1)).strip(" :")
                    if len(s_) > 260:
                        s_ = s_[:257].rsplit(" ", 1)[0] + " …"
                    out.append((t, s_.replace("|", "\\|")))
                    i = j
                else:
                    i += 1
        return out

    gates = collect("GATE:")
    if gates:
        w("## Gates")
        w("")
        w("| Task | Gate |")
        w("|---|---|")
        for t, s in gates:
            w(f"| [{t}]({rp(t)}) | {s} |")
        w("")
    decided = collect("Decided (Q")
    if decided:
        w("## Decisions")
        w("")
        w("Settled answers (full list in [DECISIONS.md](../DECISIONS.md)); each task states its rule as `Decided (Qnnn)`.")
        w("")
        w("| Task | Decision |")
        w("|---|---|")
        for t, s_ in decided:
            label = re.sub(r"^(\d{3})\)[:,]?\s*", r"\1 — ", s_)
            w(f"| [{t}]({rp(t)}) | Q{label} |")
        w("")
    still = collect("DECISION NEEDED")
    if still:
        w("## Still open")
        w("")
        w("| Task | Decision |")
        w("|---|---|")
        for t, s_ in still:
            w(f"| [{t}]({rp(t)}) | {s_} |")
        w("")
    (sdir / "README.md").write_text("\n".join(lines) + "\n")
    print(svc, "README written:", len(mine), "tasks,", len(seen), "rpcs")
