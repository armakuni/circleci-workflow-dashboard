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


class CircleCI:
    def __init__(
        self,
        api_token,
        api_url="https://circleci.com",
        job_url="https://app.circleci.com",
    ):
        self.api_token = api_token
        self.api__url = api_url
        self.job_url = job_url

    # Generic calls to each of the APIs (v1.1, v2)
    # Input: string
    # Output: Array{}
    def call_api_v1(self, api_target):
        response = requests.get(
            f"{self.api__url}/api/v1.1/{api_target}", auth=(self.api_token, "")
        )
        return response.json()

    # v2 seems to have introduced pagination.
    # Take the stress of following the tokens out of the business logic
    # Input: string, boolean
    # Output: Array{}
    def call_api_v2(self, api_target, handle_pagination):
        api_url = f"{self.api__url}/api/v2/{api_target}"
        if not handle_pagination:
            response = requests.get(api_url, auth=(self.api_token, ""))
            return response.json()
        else:
            all_results = []
            next_page = ""
            while next_page is not None:
                response = requests.get(
                    f"{api_url}?page-token={next_page}", auth=(self.api_token, "")
                )
                response_data = response.json()
                all_results.extend(response_data["items"])
                next_page = response_data["next_page_token"]
            return all_results

    # Get all the projects for the authenticated user (API_TOKEN)
    # Input: None
    # Output: Array{}
    def get_all_projects(self):
        return self.call_api_v1("projects")

    # Get all the Pipelines under a given project_slug
    # Input: string
    # Output: Array{}
    def get_all_pipelines(self, project_slug):
        return self.call_api_v2(f"project/{project_slug}/pipeline", True)

    # Get all the Workflows used by a Pipeline
    # Input: string
    # Output: Array{}
    def get_workflows_for_pipeline(self, pipeline_id):
        return self.call_api_v2(f"pipeline/{pipeline_id}/workflow", True)

    # Get all the Jobs run by a Workflow
    # Input: string
    # Output: Array{}
    def get_jobs_for_workflow(self, workflow_id):
        return self.call_api_v2(f"workflow/{workflow_id}/job", True)

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

    # Iterate through previous Pipeline runs to find the last completed run for a Workflow
    # Input: string, string
    # Output: string
    def get_previous_completed_status(self, filtered_pipelines, workflow_name):
        status = STATUS_UNKNOWN
        for pipeline in filtered_pipelines:
            workflow_in_pipeline = False
            for workflow in self.get_workflows_for_pipeline(pipeline["id"]):
                if workflow["name"] == workflow_name:
                    workflow_in_pipeline = True
                    if workflow["status"] in COMPLETED_STATUSES:
                        return workflow["status"]
            if not workflow_in_pipeline:
                return status
        return status


# Get a dictionary of branches with the associated Pipeline objects
# Input: string
# Output: Array{}
def filter_pipeline_per_branch(pipelines):
    filtered_pipelines = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"]["branch"]
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
        branch = pipeline["vcs"]["branch"]
        if branch in latest_pipeline_ids:
            continue
        latest_pipeline_ids[branch] = pipeline["id"]
    return latest_pipeline_ids


def get_dashboard_data(circleci_client):
    compound_keys = []
    dashboard_data = []
    # Iterate all projects followed in CircleCI
    for project in circleci_client.get_all_projects():
        # Generate the 'project_slug' for each one in the format :vcs_type/:vcs_username/:vcs_project_name
        project_slug = (
            f"{project['vcs_type']}/{project['username']}/{project['reponame']}"
        )
        # Get the pipelines (runs) associated with the project
        pipelines = circleci_client.get_all_pipelines(project_slug)
        # Filter the pipelines down into the VCS branches they ran for
        filtered_pipelines = filter_pipeline_per_branch(pipelines)
        # Iterate over the branches in the project, focussing on the last pipeline run of each
        for branch, pipeline_id in get_latest_pipeline_per_branch(pipelines).items():
            # Iterate over a list of the workflows used by the lastest pipeline run (normally only one)
            for workflow in circleci_client.get_workflows_for_pipeline(pipeline_id):
                # Generate a compound key to identify the pipeline uniquely
                compound_key = f"{project['reponame']}{workflow['name']}-{branch}"
                # Check to see if that key already exists, if it does - ignore this workflow / pipeline run
                if compound_key in compound_keys:
                    continue
                # For "running state" (states with animated borders), we need to append the correct secondary state
                status = workflow["status"]
                if workflow["status"] not in COMPLETED_STATUSES:
                    status += f" {circleci_client.get_previous_completed_status(filtered_pipelines[branch], workflow['name'])}"
                # Finally, put all this working out together, append that to the returned data, and start all over again
                dashboard_data.append(
                    {
                        "name": project["reponame"],
                        "workflow": workflow["name"],
                        "branch": branch,
                        "status": status,
                        "link": circleci_client.workflow_link(
                            project_slug, pipeline_id, workflow["id"]
                        ),
                    }
                )
                # Append the compound key to show we've processed this project/branch/workflow
                compound_keys.append(compound_key)
    return dashboard_data
