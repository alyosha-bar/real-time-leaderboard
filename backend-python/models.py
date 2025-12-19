from pydantic import BaseModel
from datetime import datetime
from typing import Dict

class User(BaseModel):
    username: str

class Challenge(BaseModel):
    id: int
    points: int

class Submission(BaseModel):
    user: User
    challenge: Challenge