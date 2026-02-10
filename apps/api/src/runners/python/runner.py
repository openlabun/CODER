import json, sys, time, os, subprocess, glob


def run_all_tests(time_limit_ms):
    tests = sorted(glob.glob("/tests/input*.in"))
    cases = []
    total_time = 0
    ok_count = 0
    for idx, inpath in enumerate(tests, start=1):
        outpath = inpath.replace("/input", "/output").replace(".in", ".out")
        start = time.time()
        try:
            p = subprocess.run(
                ["python3", sys.argv[1]]
                if len(sys.argv) > 1
                else ["python3", "/code/solution.py"],
                stdin=open(inpath, "rb"),
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=time_limit_ms / 1000,
            )
            elapsed = int((time.time() - start) * 1000)
            actual = p.stdout.decode().strip()
            expected = (
                open(outpath, "r").read().strip() if os.path.exists(outpath) else ""
            )
            status = (
                "OK"
                if actual == expected and p.returncode == 0
                else ("WA" if p.returncode == 0 else "RE")
            )
            if p.returncode != 0 and "Traceback" in p.stderr.decode():
                status = "RE"
            cases.append(
                {
                    "caseId": idx,
                    "status": status,
                    "timeMs": elapsed,
                    "stderr": p.stderr.decode(),
                }
            )
            total_time += elapsed
            if status == "OK":
                ok_count += 1
        except subprocess.TimeoutExpired:
            cases.append({"caseId": idx, "status": "TLE", "timeMs": time_limit_ms})
            total_time += time_limit_ms
        except Exception as e:
            cases.append({"caseId": idx, "status": "RE", "timeMs": 0, "stderr": str(e)})
    return cases, total_time, ok_count, len(tests)


if __name__ == "__main__":
    # runner receives payload on stdin for compatibility but arguments also allowed
    try:
        payload = json.load(sys.stdin)
    except Exception:
        payload = {}
    # prefer source_file from payload but also use /code/solution.py by default
    src = payload.get("source_file", "/code/solution.py")
    time_limit = int(payload.get("time_limit_ms", payload.get("time_limit", 1500)))
    # place the source path as sys.argv[1] so subprocess runs proper file
    sys.argv = [sys.argv[0], src]
    cases, total_time, ok_count, total = run_all_tests(time_limit)
    if total == ok_count:
        status = "ACCEPTED"
    elif ok_count == 0:
        status = "WRONG_ANSWER" if any(c["status"] == "WA" for c in cases) else "ERROR"
    else:
        status = "PARTIAL"
    out = {"status": status, "timeMsTotal": total_time, "cases": cases}
    print(json.dumps(out))
