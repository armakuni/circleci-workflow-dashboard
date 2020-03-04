from twisted.web import server, resource, static
from twisted.internet import reactor, endpoints, task
from datetime import datetime, timezone
import jinja2
import os
import circleci
import json


port = os.getenv("PORT", "5000")
circleci_api_url = os.getenv("CIRCLECI_API_URL", circleci.DEFAULT_API_URL)
circleci_jobs_url = os.getenv("CIRCLECI_JOBS_URL", circleci.DEFAULT_JOBS_URL)
circleci_dashboard_filter = json.loads(os.getenv("DASHBOARD_FILTER", "null"))
refresh_interval = int(os.getenv("REFERSH_INTERVAL", "30"))

api_token = os.getenv("CIRCLECI_TOKEN", None)
if api_token is None:
    print("You must provide a CIRCLECI_TOKEN as an environment variable")
    exit(1)

circleci_client = circleci.CircleCI(api_token, circleci_api_url, circleci_jobs_url)


def current_time():
    now = datetime.now(timezone.utc)
    return now.strftime("%y-%m-%d %H:%M:%S %Z")


class Dashboard(resource.Resource):
    projects = {}
    current_time = ""

    def update_projects(self):
        self.projects = circleci.get_dashboard_data(
            circleci_client, circleci_dashboard_filter
        )
        self.current_time = current_time()

    def render_GET(self, request):
        templateLoader = jinja2.FileSystemLoader(searchpath="./templates")
        templateEnv = jinja2.Environment(loader=templateLoader)
        template = templateEnv.get_template("dashboard.html")

        return template.render(
            projects=self.projects,
            refreshInterval=refresh_interval,
            now=self.current_time,
        ).encode("utf-8")


dashboard = Dashboard()

root = resource.Resource()
root.putChild(b"assets", static.File("./assets"))
root.putChild(b"", dashboard)
site = server.Site(root)

project_updater = task.LoopingCall(dashboard.update_projects)
project_updater.start(refresh_interval)

print(f"Starting server on http://127.0.0.1:{port}")
reactor.listenTCP(int(port), site)
reactor.run()
