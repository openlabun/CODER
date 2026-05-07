import asyncio
import os
import resource
import tempfile
import time
from aiohttp import web


def _make_preexec(memory_limit_mb: int):
    def preexec():
        limit = memory_limit_mb * 1024 * 1024
        resource.setrlimit(resource.RLIMIT_AS, (limit, limit))
    return preexec


async def handle_run(request: web.Request) -> web.Response:
    data = await request.json()
    code = data["code"]
    stdin_input = data.get("input", "") or ""
    time_limit_ms = int(data.get("time_limit_ms", 5000))
    memory_limit_mb = int(data.get("memory_limit_mb", 256))

    with tempfile.TemporaryDirectory() as tmpdir:
        src_path = os.path.join(tmpdir, "solution.cpp")
        bin_path = os.path.join(tmpdir, "solution")

        with open(src_path, "w") as f:
            f.write(code)

        compile_proc = await asyncio.create_subprocess_exec(
            "g++", "-O2", "-o", bin_path, src_path,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        try:
            _, compile_err_b = await asyncio.wait_for(compile_proc.communicate(), timeout=30)
        except asyncio.TimeoutError:
            compile_proc.kill()
            await compile_proc.wait()
            return web.json_response({"status": "error", "output": None, "error": "Compilation timed out", "execution_time_ms": 0})

        if compile_proc.returncode != 0:
            return web.json_response({
                "status": "error",
                "output": None,
                "error": compile_err_b.decode(errors="replace").strip(),
                "execution_time_ms": 0,
            })

        time_limit_s = time_limit_ms / 1000
        infra_timeout = time_limit_s + 10
        start = time.perf_counter()

        try:
            run_proc = await asyncio.create_subprocess_exec(
                "timeout", "--signal=SIGKILL", str(time_limit_s), bin_path,
                stdin=asyncio.subprocess.PIPE,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
                preexec_fn=_make_preexec(memory_limit_mb),
            )
            try:
                stdout_b, stderr_b = await asyncio.wait_for(
                    run_proc.communicate(stdin_input.encode()),
                    timeout=infra_timeout,
                )
            except asyncio.TimeoutError:
                run_proc.kill()
                await run_proc.wait()
                return web.json_response({
                    "status": "timeout",
                    "output": None,
                    "error": "Infrastructure timeout",
                    "execution_time_ms": time_limit_ms,
                })

            elapsed_ms = int((time.perf_counter() - start) * 1000)
            stdout = stdout_b.decode(errors="replace")
            stderr = stderr_b.decode(errors="replace")
            rc = run_proc.returncode

            if rc == 0:
                return web.json_response({"status": "executed", "output": stdout.strip(), "error": None, "execution_time_ms": elapsed_ms})
            elif rc == 124:
                return web.json_response({"status": "timeout", "output": None, "error": "Execution timed out", "execution_time_ms": time_limit_ms})
            else:
                return web.json_response({"status": "error", "output": None, "error": (stderr or stdout).strip() or "Runtime error", "execution_time_ms": elapsed_ms})

        except Exception as e:
            return web.json_response({"status": "error", "output": None, "error": str(e), "execution_time_ms": 0})


async def handle_health(request: web.Request) -> web.Response:
    return web.json_response({"status": "ok"})


app = web.Application()
app.router.add_post("/run", handle_run)
app.router.add_get("/health", handle_health)

if __name__ == "__main__":
    web.run_app(app, host="0.0.0.0", port=8080)
