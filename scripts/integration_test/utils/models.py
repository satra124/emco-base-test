from enum import Enum
from typing import Dict, List, Optional

from pydantic import UUID4, BaseModel, Field, validator

class CommandOutput(BaseModel):
    return_code: int
    output: str

class AutomationType(BaseModel):
    type: str
    description: str
    total_test_cases: int