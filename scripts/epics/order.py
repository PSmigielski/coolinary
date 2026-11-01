#!/usr/bin/env python3
"""Order every task in docs/EPICs into milestones (written as JSON for render.py and readmes.py).

Usage: order.py <docs/EPICs> <out.json>. Milestones and the epic -> milestone map are the M, BASE, LATER and
EPIC_LATER tables below; a task moves to a later milestone when it depends on one there. A dependency on a task of a
later release is a note, not an ordering constraint."""
import json
import re
import sys
from collections import defaultdict
from pathlib import Path

root = Path(sys.argv[1])
out = Path(sys.argv[2]) if len(sys.argv) > 2 else None
P = "GW|IDN|NTF|ADM|KIT|PLN|OFF|AI|MOB"
TASK_FILE = re.compile(rf"^({P})-(\d\d)-(\d\d)-.+\.md$")
REF = re.compile(rf"\b({P})-(\d\d)(?:-(\d\d))?\b(\s*\(([^)]*)\))?")
SVC_ORDER = ["IDN", "NTF", "KIT", "PLN", "OFF", "AI", "ADM", "GW", "MOB"]
DAYS = {"S": 1.0, "M": 2.5, "L": 4.5}

# ---------------------------------------------------------------- milestones
M = [
    # id, release, title, done-when
    ("M1", "R1", "Foundations", "Kitchen's event outbox and the Go consumers run, the gateway maps errors and idempotency keys, and ai answers through Claude with structured output, cost tracking and batches; the app has its typed API client."),
    ("M2", "R1", "Sign-in, households and plans", "Identity publishes household events; a phone starts anonymously, signs in with an e-mail code, keeps its household, and the gateway knows Free vs Plus; founder alerts arrive by e-mail."),
    ("M3", "R1", "Back office", "An admin signs in at admin.<domain>; every action is in the audit log."),
    ("M4", "R1", "Catalogue and kitchen profile", "Ingredients are searchable in pl/de/en, names match to the catalogue, and onboarding saves cuisines, tools, diets and staples."),
    ("M5", "R1", "Fridge and barcodes", "Typing, barcode scan and unpack mode fill the fridge, offline edits sync (Verification 14)."),
    ("M6", "R1", "Nutrition and recipe library", "Starter dishes are imported or drafted, reviewed, show calories, and 'what can I make now' works with AI off (Verification 1, 16)."),
    ("M7", "R1", "AI intake", "Fridge photo, typed sentence and label photo become items to confirm; Free quota counts scans."),
    ("M8", "R1", "AI day plan", "A day plan comes from stock, passes the rule checker, falls back without AI, and swaps work (Verification 2-4)."),
    ("M9", "R1", "Shopping list and cooking", "The list holds exactly the missing items, ticking moves them into the fridge, 'cooked it' reduces stock (Verification 12)."),
    ("M10", "R1", "Notifications and reminders", "Push reaches a development build; Plus households get expiry reminders; preferences and quiet hours apply."),
    ("M11", "R1", "Account lifecycle and launch readiness", "Deletion and export work across services, quotas and rate limits hold, costs stay in the caps (Verification 11)."),
    ("M12", "R2", "Shops and flyer pipeline", "Shops from OpenStreetMap, postcode areas and editions exist; a flyer is fetched, read by batch AI, validated and held for review (Verification 8)."),
    ("M13", "R2", "Launch chains go live", "Each launch chain passes its terms gate and the 97% sample, and syncs for two weeks (Verification 5, 7)."),
    ("M14", "R2", "Offers in the app", "Offers nearest first, my shops, own flyer scan, badges on the list, plan from deals and the price matrix (Verification 6, 9, 10, 13)."),
    ("M15", "R2", "Shared household", "Invites by QR, code and link, profiles without accounts, leaving and removing, live shared list (Verification 15)."),
    ("M16", "R2", "Family goals and portions", "Per-person targets and one-dish portions, family day plans (Verification 17, 19, 20)."),
    ("M17", "R2", "Library growth, rewrites and savings", "Well-rated AI dishes join the library, tool rewrites, savings counters and monthly report."),
    ("M18", "R3", "Week plan", "Week plan with fewest cooking sessions and leftovers (Verification 18)."),
    ("M19", "R3", "Receipts and receipt prices", "Receipt photos become items and feed chain prices into the price matrix."),
    ("M20", "R3", "Depth features", "Cheapest top-up, restyle, share card and the remaining R3 items."),
    ("M21", "R3", "More chains and licensed feed", "Next chains and the licensed feed live behind their gates."),
]
MI = {m[0]: i for i, m in enumerate(M)}

