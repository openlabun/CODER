#!/usr/bin/env bash
set -euo pipefail

# E2E worker test by language.
# Flow:
# 1) Create challenge with test cases
# 2) Register/login a throwaway student
# 3) Create submissions (one or many languages)
# 4) Poll submission status until terminal state
#
# Usage:
#   ./test_worker_languages.sh
#   ./test_worker_languages.sh --language python
#   ./test_worker_languages.sh --api-url http://localhost:3000 --timeout 180

API_URL="http://localhost:3000"
LANGUAGE="all"
TIMEOUT_SECONDS=120
POLL_INTERVAL_SECONDS=2

usage() {
  cat <<EOF
Usage: $0 [options]

Options:
  --api-url URL         API base url (default: http://localhost:3000)
  --language LANG       one of: python|node|cpp|java|all (default: all)
  --timeout SEC         max wait per submission in seconds (default: 120)
  --poll-interval SEC   polling interval in seconds (default: 2)
  --help                show this message
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --api-url) API_URL="$2"; shift 2 ;;
    --language) LANGUAGE="$2"; shift 2 ;;
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

normalize_lang() {
  case "$1" in
    python|node|cpp|java|all) echo "$1" ;;
    javascript) echo "node" ;;
    c++|cxx) echo "cpp" ;;
    *)
      echo "Invalid language: $1. Use python|node|cpp|java|all" >&2
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

is_success_status() {
  [[ "$1" == "accepted" ]]
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

if [[ "$LANGUAGE" == "all" ]]; then
  LANGUAGES=(python node cpp java)
else
  LANGUAGES=("$LANGUAGE")
fi

echo "[info] API URL: $API_URL"
echo "[info] Languages: ${LANGUAGES[*]}"

# 1) Create challenge with shared test cases.
CHALLENGE_PAYLOAD=$(build_json "
const body = {
  title: 'Worker Multi-Language Smoke',
  description: 'Given two integers, print their sum.',
  difficulty: 'easy',
  timeLimit: 2000,
  memoryLimit: 256,
  tags: ['worker', 'e2e', 'smoke'],
  publicTestCases: [{ name: 'sample-1', input: '2 3\\n', output: '5\\n' }],
  hiddenTestCases: [
    { name: 'hidden-1', input: '10 32\\n', output: '42\\n' },
    { name: 'hidden-2', input: '-5 2\\n', output: '-3\\n' }
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

# 2) Register a throwaway student and get token.
TS="$(date +%s)"
RAND="${RANDOM:-1111}"
TEST_USER="worker_test_${TS}_${RAND}"
TEST_PASS="pass_${TS}_${RAND}"

REGISTER_PAYLOAD=$(build_json "
const body = {
  username: process.env.TEST_USER,
  password: process.env.TEST_PASS,
  role: 'student'
};
process.stdout.write(JSON.stringify(body));
" )

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

# 3) Submit and poll per language.
declare -A FINAL_STATUS
FAIL_COUNT=0

for lang in "${LANGUAGES[@]}"; do
  echo "[info] submitting language: $lang"
  code="$(get_code_for_language "$lang")"

  submission_payload=$(CHALLENGE_ID="$CHALLENGE_ID" LANG="$lang" CODE="$code" build_json "
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
    echo "[error] Submission failed for $lang (HTTP $sub_status)" >&2
    echo "$sub_body" >&2
    FINAL_STATUS["$lang"]="http_${sub_status}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
    continue
  fi

  sub_id="$(json_get "$sub_body" "id")"
  echo "[ok] queued submission: $sub_id"

  start_ts="$(date +%s)"
  current_status="queued"

  while true; do
    get_resp=$(http_get "$API_URL/submissions/$sub_id")
    get_status="$(printf '%s' "$get_resp" | head -n1)"
    get_body="$(printf '%s' "$get_resp" | tail -n +2)"

    if [[ "$get_status" -lt 200 || "$get_status" -ge 300 ]]; then
      echo "[error] Cannot fetch submission $sub_id for $lang (HTTP $get_status)" >&2
      echo "$get_body" >&2
      current_status="http_${get_status}"
      break
    fi

    if status_val="$(json_get "$get_body" "status" 2>/dev/null)"; then
      current_status="$status_val"
    else
      current_status="unknown"
    fi

    elapsed=$(( $(date +%s) - start_ts ))
    echo "[info] $lang submission=$sub_id status=$current_status elapsed=${elapsed}s"

    if is_terminal_status "$current_status"; then
      break
    fi

    if [[ "$elapsed" -ge "$TIMEOUT_SECONDS" ]]; then
      current_status="timeout"
      break
    fi

    sleep "$POLL_INTERVAL_SECONDS"
  done

  FINAL_STATUS["$lang"]="$current_status"

  if ! is_success_status "$current_status"; then
    FAIL_COUNT=$((FAIL_COUNT + 1))
  fi
done

echo
echo "================ Worker Language Test Summary ================"
for lang in "${LANGUAGES[@]}"; do
  echo "$lang: ${FINAL_STATUS[$lang]:-not_run}"
done
echo "Challenge ID: $CHALLENGE_ID"

echo
if [[ "$FAIL_COUNT" -gt 0 ]]; then
  echo "[result] FAIL ($FAIL_COUNT language(s) not accepted)"
  echo "Hint: ensure API, Redis, worker and runner images are up." >&2
  exit 2
fi

echo "[result] PASS (all languages accepted)"
