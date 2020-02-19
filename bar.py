from flask import Flask, render_template

import foo

app = Flask(__name__, static_folder="assets", static_url_path="")


# mock_status = [
#     {
#         "name": "Project X",
#         "workflow": "Workflow Y",
#         "branch": "Branch Z",
#         "status": "success running",
#     },
#     {
#         "name": "Project A",
#         "workflow": "Workflow B",
#         "branch": "Branch C",
#         "status": "success",
#     },
#     {
#         "name": "Project D",
#         "workflow": "Workflow E",
#         "branch": "Branch F",
#         "status": "failed",
#     },
#     {
#         "name": "Project G",
#         "workflow": "Workflow H",
#         "branch": "Branch I",
#         "status": "failed onhold",
#     },
#     {
#         "name": "Project J",
#         "workflow": "Workflow K",
#         "branch": "Branch L",
#         "status": "success cancelled",
#     },
#     {
#         "name": "Project Foo",
#         "workflow": "Workflow Bar",
#         "branch": "Branch Baz",
#         "status": "unknown running",
#     },
# ]


@app.route("/")
def hello_world():
    return render_template(
        "dashboard.html",
        projects=foo.get_dashboard_data(),
        refreshInterval=30,
        now="The time now",
    )
