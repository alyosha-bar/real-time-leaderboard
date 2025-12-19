from fastapi import FastAPI
from typing import Dict
from models import Submission

import redis

app = FastAPI(title="Real-Time Leaderboard Analytics")


pool = redis.ConnectionPool(host='localhost', port=6379, db=0, decode_responses=True)
r = redis.Redis(connection_pool=pool)

# Basic root endpoint
@app.get("/")
def root():
    return {"message": "Real-Time Leaderboard Analytics API running."}

# Endpoint for receiving and validating submissions
@app.post("/submit", response_model=Dict[str, str])
def submit_submission(submission: Submission):
    # extract user.username and challenge.points
    username = submission.user.username
    points = submission.challenge.points

    # push to redis
    push_to_redis(username, points)

    # publish event to notify Go service
    publish_event()

    return {"status": "success", "message": "Submission received and processed."}




# function for pushing to REDIS
def push_to_redis(username: str, points: int):

    print("Pushing to Redis:", username, points)

    r.zincrby("leaderboard", points, username)
    return


# publish an event (Go will subscribe)    
def publish_event():
    r.publish("score_updates", "new_score_recorded")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="127.0.0.1", port=8000, reload=True)







