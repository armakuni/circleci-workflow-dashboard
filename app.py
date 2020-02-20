from flask import Flask, render_template
import os
import circleci

app = Flask(__name__, static_folder="assets", static_url_path="")

api_token = os.environ["CIRCLECI_TOKEN"]
circleci_client = circleci.CircleCI(api_token)


@app.route("/")
def hello_world():
    return render_template(
        "dashboard.html",
        projects=circleci.get_dashboard_data(circleci_client),
        refreshInterval=30,
        now="The time now",
    )
