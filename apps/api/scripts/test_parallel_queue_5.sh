#!/usr/bin/env bash
set -euo pipefail

# Burst queue test: enqueue N submissions quickly and poll all until completion.
# Usage:
#   ./test_parallel_queue_5.sh
#   ./test_parallel_queue_5.sh --language python --count 5
#   ./test_parallel_queue_5.sh --api-url http://localhost:3000 --timeout 180

API_URL="http://localhost:3000"
LANGUAGE="python"
SUBMISSION_COUNT=5
TIMEOUT_SECONDS=180
POLL_INTERVAL_SECONDS=2

usage() {
  cat <<EOF
Usage: $0 [options]

Options:
  --api-url URL         API base url (default: http://localhost:3000)
  --language LANG       one of: python|node|cpp|java (default: python)
  --count N             number of submissions to enqueue (default: 5)
  --timeout SEC         max total wait in seconds (default: 180)
  --poll-interval SEC   polling interval in seconds (default: 2)
  --help                show this message
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --api-url) API_URL="$2"; shift 2 ;;
    --language) LANGUAGE="$2"; shift 2 ;;
    --count) SUBMISSION_COUNT="$2"; shift 2 ;;
    --timeout) TIMEOUT_SECONDS="$2"; shift 2 ;;
    --poll-interval) POLL_INTERVAL_SECONDS="$2"; shift 2 ;;
    --help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage; exit 1 ;;
  esac
done

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2
  exit 1
fi

if ! command -v node >/dev/null 2>&1; then
  echo "node is required to build/parse JSON payloads" >&2
  exit 1
fi

if ! [[ "$SUBMISSION_COUNT" =~ ^[0-9]+$ ]] || [[ "$SUBMISSION_COUNT" -lt 1 ]]; then
  echo "--count must be a positive integer" >&2
  exit 1
fi

normalize_lang() {
  case "$1" in
    python|node|cpp|java) echo "$1" ;;
    javascript) echo "node" ;;
    c++|cxx) echo "cpp" ;;
    *)
      echo "Invalid language: $1. Use python|node|cpp|java" >&2
      exit 1
      ;;
  esac
}

LANGUAGE="$(normalize_lang "$LANGUAGE")"

http_post() {
  local url="$1"
  local json_payload="$2"
  local token="${3:-}"
  local response

  if [[ -n "$token" ]]; then
    response=$(curl -sS -X POST "$url" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $token" \
      -d "$json_payload" \
      -w "\n%{http_code}")
  else
    response=$(curl -sS -X POST "$url" \
      -H "Content-Type: application/json" \
      -d "$json_payload" \
      -w "\n%{http_code}")
  fi

  local status="${response##*$'\n'}"
  local body="${response%$'\n'*}"
  printf '%s\n%s' "$status" "$body"
}

http_get() {
  local url="$1"
  local token="${2:-}"
  local response

  if [[ -n "$token" ]]; then
    response=$(curl -sS -X GET "$url" \
      -H "Authorization: Bearer $token" \
      -w "\n%{http_code}")
  else
    response=$(curl -sS -X GET "$url" -w "\n%{http_code}")
  fi

  local status="${response##*$'\n'}"
  local body="${response%$'\n'*}"
  printf '%s\n%s' "$status" "$body"
}

json_get() {
  local json_input="$1"
  local key="$2"
  printf '%s' "$json_input" | node -e "
let data = '';
process.stdin.on('data', c => data += c);
process.stdin.on('end', () => {
  try {
    const obj = JSON.parse(data || '{}');
    const v = obj[process.argv[1]];
    if (v === undefined || v === null) process.exit(2);
    process.stdout.write(typeof v === 'string' ? v : JSON.stringify(v));
  } catch {
    process.exit(2);
  }
});
" "$key"
}

build_json() {
  node -e "$1"
}

is_terminal_status() {
  case "$1" in
    accepted|wrong_answer|error) return 0 ;;
    *) return 1 ;;
  esac
}

read -r -d '' PYTHON_CODE <<'EOF' || true
import sys
nums = list(map(int, sys.stdin.read().split()))
print(sum(nums[:2]))
EOF

read -r -d '' NODE_CODE <<'EOF' || true
const fs = require('fs');
const nums = fs.readFileSync(0, 'utf8').trim().split(/\s+/).map(Number);
console.log((nums[0] || 0) + (nums[1] || 0));
EOF

