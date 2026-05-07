import asyncio
import resource
from aiohttp import web

_RUNNER_CODE = """
import signal
import sys
import time

user_code = sys.argv[1]
timeout_ms = max(1, int(sys.argv[2]))

def _timeout_handler(signum, frame):
    raise TimeoutError("Execution timed out")

signal.signal(signal.SIGALRM, _timeout_handler)
signal.setitimer(signal.ITIMER_REAL, timeout_ms / 1000)

start = time.perf_counter()

try:
    exec(compile(user_code, "<submission>", "exec"), {})
except TimeoutError:
    elapsed = int((time.perf_counter() - start) * 1000)
    print("__WORKER_TIMEOUT__", file=sys.stderr)
    print(f"__WORKER_EXEC_MS__={elapsed}", file=sys.stderr)
    sys.exit(124)
except Exception:
    elapsed = int((time.perf_counter() - start) * 1000)
    print(f"__WORKER_EXEC_MS__={elapsed}", file=sys.stderr)
    raise
else:
    elapsed = int((time.perf_counter() - start) * 1000)
    print(f"__WORKER_EXEC_MS__={elapsed}", file=sys.stderr)
"""


def _make_preexec(memory_limit_mb: int):
    def preexec():
        limit = memory_limit_mb * 1024 * 1024
        resource.setrlimit(resource.RLIMIT_AS, (limit, limit))
    return preexec


def _extract_exec_ms(stderr: str, fallback: int) -> int:
    for line in stderr.splitlines():
        if line.startswith("__WORKER_EXEC_MS__="):
            try:
                return int(line.split("=", 1)[1].strip())
            except ValueError:
                pass
    return fallback


async def handle_run(request: web.Request) -> web.Response:
    data = await request.json()
    code = data["code"]
    stdin_input = data.get("input", "") or ""
    time_limit_ms = int(data.get("time_limit_ms", 5000))
    memory_limit_mb = int(data.get("memory_limit_mb", 256))

    infra_timeout = (time_limit_ms / 1000) + 10

    try:
        import sys as _sys
        proc = await asyncio.create_subprocess_exec(
            _sys.executable, "-c", _RUNNER_CODE, code, str(time_limit_ms),
            stdin=asyncio.subprocess.PIPE,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            preexec_fn=_make_preexec(memory_limit_mb),
        )

        try:
            stdout_b, stderr_b = await asyncio.wait_for(
                proc.communicate(stdin_input.encode()),
                timeout=infra_timeout,
            )
        except asyncio.TimeoutError:
            proc.kill()
            await proc.wait()
            return web.json_response({
                "status": "timeout",
                "output": None,
                "error": "Infrastructure timeout",
                "execution_time_ms": time_limit_ms,
            })

        stdout = stdout_b.decode(errors="replace")
        stderr = stderr_b.decode(errors="replace")
        rc = proc.returncode
        exec_ms = _extract_exec_ms(stderr, time_limit_ms)

        if rc == 0:
            return web.json_response({"status": "executed", "output": stdout.strip(), "error": None, "execution_time_ms": exec_ms})
        elif rc == 124 or "__WORKER_TIMEOUT__" in stderr:
            return web.json_response({"status": "timeout", "output": None, "error": "Execution timed out", "execution_time_ms": exec_ms})
        else:
            return web.json_response({"status": "error", "output": None, "error": (stderr or stdout).strip(), "execution_time_ms": exec_ms})

    except Exception as e:
        return web.json_response({"status": "error", "output": None, "error": str(e), "execution_time_ms": 0})


async def handle_health(request: web.Request) -> web.Response:
    return web.json_response({"status": "ok"})


app = web.Application()
app.router.add_post("/run", handle_run)
app.router.add_get("/health", handle_health)

if __name__ == "__main__":
    web.run_app(app, host="0.0.0.0", port=8080)
