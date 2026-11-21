#!/usr/bin/env bash
# MUT-GATE: the test must go RED on the pre-fix implementation.
set -uo pipefail
SRC="$(cd "$(dirname "$0")" && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
cp "$SRC"/go.mod "$SRC"/origin.go "$SRC"/origin_test.go "$TMP/"
sed -i 's/^var normalize = normalizeFixed$/var normalize = normalizeBuggy/' "$TMP/origin.go"
if grep -q '^var normalize = normalizeBuggy$' "$TMP/origin.go"; then
  echo 'seam flipped to normalizeBuggy'
else
  echo 'could not flip the seam; script is stale'
  exit 2
fi
cd "$TMP" || exit 1
if go test ./... > /tmp/x402_must_fail.log 2>&1; then
  echo 'MUTANT SURVIVED: tests passed on the pre-fix implementation'
  echo '=> the repro does not measure the bug; do not trust the README.'
  exit 1
else
  echo 'MUTANT KILLED: tests failed on the pre-fix implementation (expected)'
  grep -E '^\s+--- FAIL|FAIL\s' /tmp/x402_must_fail.log | head -5
  exit 0
fi
