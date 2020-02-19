from flask import Flask, render_template

import foo

app = Flask(__name__, static_folder="assets", static_url_path="")


@app.route("/")
def hello_world():
    return render_template(
        "dashboard.html",
        projects=foo.get_dashboard_data(),
        refreshInterval=30,
        now="The time now",
    )
