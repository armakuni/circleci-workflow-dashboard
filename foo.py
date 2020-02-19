import requests

STATUS_SUCCESS = "success"
STATUS_RUNNING = "running"
STATUS_NOT_RUN = "not_run"
STATUS_FAILED = "failed"
STATUS_ERROR = "error"
STATUS_FAILING = "failing"
STATUS_ON_HOLD = "on_hold"
STATUS_CANCELED = "canceled"
STATUS_UNAUTHORIZED = "unauthorized"
STATUS_UNKNOWN = "unknown"

COMPLETED_STATUSES = [
    STATUS_SUCCESS,
    STATUS_FAILED,
    STATUS_FAILING,
    STATUS_ERROR,
    STATUS_UNKNOWN,
]

api_token = ""
api_response = {}


# Generic calls to each of the APIs
def call_api_v1(api_target):
    response = requests.get(
        f"https://circleci.com/api/v1.1/{api_target}", auth=(api_token, "")
    )
    return response.json()


def call_api_v2(api_target, handle_pagination):
    api_url = f"https://circleci.com/api/v2/{api_target}"
    if not handle_pagination:
        response = requests.get(api_url, auth=(api_token, ""))
        return response.json()
    else:
        all_results = []
        next_page = ""
        while next_page is not None:
            response = requests.get(
                f"{api_url}?page-token={next_page}", auth=(api_token, "")
            )
            response_data = response.json()
            all_results.extend(response_data["items"])
            next_page = response_data["next_page_token"]
        return all_results


# For each project get all pipelines (using project slug)
def get_all_projects():
    return call_api_v1("projects")


def get_all_pipelines(project_slug):
    return call_api_v2(f"project/{project_slug}/pipeline", True)


def filter_pipeline_per_branch(pipelines):
    filtered_pipelines = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"]["branch"]
        if branch not in filtered_pipelines:
            filtered_pipelines[branch] = []
        filtered_pipelines[branch].append(pipeline)
    return filtered_pipelines


def get_latest_pipeline_per_branch(pipelines):
    latest_pipeline_ids = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"]["branch"]
        if branch in latest_pipeline_ids:
            continue
        latest_pipeline_ids[branch] = pipeline["id"]
    return latest_pipeline_ids


def get_workflows_for_pipeline(pipeline_id):
    return call_api_v2(f"pipeline/{pipeline_id}/workflow", True)


def get_jobs_for_workflow(workflow_id):
    return call_api_v2(f"workflow/{workflow_id}/job", True)


def workflow_link(project_slug, pipeline_id, workflow_id):

    #   https://app.circleci.com/github/srbry/hello-world-circle/pipelines/666942a9-41d2-4239-8fc3-babd6e0e4cd0/workflows/fa0d75e8-d594-4e37-832f-28a4627da3bb
    #   https://app.circleci.com/jobs/github/srbry/hello-world-circle/98
    return f"https://app.circleci.com/{project_slug}/pipelines/{pipeline_id}/workflows/{workflow_id}"


def get_previous_completed_status(filtered_pipelines, workflow_name):
    status = STATUS_UNKNOWN
    for pipeline in filtered_pipelines:
        worklow_in_pipeline = False
        for workflow in get_workflows_for_pipeline(pipeline["id"]):
            if workflow["name"] == workflow_name:
                worklow_in_pipeline = True
                if workflow["status"] in COMPLETED_STATUSES:
                    return workflow["status"]
        if not worklow_in_pipeline:
            return status
    return status


def get_dashboard_data():
    # {
    #       "name": "Project X",
    #       "workflow": "Workflow Y",
    #       "branch": "Branch Z",
    #       "status": "success running",
    #   }
    compound_keys = []
    dashboard_data = []
    for project in get_all_projects():
        project_slug = (
            f"{project['vcs_type']}/{project['username']}/{project['reponame']}"
        )
        pipelines = get_all_pipelines(project_slug)
        filtered_pipelines = filter_pipeline_per_branch(pipelines)
        for branch, pipeline_id in get_latest_pipeline_per_branch(pipelines).items():
            # print(f"Branch: {branch}")
            for workflow in get_workflows_for_pipeline(pipeline_id):
                compound_key = f"{project['reponame']}{workflow['name']}-{branch}"
                if compound_key in compound_keys:
                    continue
                status = workflow["status"]
                if workflow["status"] not in COMPLETED_STATUSES:
                    status += f" {get_previous_completed_status(filtered_pipelines[branch], workflow['name'])}"
                dashboard_data.append(
                    {
                        "name": project["reponame"],
                        "workflow": workflow["name"],
                        "branch": branch,
                        "status": status,
                        "link": workflow_link(
                            project_slug, pipeline_id, workflow["id"]
                        ),
                    }
                )
                # print(
                #     f"  Name: {workflow['name']}, Status: {workflow['status']}, ID: {workflow['id']}"
                # )
                compound_keys.append(compound_key)

            # for job in get_jobs_for_workflow(workflow["id"]):
            #     print(f"    Job: {job['name']}, Status: {job['status']}")
    return dashboard_data


# Get latest completed build
# Previous versions with matching workflow name within same pipeline
