from flask import Flask, render_template
from datetime import datetime, timezone
import threading
import atexit
import os
import circleci

REFRESH_INTERVAL = 30
projects = {}

# lock to control access to variable
circleLock = threading.Lock()
# thread handler
fetchCircleDataThread = threading.Thread()


def create_app(circleci_client):
    app = Flask(__name__, static_folder="assets", static_url_path="")

    def interrupt():
        global fetchCircleDataThread
        fetchCircleDataThread.cancel()

    def update_projects():
        global projects
        global fetchCircleDataThread
        with circleLock:
            projects = circleci.get_dashboard_data(circleci_client)

        fetchCircleDataThread = threading.Timer(REFRESH_INTERVAL, update_projects, ())
        fetchCircleDataThread.start()

    def init_projects():
        global fetchCircleDataThread
        global projects
        projects = circleci.get_dashboard_data(circleci_client)
        fetchCircleDataThread = threading.Timer(REFRESH_INTERVAL, update_projects, ())
        fetchCircleDataThread.start()

    # Initiate
    init_projects()
    # When you kill Flask (SIGTERM), clear the trigger for the next thread
    atexit.register(interrupt)
    return app


def current_time():
    now = datetime.now(timezone.utc)
    return now.strftime("%y-%m-%d %H:%M:%S %Z")


api_token = os.environ["CIRCLECI_TOKEN"]
circleci_client = circleci.CircleCI(api_token)
app = create_app(circleci_client)


@app.route("/")
def hello_world():
    global projects
    return render_template(
        "dashboard.html",
        projects=projects,
        refreshInterval=REFRESH_INTERVAL,
        now=current_time(),
    )
