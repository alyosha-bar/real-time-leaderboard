from fastapi import FastAPI, BackgroundTasks
from fastapi.encoders import jsonable_encoder
import httpx
import asyncio
from typing import List
from models import Submission
from datetime import datetime
import statistics
import logging
import json

app = FastAPI(title="Real-Time Leaderboard Analytics")

submissions: List[Submission] = []


@app.get("/")
def root():
    return {"message": "Real-Time Leaderboard Analytics API running."}

@app.post("/submit")
def receive_submission(sub: Submission, background_tasks: BackgroundTasks):
    submissions.append(sub)
    print(f"Received submission: {sub}")
    # analytics = get_analytics()
    # print(f"Current analytics: {analytics}")
    # background_tasks.add_task(push_analytics, analytics)

    score = get_score(sub)
    print(f"Score for submission: {score}")
    background_tasks.add_task(push_score, score)

    return {"status": "Submission received"}


@app.get("/submissions", response_model=List[Submission])
async def get_submissions():

    serialized_submissions = serialize_submissions(submissions)
    print(f"Serialized submissions: {serialized_submissions}")

    return serialized_submissions

def get_analytics():

    if not submissions:
        return {"message": "No submissions yet."}

    # avg_time = round(statistics.mean(sub.time_to_complete for sub in submissions), 2)

    # user_scores = {}
    # for s in submissions:
    #     user_scores[s.user.username] = user_scores.get(s.user.username, 0) + s.challenge.points

    # top_user = max(user_scores, key=user_scores.get)

    # return {
    #     "total_submissions": len(submissions),
    #     "average_time_to_complete": avg_time,
    #     "leaderboard": user_scores,
    #     "top_user": top_user
    # }

    return jsonable_encoder(submissions)


def get_score(sub: Submission):
    return {
        "username": sub.user.username,
        "points": sub.challenge.points,
    }

async def push_analytics(analytics):
    try:
        async with httpx.AsyncClient() as client:
            response = await client.post("http://backend-go:8080/analytics", json=analytics)
            response.raise_for_status()
    except Exception as e:
        logging.error(f"Failed to push analytics: {e}")


async def push_score(score):
    try:
        async with httpx.AsyncClient() as client:
            response = await client.post("http://backend-go:8080/score", json=score)
            response.raise_for_status()
    except Exception as e:
        logging.error(f"Failed to push score: {e}")

