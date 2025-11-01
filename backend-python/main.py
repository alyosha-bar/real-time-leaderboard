from fastapi import FastAPI, BackgroundTasks
from fastapi.encoders import jsonable_encoder
import httpx
import asyncio
from typing import List, Dict
from models import Submission
from models import Analytics
from datetime import datetime
import statistics
import logging
import json

app = FastAPI(title="Real-Time Leaderboard Analytics")


# global variables to replace databases for development purposes
submissions: List[Submission] = []

analytics : Analytics = Analytics(
    total_submissions=0,
    avg_completion=0,
    topics={
        "Sorting Algorithms": 0,
        "REST API Design": 0,
        "Concurrency": 0,
        "Data Structures": 0,
        "Graph Theory": 0,
        "Dynamic Programming": 0,
        "File I/O": 0,
        "Unit Testing": 0,
    }
)


@app.get("/")
def root():
    return {"message": "Real-Time Leaderboard Analytics API running."}

@app.post("/submit")
def receive_submission(sub: Submission, background_tasks: BackgroundTasks):

    global analytics

    submissions.append(sub)
    print(f"Received submission: {sub}")

    # update analytics based on submissions
    total_submissions = len(submissions)
    
    # generate mapping
    topics = {
        "Sorting Algorithms": 0,
        "REST API Design": 0,
        "Concurrency": 0,
        "Data Structures": 0,
        "Graph Theory": 0,
        "Dynamic Programming": 0,
        "File I/O": 0,
        "Unit Testing": 0,
    }

    # avg time to complete
    total_time = 0


    for sub in submissions:
        total_time += sub.time_to_complete
        # increment dict where key is sub.topic
        topic = sub.challenge.topic
        topics[topic] = topics.get(topic, 0) + 1
    
    avg_completion = int(total_time / total_submissions) # sum of submission times / total_submissions



    analytics = Analytics(
        total_submissions = total_submissions,
        avg_completion = avg_completion,
        topics = topics
    )



    score = get_score(sub)
    print(f"Score for submission: {score}")
    background_tasks.add_task(push_score, score)

    return {"status": "Submission received"}


@app.get("/submissions", response_model=List[Submission])
async def get_submissions():

    serialized_submissions = serialize_submissions(submissions)
    print(f"Serialized submissions: {serialized_submissions}")

    return serialized_submissions

@app.get("/analytics", response_model=Analytics)
def get_analytics():

    if not analytics:
        return {"message": "No analytics yet."}
    

    return jsonable_encoder(analytics)


def get_score(sub: Submission):
    return {
        "username": sub.user.username,
        "points": sub.challenge.points,
    }


async def push_score(score):
    try:
        async with httpx.AsyncClient() as client:
            response = await client.post("http://backend-go:8080/score", json=score)
            response.raise_for_status()
    except Exception as e:
        logging.error(f"Failed to push score: {e}")