# Base milestone per epic for its R1 (or own-release) tasks.
BASE = {
    "GW-01": "M1", "IDN-01": "M1", "KIT-01": "M1", "PLN-01": "M1", "NTF-01": "M1", "AI-01": "M1",
    "AI-02": "M1", "AI-03": "M1", "AI-04": "M1", "MOB-02": "M1", "MOB-03": "M1",
    "IDN-02": "M2", "IDN-03": "M2", "IDN-04": "M2", "IDN-05": "M2", "IDN-10": "M2", "NTF-02": "M2",
    "GW-02": "M2", "MOB-01": "M2", "GW-04": "M2", "NTF-05": "M2",
    "ADM-01": "M3", "ADM-02": "M3", "ADM-05": "M3", "ADM-07": "M3",
    "KIT-02": "M4", "KIT-03": "M4", "KIT-10": "M4", "AI-08": "M4", "AI-09": "M4", "GW-05": "M4", "ADM-03": "M4",
    "KIT-04": "M5", "KIT-05": "M5", "KIT-07": "M5", "GW-06": "M5", "MOB-04": "M5",
    "KIT-06": "M6", "PLN-02": "M6", "PLN-03": "M6", "PLN-04": "M6", "PLN-05": "M6", "PLN-06": "M6",
    "PLN-07": "M6", "ADM-04": "M6", "GW-07": "M6",
    "AI-05": "M7", "AI-06": "M7", "AI-07": "M7", "KIT-08": "M7", "GW-03": "M7", "MOB-07": "M7",
    "AI-10": "M8", "AI-11": "M8", "PLN-08": "M8", "PLN-09": "M8",
    "PLN-15": "M9", "PLN-16": "M9", "PLN-10": "M9", "GW-08": "M9",
    "NTF-03": "M10", "NTF-04": "M10", "NTF-06": "M10", "KIT-09": "M10", "GW-11": "M10", "MOB-06": "M10",
    "IDN-09": "M11", "KIT-13": "M11", "PLN-21": "M11",
    # R2
    "OFF-01": "M12", "OFF-02": "M12", "OFF-03": "M12", "OFF-05": "M12", "OFF-07": "M12", "OFF-08": "M12",
    "AI-13": "M12", "ADM-06": "M12",
    "OFF-06": "M13",
    "OFF-04": "M14", "OFF-09": "M14", "OFF-10": "M14", "OFF-11": "M14", "GW-10": "M14", "PLN-17": "M14",
    "PLN-18": "M14", "AI-14": "M14", "PLN-19": "M14",
    "IDN-06": "M15", "IDN-07": "M15", "IDN-08": "M15", "GW-09": "M15", "MOB-05": "M15", "OFF-14": "M15",
    "KIT-11": "M16", "PLN-12": "M16", "PLN-13": "M16", "AI-12": "M16",
    "PLN-11": "M17", "KIT-12": "M17", "GW-12": "M17",
    # R3
    "PLN-14": "M18", "AI-15": "M18",
    "AI-16": "M19",
    "PLN-20": "M20", "AI-17": "M20",
    "OFF-12": "M21", "OFF-13": "M21",
}
# Where a later-release task of an earlier-release epic goes.
LATER = {
    ("R2", "IDN"): "M15", ("R2", "GW"): "M14", ("R2", "MOB"): "M14", ("R2", "NTF"): "M15",
    ("R2", "KIT"): "M15", ("R2", "PLN"): "M17", ("R2", "AI"): "M17", ("R2", "ADM"): "M17", ("R2", "OFF"): "M14",
    ("R3", "KIT"): "M19", ("R3", "OFF"): "M19", ("R3", "GW"): "M20", ("R3", "MOB"): "M20",
    ("R3", "PLN"): "M20", ("R3", "AI"): "M20", ("R3", "ADM"): "M21", ("R3", "IDN"): "M20", ("R3", "NTF"): "M20",
}
EPIC_LATER = {  # finer placement for specific epics
    ("R2", "KIT-11"): "M16", ("R2", "KIT-12"): "M17", ("R2", "KIT-10"): "M16", ("R2", "GW-05"): "M16",
    ("R2", "PLN-15"): "M15", ("R2", "GW-08"): "M15", ("R2", "KIT-07"): "M15", ("R2", "GW-07"): "M16",
    ("R2", "MOB-03"): "M14", ("R2", "GW-04"): "M15", ("R2", "PLN-09"): "M17", ("R2", "AI-11"): "M17",
    ("R2", "GW-06"): "M14", ("R2", "MOB-07"): "M14", ("R2", "AI-09"): "M17", ("R2", "ADM-04"): "M17",
    ("R2", "GW-03"): "M14", ("R2", "NTF-03"): "M15", ("R2", "NTF-02"): "M15",
    ("R3", "KIT-08"): "M19", ("R3", "OFF-11"): "M19", ("R3", "PLN-19"): "M19", ("R3", "GW-06"): "M19",
    ("R3", "GW-07"): "M18", ("R3", "MOB-03"): "M18", ("R3", "PLN-09"): "M20", ("R3", "AI-11"): "M20",
    ("R3", "KIT-12"): "M20", ("R3", "GW-12"): "M20", ("R3", "GW-10"): "M20", ("R3", "MOB-07"): "M19",
    ("R3", "ADM-06"): "M21", ("R3", "GW-03"): "M19",
}

