#!/usr/bin/env bash
# Rebuilds docs/EPICs/ROADMAP.md and the generated kitchen/planner READMEs from the task files.
# Run after adding tasks or changing any "Depends on" field: scripts/epics/roadmap.sh
set -euo pipefail
cd "$(dirname "$0")/../.."
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
python3 -I scripts/epics/check.py docs/EPICs
python3 -I scripts/epics/order.py docs/EPICs "$tmp" > /dev/null
python3 -I scripts/epics/render.py "$tmp" docs/EPICs
python3 -I scripts/epics/readmes.py "$tmp" docs/EPICs
