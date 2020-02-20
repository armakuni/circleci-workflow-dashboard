from flask import Flask, render_template
from datetime import datetime, timezone
import os
import circleci

app = Flask(__name__, static_folder="assets", static_url_path="")

api_token = os.environ["CIRCLECI_TOKEN"]
circleci_client = circleci.CircleCI(api_token)


def current_time():
    now = datetime.now(timezone.utc)
    return now.strftime("%y-%m-%d %H:%M:%S %Z")


@app.route("/")
def hello_world():
    return render_template(
        "dashboard.html",
        projects=circleci.get_dashboard_data(circleci_client),
        refreshInterval=30,
        now=current_time(),
    )