read -r -d '' CPP_CODE <<'EOF' || true
#include <bits/stdc++.h>
using namespace std;
int main() {
    long long a, b;
    if (!(cin >> a >> b)) return 0;
    cout << (a + b) << "\n";
    return 0;
}
EOF

read -r -d '' JAVA_CODE <<'EOF' || true
import java.util.*;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        if (!sc.hasNextLong()) return;
        long a = sc.nextLong();
        long b = sc.hasNextLong() ? sc.nextLong() : 0;
        System.out.println(a + b);
    }
}
EOF

get_code_for_language() {
  case "$1" in
    python) printf '%s' "$PYTHON_CODE" ;;
    node) printf '%s' "$NODE_CODE" ;;
    cpp) printf '%s' "$CPP_CODE" ;;
    java) printf '%s' "$JAVA_CODE" ;;
    *)
      echo "Unsupported language: $1" >&2
      exit 1
      ;;
  esac
}

echo "[info] API URL: $API_URL"
echo "[info] Language: $LANGUAGE"
echo "[info] Submission burst size: $SUBMISSION_COUNT"

CHALLENGE_PAYLOAD=$(build_json "
const body = {
  title: 'Parallel Queue Burst Test',
  description: 'Given two integers, print their sum.',
  difficulty: 'easy',
  timeLimit: 2000,
  memoryLimit: 256,
  tags: ['queue', 'parallel', 'burst'],
  publicTestCases: [{ name: 'sample-1', input: '2 3\\n', output: '5\\n' }],
  hiddenTestCases: [
    { name: 'hidden-1', input: '10 32\\n', output: '42\\n' },
    { name: 'hidden-2', input: '100 200\\n', output: '300\\n' }
  ]
};
process.stdout.write(JSON.stringify(body));
")

challenge_resp=$(http_post "$API_URL/challenges" "$CHALLENGE_PAYLOAD")
challenge_status="$(printf '%s' "$challenge_resp" | head -n1)"
challenge_body="$(printf '%s' "$challenge_resp" | tail -n +2)"

if [[ "$challenge_status" -lt 200 || "$challenge_status" -ge 300 ]]; then
  echo "[error] Challenge creation failed (HTTP $challenge_status)" >&2
  echo "$challenge_body" >&2
  exit 1
fi

CHALLENGE_ID="$(json_get "$challenge_body" "id")"
echo "[ok] challenge created: $CHALLENGE_ID"

TS="$(date +%s)"
RAND="${RANDOM:-1111}"
TEST_USER="queue_test_${TS}_${RAND}"
TEST_PASS="pass_${TS}_${RAND}"

REGISTER_PAYLOAD=$(build_json "
const body = {
  username: process.env.TEST_USER,
  password: process.env.TEST_PASS,
  role: 'student'
};
process.stdout.write(JSON.stringify(body));
")

reg_resp=$(TEST_USER="$TEST_USER" TEST_PASS="$TEST_PASS" http_post "$API_URL/auth/register" "$REGISTER_PAYLOAD")
reg_status="$(printf '%s' "$reg_resp" | head -n1)"
reg_body="$(printf '%s' "$reg_resp" | tail -n +2)"

if [[ "$reg_status" -lt 200 || "$reg_status" -ge 300 ]]; then
  echo "[error] User registration failed (HTTP $reg_status)" >&2
  echo "$reg_body" >&2
  exit 1
fi

ACCESS_TOKEN="$(json_get "$reg_body" "accessToken")"
echo "[ok] student registered: $TEST_USER"

SOURCE_CODE="$(get_code_for_language "$LANGUAGE")"
submission_ids=()
submission_statuses=()
start_times=()

enqueue_start=$(date +%s)
for ((i=1; i<=SUBMISSION_COUNT; i++)); do
  submission_payload=$(CHALLENGE_ID="$CHALLENGE_ID" LANG="$LANGUAGE" CODE="$SOURCE_CODE" build_json "
const body = {
  challengeId: process.env.CHALLENGE_ID,
  language: process.env.LANG,
  code: process.env.CODE
};
process.stdout.write(JSON.stringify(body));
")

  sub_resp=$(http_post "$API_URL/submissions" "$submission_payload" "$ACCESS_TOKEN")
  sub_status="$(printf '%s' "$sub_resp" | head -n1)"
  sub_body="$(printf '%s' "$sub_resp" | tail -n +2)"

  if [[ "$sub_status" -lt 200 || "$sub_status" -ge 300 ]]; then
    echo "[error] Submission #$i failed (HTTP $sub_status)" >&2
    echo "$sub_body" >&2
    exit 1
  fi

  sub_id="$(json_get "$sub_body" "id")"
  submission_ids+=("$sub_id")
  submission_statuses+=("queued")
  start_times+=("$(date +%s)")
  echo "[ok] queued #$i -> $sub_id"
done
enqueue_end=$(date +%s)

echo "[info] enqueue burst completed in $((enqueue_end - enqueue_start))s"

global_start=$(date +%s)
while true; do
  all_done=1
  line_parts=()
  now=$(date +%s)

  for ((i=0; i<SUBMISSION_COUNT; i++)); do
    current="${submission_statuses[$i]}"
    if is_terminal_status "$current" || [[ "$current" == "timeout" || "$current" == "poll_error" ]]; then
      line_parts+=("#$((i+1))=$current")
      continue
    fi

    all_done=0
    sub_id="${submission_ids[$i]}"
    get_resp=$(http_get "$API_URL/submissions/$sub_id")
    get_status="$(printf '%s' "$get_resp" | head -n1)"
    get_body="$(printf '%s' "$get_resp" | tail -n +2)"

    if [[ "$get_status" -lt 200 || "$get_status" -ge 300 ]]; then
      current="poll_error"
    else
      if current_val="$(json_get "$get_body" "status" 2>/dev/null)"; then
        current="$current_val"
      else
        current="unknown"
      fi
    fi

    elapsed=$(( now - ${start_times[$i]} ))
    global_elapsed=$(( now - global_start ))
    if ! is_terminal_status "$current" && [[ "$elapsed" -ge "$TIMEOUT_SECONDS" || "$global_elapsed" -ge "$TIMEOUT_SECONDS" ]]; then
      current="timeout"
    fi

    submission_statuses[$i]="$current"
    line_parts+=("#$((i+1))=$current")
  done

  echo "[info] poll t=$(( $(date +%s) - global_start ))s :: ${line_parts[*]}"

  unfinished=0
  for status in "${submission_statuses[@]}"; do
    if ! is_terminal_status "$status" && [[ "$status" != "timeout" && "$status" != "poll_error" ]]; then
      unfinished=1
      break
    fi
  done

  if [[ "$all_done" -eq 1 || "$unfinished" -eq 0 ]]; then
    break
  fi

  if [[ $(( $(date +%s) - global_start )) -ge "$TIMEOUT_SECONDS" ]]; then
    break
  fi

  sleep "$POLL_INTERVAL_SECONDS"
done

accepted=0
wrong_answer=0
error_count=0
timeout_count=0
poll_error_count=0

for status in "${submission_statuses[@]}"; do
  case "$status" in
    accepted) accepted=$((accepted + 1)) ;;
    wrong_answer) wrong_answer=$((wrong_answer + 1)) ;;
    error) error_count=$((error_count + 1)) ;;
    timeout) timeout_count=$((timeout_count + 1)) ;;
    poll_error) poll_error_count=$((poll_error_count + 1)) ;;
  esac
done

echo
echo "================ Parallel Queue Summary ================"
for ((i=0; i<SUBMISSION_COUNT; i++)); do
  elapsed=$(( $(date +%s) - ${start_times[$i]} ))
  echo "#$((i+1)) id=${submission_ids[$i]} status=${submission_statuses[$i]} elapsed=${elapsed}s"
done

echo ""
echo "accepted=$accepted wrong_answer=$wrong_answer error=$error_count timeout=$timeout_count poll_error=$poll_error_count"
echo "challengeId=$CHALLENGE_ID"

echo ""
if [[ "$accepted" -ne "$SUBMISSION_COUNT" ]]; then
  echo "[result] FAIL (expected $SUBMISSION_COUNT accepted, got $accepted)." >&2
  exit 2
fi

echo "[result] PASS (all $SUBMISSION_COUNT submissions accepted)"
