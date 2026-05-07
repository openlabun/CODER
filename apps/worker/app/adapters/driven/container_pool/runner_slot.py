from datetime import datetime, timezone
from app.adapters.driven.container_pool.slot_state import SlotState

class RunnerSlot:
    MAX_RESTARTS = 3

    def __init__(self, container_name: str, image: str, network: str):
        self.container_name = container_name
        self.image = image
        self.network = network
        self.url: str = ""
        self.state: SlotState = SlotState.RESTARTING
        self.state_changed_at: datetime = datetime.now(timezone.utc)
        self.restart_count: int = 0

    def set_state(self, state: SlotState) -> None:
        self.state = state
        self.state_changed_at = datetime.now(timezone.utc)

    @property
    def seconds_in_state(self) -> float:
        return (datetime.now(timezone.utc) - self.state_changed_at).total_seconds()

    def is_dead(self) -> bool:
        return self.restart_count >= self.MAX_RESTARTS