# ---------------------------------------------------------------- parse tasks
tasks = {}
by_epic = defaultdict(list)
for f in sorted(root.rglob("*.md")):
    m = TASK_FILE.match(f.name)
    if not m:
        continue
    pre, ep, tk = m.groups()
    tid = f"{pre}-{ep}-{tk}"
    text = f.read_text()
    title = text.splitlines()[0].removeprefix(f"# {tid}").strip(" —-")
    head = {}
    for line in text.splitlines()[:16]:
        mm = re.match(r"^\|\s*([A-Za-z ]+?)\s*\|\s*(.*?)\s*\|$", line)
        if mm:
            head[mm.group(1)] = mm.group(2)
    rel = (re.search(r"R[123]", head.get("Release", "")) or re.search("R1", "R1")).group(0)
    size = (re.match(r"[SML]", head.get("Size", "M")) or re.match("M", "M")).group(0)
    tasks[tid] = dict(id=tid, pre=pre, epic=f"{pre}-{ep}", title=title, rel=rel, size=size, status=head.get("Status", "todo"),
                      path=f.relative_to(root).as_posix(), deps_raw=head.get("Depends on", ""), text=text)
    by_epic[f"{pre}-{ep}"].append(tid)

RELN = {"R1": 1, "R2": 2, "R3": 3}

def sections(text):
    out, cur = {}, "head"
    for line in text.splitlines():
        if line.startswith("## "):
            cur = line[3:].strip().lower()
        out[cur] = out.get(cur, "") + line + "\n"
    return out

SEC = {tid: sections(t["text"]) for tid, t in tasks.items()}

def tokens(paren):
    toks = set(re.findall(r"`([^`]+)`", paren))
    toks |= set(re.findall(r"\b[A-Z][a-z]+[A-Z][A-Za-z0-9]*\b", paren))
    toks |= set(re.findall(r"\b[a-z]+\.[a-z_]+\b", paren))
    toks |= set(re.findall(r"/(?:v1|api)/[^\s,)`]*", paren))
    toks |= set(re.findall(r"\b[a-z]+_[a-z_]+\b", paren))
    out = set()
    for x in toks:
        if re.fullmatch(r"[A-Z]\w+\.[A-Z]\w+", x):
            x = x.split(".")[-1]
        if len(x) > 3:
            out.add(x)
    return out

def pick(epic, paren, own_rel):
    cands = [x for x in by_epic[epic] if RELN[tasks[x]["rel"]] <= RELN[own_rel]] or by_epic[epic]
    best = []
    for tok in tokens(paren or ""):
        strong, weak = [], []
        for x in cands:
            s = SEC[x]
            if re.search(rf"\brpc\s+{re.escape(tok)}\s*\(", tasks[x]["text"]) or tok.lower() in tasks[x]["title"].lower() \
               or any(tok in s.get(k, "") for k in ("contract", "data", "events", "goal")):
                strong.append(x)
            elif tok in tasks[x]["text"]:
                weak.append(x)
        hit = strong or weak
        if hit:
            best.append(min(hit))
    return best

