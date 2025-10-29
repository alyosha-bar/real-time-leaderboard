from pydantic import BaseModel
from datetime import datetime

class CodingChallenge(BaseModel):
    id: int
    points: int
    topic: str


class User(BaseModel):
    id: int
    username: str
    email: str


class Submission(BaseModel):
    id: int
    user: User
    challenge: CodingChallenge
    time_to_complete: int
    submitted_at: datetime