import asyncio
import textwrap
from app.application.ports.executor import Executor
from app.domain.models import SubmissionResult
from config import DOCKER_IMAGE, EXECUTION_TIMEOUT


class DockerExecutor(Executor):

    _RUNNER_CODE = textwrap.dedent(
        """
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
    )

    async def ensure_image_cached(self) -> None:
        proc = await asyncio.create_subprocess_exec(
            "docker", "image", "inspect", DOCKER_IMAGE,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.DEVNULL,
        )
        await proc.wait()

        if proc.returncode == 0:
            return

        proc = await asyncio.create_subprocess_exec(
            "docker", "pull", DOCKER_IMAGE,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.DEVNULL,
        )
        await proc.wait()
        if proc.returncode != 0:
            raise RuntimeError(f"Failed to pull Docker image {DOCKER_IMAGE}")

    async def execute(self, submission: SubmissionResult) -> SubmissionResult:
        cmd = [
            "docker", "run", "--rm", "-i",
            "--network", "none",
            "--memory", f"{submission.memory_limit_mb}m",
            DOCKER_IMAGE,
            "python", "-c", self._RUNNER_CODE,
            submission.code,
            str(submission.time_limit_ms),
        ]

        try:
            # Do not count container startup overhead as execution timeout.
            # Timeout inside the container is enforced by _RUNNER_CODE.
            infra_timeout_s = max(EXECUTION_TIMEOUT, (submission.time_limit_ms / 1000) + 30)

            proc = await asyncio.create_subprocess_exec(
                *cmd,
                stdin=asyncio.subprocess.PIPE,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
            )

            try:
                stdout_bytes, stderr_bytes = await asyncio.wait_for(
                    proc.communicate((submission.input or "").encode()),
                    timeout=infra_timeout_s,
                )
            except asyncio.TimeoutError:
                proc.kill()
                await proc.wait()
                submission.status = "timeout"
                submission.output = None
                submission.error = "Infrastructure timeout while starting container"
                submission.execution_time_ms = submission.time_limit_ms
                return submission

            stdout = stdout_bytes.decode()
            stderr = stderr_bytes.decode()
            returncode = proc.returncode

            execution_time = self._extract_exec_time_ms(stderr, submission.time_limit_ms)

            if returncode == 0:
                submission.status = "executed"
                submission.output = stdout.strip()
                submission.execution_time_ms = execution_time
                submission.error = None
            elif self._is_timeout_result(returncode, stderr):
                submission.status = "timeout"
                submission.output = None
                submission.error = "Execution timed out"
                submission.execution_time_ms = execution_time
            else:
                submission.status = "error"
                submission.output = None
                submission.execution_time_ms = execution_time
                submission.error = (stderr or stdout).strip() or "Execution failed"

            return submission

        except Exception as e:
            submission.status = "error"
            submission.output = None
            submission.error = str(e)
            submission.execution_time_ms = 0
            return submission

    def _extract_exec_time_ms(self, stderr: str, fallback: int) -> int:
        for line in stderr.splitlines():
            if line.startswith("__WORKER_EXEC_MS__="):
                raw = line.split("=", 1)[1].strip()
                try:
                    return int(raw)
                except ValueError:
                    return fallback
        return fallback

    def _is_timeout_result(self, returncode: int, stderr: str) -> bool:
        if returncode == 124:
            return True
        return "__WORKER_TIMEOUT__" in (stderr or "")