unresolved = []
deps = defaultdict(set)
for tid, t in tasks.items():
    for pre, ep, tk, _, paren in REF.findall(t["deps_raw"]):
        if tk:
            ref = f"{pre}-{ep}-{tk}"
            if ref in tasks and ref != tid:
                deps[tid].add(ref)
            continue
        epic = f"{pre}-{ep}"
        if epic == t["epic"] or epic not in by_epic:
            continue
        found = pick(epic, paren, t["rel"])
        if not found:
            first = [x for x in by_epic[epic] if RELN[tasks[x]["rel"]] <= RELN[t["rel"]]] or by_epic[epic]
            found = [min(first)]
            if paren:
                unresolved.append(f"{tid}: {epic} ({paren})")
        deps[tid].update(x for x in found if x != tid)

# ---------------------------------------------------------------- forward references
# A dependency on a task of a later release is not an ordering constraint: the earlier task ships in its own
# release and the later task extends it. Keep them as notes.
forward = []
for tid in list(tasks):
    for d in sorted(deps[tid]):
        if RELN[tasks[d]["rel"]] > RELN[tasks[tid]["rel"]]:
            deps[tid].discard(d)
            forward.append((tid, d))

# ---------------------------------------------------------------- break cycles (drop conservative edges first)
def find_cycle():
    color, stack = {}, []
    def dfs(u):
        color[u] = 1
        stack.append(u)
        for v in sorted(deps[u]):
            if color.get(v) == 1:
                return stack[stack.index(v):] + [v]
            if color.get(v) is None:
                c = dfs(v)
                if c:
                    return c
        color[u] = 2
        stack.pop()
        return None
    sys.setrecursionlimit(10000)
    for n in sorted(tasks):
        if color.get(n) is None:
            c = dfs(n)
            if c:
                return c
    return None

cycles = []
while (c := find_cycle()):
    cycles.append(" -> ".join(c))
    # drop the cross-service edge in the cycle (or the last edge)
    edges = list(zip(c, c[1:]))
    cross = [e for e in edges if tasks[e[0]]["pre"] != tasks[e[1]]["pre"]]
    a, b = (cross or edges)[-1]
    deps[a].discard(b)

# ---------------------------------------------------------------- milestones
ms = {}
for tid, t in tasks.items():
    base = BASE.get(t["epic"], "M11")
    if RELN[t["rel"]] > RELN[M[MI[base]][1]]:
        base = EPIC_LATER.get((t["rel"], t["epic"]), LATER.get((t["rel"], t["pre"]), base))
    ms[tid] = MI[base]

changed = True
while changed:
    changed = False
    for tid in tasks:
        for d in deps[tid]:
            if ms[d] > ms[tid]:
                ms[tid] = ms[d]
                changed = True

inversions = [tid for tid in tasks if RELN[tasks[tid]["rel"]] < RELN[M[ms[tid]][1]]]

# ---------------------------------------------------------------- order inside milestones
order = []
for mi in range(len(M)):
    group = {t for t in tasks if ms[t] == mi}
    indeg = {t: len([d for d in deps[t] if d in group]) for t in group}
    key = lambda t: (SVC_ORDER.index(tasks[t]["pre"]), t)
    ready = sorted([t for t in group if indeg[t] == 0], key=key)
    seq = []
    while ready:
        u = ready.pop(0)
        seq.append(u)
        for v in sorted(group):
            if u in deps[v]:
                indeg[v] -= 1
                if indeg[v] == 0:
                    ready.append(v)
                    ready.sort(key=key)
    assert len(seq) == len(group), "cycle inside milestone"
    order.append(seq)

stats = dict(
    total=len(tasks),
    per_milestone=[(M[i][0], len(s), round(sum(DAYS[tasks[t]["size"]] for t in s), 1)) for i, s in enumerate(order)],
    unresolved=len(unresolved), cycles=len(cycles), inversions=len(inversions),
)
print(json.dumps(stats, indent=1))
print("cycles:", *cycles[:20], sep="\n  ")
print("inversions:", *[f"{t} ({tasks[t]['rel']}) -> {M[ms[t]][0]} via {sorted(d for d in deps[t] if ms[d] == ms[t] and RELN[tasks[d]['rel']] > RELN[tasks[t]['rel']])}" for t in inversions[:40]], sep="\n  ")
print("unresolved sample:", *unresolved[:25], sep="\n  ")

if out:
    json.dump(dict(order=order, ms={k: M[v][0] for k, v in ms.items()},
                   tasks={k: {x: v[x] for x in ("title", "rel", "size", "path", "pre", "status")} for k, v in tasks.items()},
                   deps={k: sorted(v) for k, v in deps.items()}, inversions=inversions, cycles=cycles, forward=forward,
                   unresolved=unresolved, M=M), out.open("w"), indent=1)
