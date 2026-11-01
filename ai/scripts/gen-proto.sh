#!/bin/sh
# Generates protobuf + grpcio code for every contract (committed under src/coolinary).
# Runs inside the coolinary/protogen image (make proto); grpcio-tools there is pinned
# to the same version as this service's grpcio runtime. OUT defaults to src/.
set -eu
cd "$(dirname "$0")/.."
OUT="${OUT:-src}"
mkdir -p "$OUT"
rm -rf "$OUT/coolinary"
python -m grpc_tools.protoc -I ../contracts \
  --python_out="$OUT" --pyi_out="$OUT" --grpc_python_out="$OUT" \
  $(cd ../contracts && find coolinary -name '*.proto' | sort)
# Regular packages (with __init__.py) so packaging tools pick them up.
find "$OUT/coolinary" -type d -exec touch {}/__init__.py \;
