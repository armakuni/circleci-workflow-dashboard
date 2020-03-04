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

DEFAULT_API_URL = "https://circleci.com"
DEFAULT_JOBS_URL = "https://app.circleci.com"


class CircleCI:
    def __init__(self, api_token, api_url=DEFAULT_API_URL, job_url=DEFAULT_JOBS_URL):
        self.api_token = api_token
        self.api_url = api_url
        self.job_url = job_url

    # Generic calls to each of the APIs (v1.1, v2)
    # Input: string
    # Output: Array{}
    def __call_api_v1(self, api_target):
        response = requests.get(
            f"{self.api_url}/api/v1.1/{api_target}", auth=(self.api_token, "")
        )
        if not response.ok:
            if response.status_code in [401, 403]:
                raise CircleCIAuthError(
                    f"Authentication error hitting {api_target}, do you have permission?"
                )
            raise CircleCIRequestError(f"Error {response.status_code}: {response.text}")
        return response.json()

    # v2 seems to have introduced pagination.
    # Take the stress of following the tokens out of the business logic
    # Input: string, boolean
    # Output: Array{}
    def __call_api_v2(self, api_target):
        api_url = f"{self.api_url}/api/v2/{api_target}"
        all_results = []
        next_page = ""
        while next_page is not None:

            response = requests.get(
                f"{api_url}?page-token={next_page}", auth=(self.api_token, "")
            )
            if not response.ok:
                if response.status_code in [401, 403]:
                    raise CircleCIAuthError(
                        f"Authentication error hitting {api_target}, do you have permission?"
                    )
                raise CircleCIRequestError(
                    f"Error {response.status_code}: {response.text}"
                )
            response_data = response.json()
            all_results.extend(response_data["items"])
            next_page = response_data["next_page_token"]
        return all_results

    # Get all the projects for the authenticated user (API_TOKEN)
    # Input: None
    # Output: Array{}
    def get_all_projects(self):
        return self.__call_api_v1("projects")

    # Get all the Pipelines under a given project_slug
    # Input: string
    # Output: Array{}
    def get_all_pipelines(self, project_slug):
        return self.__call_api_v2(f"project/{project_slug}/pipeline")

    # Get all the Workflows used by a Pipeline
    # Input: string
    # Output: Array{}
    def get_workflows_for_pipeline(self, pipeline_id):
        return self.__call_api_v2(f"pipeline/{pipeline_id}/workflow")

    # Get all the Jobs run by a Workflow
    # Input: string
    # Output: Array{}
    def get_jobs_for_workflow(self, workflow_id):
        return self.__call_api_v2(f"workflow/{workflow_id}/job")

    # Iterate through previous Pipeline runs to find the last completed run for a Workflow
    # Input: string, string
    # Output: string
    def get_previous_completed_status(self, filtered_pipelines, workflow_name):
        status = STATUS_UNKNOWN
        for pipeline in filtered_pipelines:
            workflow_in_pipeline = False
            for workflow in self.get_workflows_for_pipeline(pipeline["id"]):
                if workflow["name"] == "Build Error":
                    workflow_in_pipeline = True
                if workflow["name"] == workflow_name:
                    workflow_in_pipeline = True
                    if workflow["status"] in COMPLETED_STATUSES:
                        return workflow["status"]
            if not workflow_in_pipeline:
                return status
        return status

    # Generate a URI to the targeted Pipeline run
    # Input: string, string, string
    # Output: string
    def workflow_link(self, project_slug, pipeline_id, workflow_id):
        return f"{self.job_url}/{project_slug}/pipelines/{pipeline_id}/workflows/{workflow_id}"

    # Generate a URI to the targeted Job
    # Input: string, string
    # Output: string
    def job_link(self, project_slug, job_id):
        return f"{self.job_url}/jobs/{project_slug}/{job_id}"

    def generate_project_slug(self, project):
        # Generate the 'project_slug' for each one in the format :vcs_type/:vcs_username/:vcs_project_name
        return f"{project['vcs_type']}/{project['username']}/{project['reponame']}"

    def workflow_status(self, workflow, branch_pipelines):
        # For "running state" (states with animated borders), we need to append the correct secondary state
        # e.g "success running" will be green with a blue animated border
        status = workflow["status"]
        if workflow["status"] not in COMPLETED_STATUSES:
            status += f" {self.get_previous_completed_status(branch_pipelines, workflow['name'])}"
        return status


# Get a dictionary of branches with the associated Pipeline objects
# Input: string
# Output: Array{}
def filter_pipeline_per_branch(pipelines):
    filtered_pipelines = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"].get("branch", None)
        if branch is None:
            continue
        if branch not in filtered_pipelines:
            filtered_pipelines[branch] = []
        filtered_pipelines[branch].append(pipeline)
    return filtered_pipelines


# Get a dictionary of branches with the associated Pipeline IDs
# Input: string
# Output: Array{}
def get_latest_pipeline_per_branch(pipelines):
    latest_pipeline_ids = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"].get("branch", None)
        if branch in latest_pipeline_ids or branch is None:
            continue
        latest_pipeline_ids[branch] = pipeline["id"]
    return latest_pipeline_ids


def create_dashboard_monitor(project, workflow, branch, status, link):
    return {
        "name": _project_name(project),
        "workflow": workflow["name"],
        "branch": branch,
        "status": status,
        "link": link,
    }


def get_dashboard_data(circleci_client, project_filters=None):
    compound_keys = []
    dashboard_data = []
    # Iterate all projects followed in CircleCI
    for project in circleci_client.get_all_projects():
        if (
            project_filters is not None
            and _project_name(project) not in project_filters
        ):
            continue
        project_slug = circleci_client.generate_project_slug(project)
        # Get the pipelines (runs) associated with the project
        pipelines = circleci_client.get_all_pipelines(project_slug)
        # Filter the pipelines down into the VCS branches they ran for
        filtered_pipelines = filter_pipeline_per_branch(pipelines)
        # Iterate over the branches in the project, focussing on the last pipeline run of each
        for branch, pipeline_id in get_latest_pipeline_per_branch(pipelines).items():
            # Iterate over a list of the workflows used by the lastest pipeline run (normally only one)
            for workflow in circleci_client.get_workflows_for_pipeline(pipeline_id):
                # Generate a compound key to identify the pipeline uniquely
                compound_key = f"{project['username']}/{project['reponame']}{workflow['name']}-{branch}"
                # Check to see if that key already exists, if it does - ignore this workflow / pipeline run
                if compound_key in compound_keys:
                    continue
                status = circleci_client.workflow_status(
                    workflow, filtered_pipelines[branch]
                )
                link = circleci_client.workflow_link(
                    project_slug, pipeline_id, workflow["id"]
                )
                dashboard_data.append(
                    create_dashboard_monitor(project, workflow, branch, status, link)
                )
                # Append the compound key to show we've processed this project/branch/workflow
                compound_keys.append(compound_key)
    return sort_dashboard_data(dashboard_data)


def sort_dashboard_data(dashboard_data):
    return sorted(
        dashboard_data, key=lambda i: f"{i['name']}-{i['workflow']}-{i['branch']}"
    )


def _project_name(project):
    return f"{project['username']}/{project['reponame']}"


class CircleCIAuthError(Exception):
    pass


class CircleCIRequestError(Exception):
    pass
