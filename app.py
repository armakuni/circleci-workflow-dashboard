from twisted.web import server, resource, static
from twisted.internet import reactor, endpoints, task
from datetime import datetime, timezone
import jinja2
import os
import circleci


REFRESH_INTERVAL = 30

port = os.getenv("PORT", "5000")
api_token = os.environ["CIRCLECI_TOKEN"]
circleci_client = circleci.CircleCI(api_token)


def current_time():
    now = datetime.now(timezone.utc)
    return now.strftime("%y-%m-%d %H:%M:%S %Z")


class Dashboard(resource.Resource):
    projects = {}
    current_time = ""

    def update_projects(self):
        self.projects = circleci.get_dashboard_data(circleci_client)
        self.current_time = current_time()

    def render_GET(self, request):
        templateLoader = jinja2.FileSystemLoader(searchpath="./templates")
        templateEnv = jinja2.Environment(loader=templateLoader)
        template = templateEnv.get_template("dashboard.html")

        return template.render(
            projects=self.projects,
            refreshInterval=REFRESH_INTERVAL,
            now=self.current_time,
        ).encode("utf-8")


dashboard = Dashboard()

root = resource.Resource()
root.putChild(b"assets", static.File("./assets"))
root.putChild(b"", dashboard)
site = server.Site(root)

project_updater = task.LoopingCall(dashboard.update_projects)
project_updater.start(REFRESH_INTERVAL)

print(f"Starting server on http://127.0.0.1:{port}")
reactor.listenTCP(int(port), site)
reactor.run()
