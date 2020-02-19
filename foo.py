import requests

api_token = ""

response = requests.get("https://circleci.com/api/v1.1/projects", auth=(api_token, ""))

projects = response.json()

# For each project get all pipelines (using project slug)


def get_pipelines(project_slug):
    pipelines = []
    next_page = ""
    while next_page is not None:
        response = requests.get(
            f"https://circleci.com/api/v2/project/{project_slug}/pipeline?page-token={next_page}",
            auth=(api_token, ""),
        )
        pipeline_info = response.json()
        pipelines.extend(pipeline_info["items"])
        next_page = pipeline_info["next_page_token"]
    return pipelines


def get_latest_pipeline_per_branch(pipelines):
    latest_pipeline_ids = {}
    for pipeline in pipelines:
        branch = pipeline["vcs"]["branch"]
        if branch in latest_pipeline_ids:
            continue
        latest_pipeline_ids[branch] = pipeline["id"]
    return latest_pipeline_ids


def get_workflows_for_pipeline(pipeline_id):
    workflows = []
    next_page = ""
    while next_page is not None:
        response = requests.get(
            f"https://circleci.com/api/v2/pipeline/{pipeline_id}/workflow",
            auth=(api_token, ""),
        )
        pipeline_info = response.json()
        workflows.extend(pipeline_info["items"])
        next_page = pipeline_info["next_page_token"]
    return workflows


def get_workflow_jobs(workflow_id):
    jobs = []
    next_page = ""
    while next_page is not None:
        response = requests.get(
            f"https://circleci.com/api/v2/workflow/{workflow_id}/job",
            auth=(api_token, ""),
        )
        pipeline_info = response.json()
        jobs.extend(pipeline_info["items"])
        next_page = pipeline_info["next_page_token"]
    return jobs


for project in projects:
    project_slug = f"{project['vcs_type']}/{project['username']}/{project['reponame']}"
    pipelines = get_pipelines(project_slug)
    print(len(pipelines))
    print(get_latest_pipeline_per_branch(pipelines))
    for branch, id in get_latest_pipeline_per_branch(pipelines).items():
        workflows = get_workflows_for_pipeline(id)
        print(branch)
        for workflow in workflows:
            print(f"  Name: {workflow['name']}, Status: {workflow['status']}")
            print(get_workflow_jobs(workflow["id"]))

# Get latest completed build
