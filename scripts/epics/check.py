#!/usr/bin/env python3
"""Validate docs/EPICs task files: names, header table, sections, dependency references."""
import re
import sys
from collections import Counter, defaultdict
from pathlib import Path

root = Path(sys.argv[1])
PREFIXES = "GW|IDN|NTF|ADM|KIT|PLN|OFF|AI|MOB"
TASK_RE = re.compile(rf"^({PREFIXES})-(\d\d)-(\d\d)-[a-z0-9-]+\.md$")
REF_RE = re.compile(rf"\b({PREFIXES})-(\d\d)(?:-(\d\d))?\b")
REQUIRED = ["## Goal", "## Business rules", "## Acceptance criteria", "## Errors"]
FIELDS = ["Service", "Epic", "Release", "Size", "Depends on"]

tasks, epics, problems = {}, set(), []
stats = Counter()
sizes = Counter()
releases = Counter()
decisions = defaultdict(int)

for f in sorted(root.rglob("*.md")):
    rel = f.relative_to(root)
    if f.name in ("README.md", "TEMPLATE.md", "ROADMAP.md", "DECISIONS.md"):
        continue
    m = TASK_RE.match(f.name)
    if not m:
        problems.append(f"bad file name: {rel}")
        continue
    pre, ep, tk = m.groups()
    tid = f"{pre}-{ep}-{tk}"
    if tid in tasks:
        problems.append(f"duplicate id {tid}: {rel} and {tasks[tid]['path']}")
    folder = f.parent.name
    if not folder.startswith(f"E{ep}-"):
        problems.append(f"{tid} sits in folder {folder}")
    if not (f.parent / "README.md").exists():
        problems.append(f"missing epic README: {f.parent.relative_to(root)}")
    text = f.read_text()
    first = text.splitlines()[0] if text else ""
    if not first.startswith(f"# {tid} "):
        problems.append(f"{tid}: first line is not '# {tid} <title>'")
    for s in REQUIRED:
        if s not in text:
            problems.append(f"{tid}: missing section '{s}'")
    head = {}
    for line in text.splitlines()[:16]:
        mm = re.match(r"^\|\s*([A-Za-z ]+?)\s*\|\s*(.*?)\s*\|$", line)
        if mm:
            head[mm.group(1)] = mm.group(2)
    for fld in FIELDS:
        if fld not in head:
            problems.append(f"{tid}: header field '{fld}' missing")
    size = head.get("Size", "")
    rel_tag = head.get("Release", "")
    if size not in ("S", "M", "L"):
        problems.append(f"{tid}: size '{size}'")
    if not re.match(r"R[123]\b", rel_tag):
        problems.append(f"{tid}: release '{rel_tag}'")
    sizes[size] += 1
    releases[rel_tag] += 1
    decisions[pre] += text.count("DECISION NEEDED")
    tasks[tid] = {"path": rel, "deps": head.get("Depends on", ""), "pre": pre}
    epics.add(f"{pre}-{ep}")
    stats[pre] += 1

# Epic-level ids that exist as folders (even if empty) count too.
for d in root.glob("*/E*-*"):
    if d.is_dir():
        for f in d.glob("*.md"):
            m = TASK_RE.match(f.name)
            if m:
                epics.add(f"{m.group(1)}-{m.group(2)}")

dangling = []
for tid, t in tasks.items():
    for pre, ep, tk in REF_RE.findall(t["deps"]):
        ref = f"{pre}-{ep}-{tk}" if tk else f"{pre}-{ep}"
        ok = ref in tasks if tk else ref in epics
        if not ok:
            dangling.append(f"{tid} depends on missing {ref}")

print("tasks per service:", dict(stats), "total", sum(stats.values()))
print("sizes:", dict(sizes), " releases:", dict(releases))
print("DECISION NEEDED mentions:", dict(decisions))
print(f"\n{len(problems)} format problems")
for p in problems[:80]:
    print("  ", p)
print(f"\n{len(dangling)} dangling dependencies")
for d in dangling[:120]:
    print("  ", d)
