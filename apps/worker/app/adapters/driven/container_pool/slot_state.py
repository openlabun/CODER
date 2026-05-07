from enum import Enum

class SlotState(Enum):
    AVAILABLE = "available"
    BUSY = "busy"
    RESTARTING = "restarting"
    DEAD = "dead"