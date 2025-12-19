from locust import HttpUser, task, between
import random

class LeaderboardUser(HttpUser):
    host = "http://localhost:8000"
    wait_time = between(1, 2)

    @task
    def submit_score(self):
        username = f"user_{random.randint(1, 1000)}"
        challenge_id = random.randint(1, 10)
        points = random.randint(90, 250)

        submission_data = {
            "user": {
                "username": username
            },
            "challenge": {
                "id": challenge_id,
                "points": points
            }
        }

        self.client.post("/submit", json=submission_data)