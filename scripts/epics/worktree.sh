#!/usr/bin/env bash
# Git worktrees for building roadmap tasks in parallel (/next-task batches). Each task gets
# .claude/worktrees/<ID> on branch task/<ID>, made from the current HEAD. Its tests use their own
# <svc>_test_<id> databases (scripts/test-service.sh), so parallel tasks never share a test database.
# The main session then lands each finished task on the milestone branch as one commit.
#
# Usage (from the main checkout, not from a worktree):
#   worktree.sh add <ID>            create the task's worktree (or reuse it) and print its path
#   worktree.sh land <ID>           apply the task's commits to the main checkout, staged but not committed;
#                                   on a conflict exit 1 and list the conflicting files
#   worktree.sh abort <ID>          undo a land that hit a conflict (other changes in the tree are kept)
#   worktree.sh finish <ID> [note]  after land (and any conflict fixes, staged): mark the task done, update the
#                                   roadmap, commit with the task's own message (note becomes a "Landed:"
#                                   trailer), then remove the worktree
#   worktree.sh remove <ID>         delete the worktree, its branch and its test databases
#   worktree.sh list                task worktrees and how many commits each has
set -euo pipefail
cd "$(dirname "$0")/../.."

die() { echo "worktree.sh: $*" >&2; exit 1; }
tasks() { python3 -I scripts/epics/tasks.py "$@"; }

[[ "$(git rev-parse --path-format=absolute --git-dir)" == "$(git rev-parse --path-format=absolute --git-common-dir)" ]] ||
  die "run it from the main checkout, not from a worktree"

state=.claude/state
landing="$state/landing.json"

# Same rule as scripts/test-service.sh: the test-database suffix of a worktree, from its folder name.
slug() { tr 'A-Z' 'a-z' <<<"$1" | tr -c 'a-z0-9\n' '_' | cut -c1-40; }

# Task files are staged only by their own task; the roadmap and generated READMEs go with every task.
epic_files() {
  local file="$1"
  { echo "$file"; git diff --name-only -- docs/EPICs | grep -v -E '/[A-Z]+-[0-9]{2}-[0-9]{2}-[^/]*\.md$' || true; } | sort -u
}

cmd_add() {
  local id="$1" wt=".claude/worktrees/$1"
  tasks show "$id" >/dev/null
  if [[ -d "$wt" ]]; then
    echo "$PWD/$wt"
    return
  fi
  git worktree prune
  if git rev-parse -q --verify "refs/heads/task/$id" >/dev/null; then
    git worktree add -q "$wt" "task/$id"
  else
    git worktree add -q -b "task/$id" "$wt" HEAD
  fi
  echo "$PWD/$wt"
}

