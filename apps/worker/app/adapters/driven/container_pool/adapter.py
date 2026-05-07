import asyncio
from typing import Optional
import aiohttp

from app.adapters.driven.container_pool.runner_slot import RunnerSlot
from app.adapters.driven.container_pool.slot_state import SlotState
from app.application.ports.executor import Executor
from app.domain.models import SubmissionResult


class ContainerPoolAdapter(Executor):

    def __init__(
        self,
        language: str,
        image: str,
        count: int,
        network: str,
        name_prefix: str,
        health_check_timeout: float,
        restart_timeout: float,
        repair_interval: float,
    ):
        self.language = language
        self.image = image
        self.count = count
        self.network = network
        self.name_prefix = name_prefix
        self.health_check_timeout = health_check_timeout
        self.restart_timeout = restart_timeout
        self.repair_interval = repair_interval

        self._all_slots: list[RunnerSlot] = []
        self._queue: asyncio.Queue[RunnerSlot] = asyncio.Queue()
        self._repair_task: Optional[asyncio.Task] = None

    async def startup(self) -> None:
        print(f"[Pool:{self.language}] Starting {self.count} runner containers", flush=True)
        for i in range(self.count):
            name = f"{self.name_prefix}-{self.language}-{i}"
            slot = RunnerSlot(container_name=name, image=self.image, network=self.network)
            self._all_slots.append(slot)
            asyncio.create_task(self._boot_slot(slot))

        self._repair_task = asyncio.create_task(self._repair_loop())

    async def shutdown(self) -> None:
        if self._repair_task:
            self._repair_task.cancel()
            try:
                await self._repair_task
            except asyncio.CancelledError:
                pass

        await asyncio.gather(
            *[self._kill_container(s.container_name) for s in self._all_slots],
            return_exceptions=True,
        )
        print(f"[Pool:{self.language}] All containers stopped", flush=True)

    async def acquire(self) -> RunnerSlot:
        while True:
            slot = await self._queue.get()

            if slot.state == SlotState.DEAD:
                continue

            if await self._is_healthy(slot):
                slot.set_state(SlotState.BUSY)
                return slot

            print(f"[Pool:{self.language}] {slot.container_name} failed health check on acquire — restarting", flush=True)
            asyncio.create_task(self._restart_slot(slot))

    async def release(self, slot: RunnerSlot, crashed: bool = False) -> None:
        if crashed:
            print(f"[Pool:{self.language}] {slot.container_name} crashed during execution — restarting", flush=True)
            asyncio.create_task(self._restart_slot(slot))
        else:
            slot.set_state(SlotState.AVAILABLE)
            await self._queue.put(slot)

    async def execute(self, submission: SubmissionResult) -> SubmissionResult:
        slot = await self.acquire()
        crashed = False

        try:
            payload = {
                "code": submission.code,
                "input": submission.input,
                "time_limit_ms": submission.time_limit_ms,
                "memory_limit_mb": submission.memory_limit_mb,
            }
            timeout = aiohttp.ClientTimeout(total=(submission.time_limit_ms / 1000) + 15)

            async with aiohttp.ClientSession() as session:
                async with session.post(f"{slot.url}/run", json=payload, timeout=timeout) as resp:
                    if resp.status != 200:
                        crashed = True
                        raise RuntimeError(f"Runner returned HTTP {resp.status}")
                    data = await resp.json()

            submission.status = data["status"]
            submission.output = data.get("output")
            submission.error = data.get("error")
            submission.execution_time_ms = data.get("execution_time_ms", 0)

        except Exception as e:
            crashed = True
            submission.status = "error"
            submission.output = None
            submission.error = str(e)
            submission.execution_time_ms = 0

        finally:
            await self.release(slot, crashed=crashed)

        return submission

    # ------------------------------------------------------------------
    # Container lifecycle
    # ------------------------------------------------------------------

    async def _boot_slot(self, slot: RunnerSlot) -> None:
        slot.set_state(SlotState.RESTARTING)
        await self._kill_container(slot.container_name)

        proc = await asyncio.create_subprocess_exec(
            "docker", "run", "-d",
            "--name", slot.container_name,
            "--network", self.network,
            slot.image,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        _, stderr_b = await proc.communicate()

        if proc.returncode != 0:
            err = stderr_b.decode(errors="replace").strip()
            print(f"[Pool:{self.language}] Failed to start {slot.container_name}: {err}", flush=True)
            await self._handle_boot_failure(slot)
            return

        slot.url = await self._resolve_url(slot.container_name)
        if not slot.url:
            print(f"[Pool:{self.language}] Could not resolve IP for {slot.container_name}", flush=True)
            await self._handle_boot_failure(slot)
            return

        # Poll until healthy (up to 60 attempts × 1s)
        for _ in range(60):
            await asyncio.sleep(1)
            if await self._is_healthy(slot):
                slot.restart_count = 0
                slot.set_state(SlotState.AVAILABLE)
                await self._queue.put(slot)
                print(f"[Pool:{self.language}] {slot.container_name} ready at {slot.url}", flush=True)
                return

        print(f"[Pool:{self.language}] {slot.container_name} never became healthy", flush=True)
        await self._handle_boot_failure(slot)

    async def _restart_slot(self, slot: RunnerSlot) -> None:
        slot.restart_count += 1
        if slot.is_dead():
            slot.set_state(SlotState.DEAD)
            print(f"[Pool:{self.language}] {slot.container_name} marked DEAD after {slot.restart_count} failures", flush=True)
            return

        print(f"[Pool:{self.language}] Restarting {slot.container_name} (attempt {slot.restart_count})", flush=True)
        await self._boot_slot(slot)

    async def _handle_boot_failure(self, slot: RunnerSlot) -> None:
        slot.restart_count += 1
        if slot.is_dead():
            slot.set_state(SlotState.DEAD)
            print(f"[Pool:{self.language}] {slot.container_name} marked DEAD after {slot.restart_count} failures", flush=True)
        else:
            slot.set_state(SlotState.RESTARTING)

    async def _kill_container(self, name: str) -> None:
        proc = await asyncio.create_subprocess_exec(
            "docker", "rm", "-f", name,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.DEVNULL,
        )
        await proc.wait()

    # ------------------------------------------------------------------
    # Health / introspection
    # ------------------------------------------------------------------

    async def _is_healthy(self, slot: RunnerSlot) -> bool:
        if not slot.url:
            return False
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(
                    f"{slot.url}/health",
                    timeout=aiohttp.ClientTimeout(total=self.health_check_timeout),
                ) as resp:
                    return resp.status == 200
        except Exception:
            return False

    async def _resolve_url(self, container_name: str) -> str:
        proc = await asyncio.create_subprocess_exec(
            "docker", "inspect",
            "--format", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
            container_name,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.DEVNULL,
        )
        stdout, _ = await proc.communicate()
        ip = stdout.decode().strip()
        return f"http://{ip}:8080" if ip else ""

    # ------------------------------------------------------------------
    # Background repair loop
    # ------------------------------------------------------------------

    async def _repair_loop(self) -> None:
        while True:
            await asyncio.sleep(self.repair_interval)
            for slot in self._all_slots:
                if slot.state != SlotState.RESTARTING:
                    continue
                if slot.seconds_in_state > self.restart_timeout:
                    print(
                        f"[Pool:{self.language}] {slot.container_name} stuck in RESTARTING "
                        f"for {slot.seconds_in_state:.0f}s (limit {self.restart_timeout}s) — forcing restart",
                        flush=True,
                    )
                    asyncio.create_task(self._restart_slot(slot))
