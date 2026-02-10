#!/usr/bin/env bash
set -e

# === Leer configuración desde stdin ===
payload=$(cat -)
TIME_LIMIT_MS=$(echo "$payload" | python3 -c "import sys,json; print(json.load(sys.stdin).get('time_limit_ms',1500))")
SRC=$(echo "$payload" | python3 -c "import sys,json; print(json.load(sys.stdin).get('source_file','/code/solution.cpp'))")

# === Preparar entorno ===
mkdir -p /runner
shopt -s nullglob
cases=()

# === Compilar el código ===
if ! g++ -O2 -std=c++17 -o /runner/main "$SRC" 2> /runner/compile.err; then
  python3 - <<'PY'
import json,sys
print(json.dumps({
  "status": "COMPILATION_ERROR",
  "timeMsTotal": 0,
  "cases": [],
  "stderr": open("/runner/compile.err").read()
}))
PY
  exit 0
fi

# === Ejecutar tests ===
total=0
ok=0
i=0

for inpath in $(ls /tests/*.in 2>/dev/null | sort); do
  i=$((i+1))
  # Emparejar output1.out o 1.out según formato
  outpath=$(echo "$inpath" | sed -E 's/(input)?([0-9]+)\.in$/output\2.out/')
  [ -f "$outpath" ] || outpath=$(echo "$inpath" | sed 's/\.in$/.out/')

  start=$(date +%s%3N)
  TIME_LIMIT_SEC=$(( (TIME_LIMIT_MS + 999) / 1000 ))
  timeout --signal=KILL "${TIME_LIMIT_SEC}s" /runner/main < "$inpath" > /runner/exec.out 2> /runner/exec.err || true
  end=$(date +%s%3N)
  elapsed=$((end - start))
  total=$((total + elapsed))

  actual=$(cat /runner/exec.out 2>/dev/null | sed -e 's/[ \t]*$//')
  expected=$(cat "$outpath" 2>/dev/null | sed -e 's/[ \t]*$//')

  if [ "$actual" = "$expected" ]; then
    status="OK"
    ok=$((ok + 1))
  else
    if [ -s /runner/exec.err ]; then
      status="RUNTIME_ERROR"
    else
      status="WRONG_ANSWER"
    fi
  fi

  stderr_escaped=$(cat /runner/exec.err 2>/dev/null | sed ':a;N;$!ba;s/\n/\\n/g')
  cases+=("{\"caseId\":$i,\"status\":\"$status\",\"timeMs\":$elapsed,\"stderr\":\"$stderr_escaped\"}")
done

# === Determinar resultado final ===
if [ $i -eq 0 ]; then
  final="NO_TESTS"
elif [ $i -eq $ok ]; then
  final="ACCEPTED"
elif [ $ok -eq 0 ]; then
  final="WRONG_ANSWER"
else
  final="PARTIAL"
fi

# === Salida final JSON (una sola línea) ===
echo -n "{\"status\":\"$final\",\"timeMsTotal\":$total,\"cases\":["$(IFS=,; echo "${cases[*]}")"]}"