cmd_land() {
  local id="$1" br="task/$1"
  git rev-parse -q --verify "refs/heads/$br" >/dev/null || die "$id: no branch $br"
  [[ ! -f "$landing" ]] || die "another land is open ($(jq -r .id "$landing")): finish or abort it first"
  [[ -z "$(git diff --cached --name-only)" ]] || die "the index has staged changes; commit or unstage them first"
  local base pick
  base="$(git merge-base HEAD "$br")"
  mapfile -t commits < <(git rev-list --reverse "$base..$br")
  ((${#commits[@]})) || die "$id: nothing committed on $br"
  pick="${commits[0]}"
  if ((${#commits[@]} > 1)); then
    # Squash the task's commits into one, with the first commit's message (no reset, no rebase).
    pick="$(git log -1 --format=%B "${commits[0]}" | git commit-tree "$br^{tree}" -p "$base")"
  fi
  mkdir -p "$state"
  jq -n --arg id "$id" --arg base "$base" --arg pick "$pick" '{id: $id, base: $base, pick: $pick}' >"$landing"
  if git cherry-pick --no-commit "$pick" >/dev/null 2>&1; then
    echo "landed $id (staged, not committed): $(git diff --cached --name-only | wc -l) files"
    return
  fi
  local conflicts
  conflicts="$(git diff --name-only --diff-filter=U)"
  if [[ -z "$conflicts" ]]; then
    rm -f "$landing"
    die "$id: cherry-pick refused to start (local changes or untracked files in the way?); nothing was changed"
  fi
  echo "conflict landing $id; resolve these, git add them, then: worktree.sh finish $id (or abort $id)"
  while read -r f; do
    case "$f" in
      */gen/* | */src/gen/* | ai/src/coolinary/* | contracts/openapi/* | */store/* | */auth-schema.ts | mobile/src/api/*) echo "  generated  $f" ;;
      *.lock | */bun.lock | */uv.lock | */go.sum) echo "  lockfile   $f" ;;
      */migrations/* | */drizzle/*) echo "  migration  $f" ;;
      *) echo "  code       $f" ;;
    esac
  done <<<"$conflicts"
  exit 1
}

cmd_abort() {
  local id="$1"
  [[ -f "$landing" && "$(jq -r .id "$landing")" == "$id" ]] || die "$id: no open land"
  local base pick
  base="$(jq -r .base "$landing")"
  pick="$(jq -r .pick "$landing")"
  # The land started only on paths without local changes, so putting the pick's paths back to HEAD
  # undoes it and keeps everything else.
  mapfile -t paths < <(git diff --no-renames --name-only "$base" "$pick")
  for f in "${paths[@]}"; do
    if git cat-file -e "HEAD:$f" 2>/dev/null; then
      git restore --source=HEAD --staged --worktree -- "$f"
    else
      git rm -q -f --cached --ignore-unmatch -- "$f"
      rm -f -- "$f"
    fi
  done
  rm -f "$landing"
  echo "aborted landing $id"
}

cmd_finish() {
  local id="$1" note="${2:-}"
  [[ -f "$landing" && "$(jq -r .id "$landing")" == "$id" ]] || die "$id: no open land (run worktree.sh land $id first)"
  [[ -z "$(git diff --name-only --diff-filter=U)" ]] || die "$id: unresolved conflicts: $(git diff --name-only --diff-filter=U | xargs)"
  local pick file msg
  pick="$(jq -r .pick "$landing")"
  file="$(tasks show "$id" | jq -r .file)"
  tasks done "$id" >/dev/null
  # Render the roadmap from the statuses this commit holds: the batch's other task files (still "in progress",
  # maybe with a decision recorded) are put back to HEAD while it renders, then restored.
  keep="$(mktemp -d)"
  mapfile -t others < <(git diff --name-only -- docs/EPICs | grep -E '/[A-Z]+-[0-9]{2}-[0-9]{2}-[^/]*\.md$' | grep -v -x -F "$file" || true)
  for f in "${others[@]}"; do
    mkdir -p "$keep/$(dirname "$f")" && cp -p "$f" "$keep/$f" && git show "HEAD:$f" >"$f"
  done
  restore() { for f in "${others[@]}"; do cp -p "$keep/$f" "$f"; done; rm -rf "$keep"; }
  trap restore EXIT
  scripts/epics/roadmap.sh >/dev/null
  mapfile -t docs < <(epic_files "$file")
  git add -- "${docs[@]}"
  restore
  trap - EXIT
  msg="$(mktemp)"
  git log -1 --format=%B "$pick" >"$msg"
  if [[ -n "$note" ]]; then
    git commit -q -F "$msg" --trailer "Landed: $note"
  else
    git commit -q -F "$msg"
  fi
  rm -f "$msg" "$landing"
  echo "committed $id: $(git log -1 --format='%h %s')"
  cmd_remove "$id"
}

cmd_remove() {
  local id="$1" wt=".claude/worktrees/$1"
  if [[ -d "$wt" ]]; then
    git worktree remove --force "$wt"
  fi
  git worktree prune
  git branch -q -D "task/$id" 2>/dev/null || true
  local dev=(docker compose -f docker-compose.dev.yml)
  if [[ -n "$("${dev[@]}" ps --status running -q postgres 2>/dev/null)" ]]; then
    local suffix="_test_$(slug "$id")"
    while read -r db; do
      [[ -n "$db" ]] && "${dev[@]}" exec -T postgres psql -U postgres -qtA -c "DROP DATABASE IF EXISTS \"$db\" WITH (FORCE)" >/dev/null
    done < <("${dev[@]}" exec -T postgres psql -U postgres -qtA \
      -c "SELECT datname FROM pg_database WHERE right(datname, ${#suffix}) = '$suffix'")
  fi
  echo "removed worktree $id"
}

cmd_list() {
  git worktree list --porcelain | awk '/^worktree /{p=$2} /^branch refs\/heads\/task\//{sub("refs/heads/","",$2); print p, $2}' |
    while read -r path br; do
      echo "${br#task/}  $(git rev-list --count "$(git merge-base HEAD "$br")..$br") commits  $path"
    done
}

cmd="${1:-}"
case "$cmd" in
  add | land | abort | remove) [[ -n "${2:-}" ]] || die "usage: worktree.sh $cmd <ID>"; "cmd_$cmd" "$2" ;;
  finish) [[ -n "${2:-}" ]] || die "usage: worktree.sh finish <ID> [note]"; cmd_finish "$2" "${3:-}" ;;
  list) cmd_list ;;
  *) sed -n '2,16p' "$0" | sed 's/^# \{0,1\}//'; exit 2 ;;
esac
