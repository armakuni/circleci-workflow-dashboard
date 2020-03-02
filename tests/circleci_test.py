import circleci
import pytest


def test_get_all_projects(requests_mock, circleci_client, projects):
    requests_mock.get(f"{circleci_client.api_url}/api/v1.1/projects", json=projects)
    projects_resp = circleci_client.get_all_projects()
    assert projects_resp == projects


def test_get_all_projects_auth_error(requests_mock, circleci_client, projects):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v1.1/projects", json=projects, status_code=403
    )
    with pytest.raises(circleci.CircleCIAuthError) as err:
        circleci_client.get_all_projects()
    assert (
        str(err.value)
        == "Authentication error hitting projects, do you have permission?"
    )


def test_get_all_projects_error(requests_mock, circleci_client):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v1.1/projects",
        text="Something went wrong",
        status_code=500,
    )
    with pytest.raises(circleci.CircleCIRequestError) as err:
        circleci_client.get_all_projects()
    assert str(err.value) == "Error 500: Something went wrong"


def test_workflow_link(circleci_client, project_slug, pipeline_id, workflow_id):
    assert (
        circleci_client.workflow_link(project_slug, pipeline_id, workflow_id)
        == f"https://app.circleci.com/{project_slug}/pipelines/{pipeline_id}/workflows/{workflow_id}"
    )


def test_job_link(circleci_client, project_slug, job_id):
    assert f"https://app.circleci.com/jobs/{project_slug}/{job_id}"


def test_generate_project_slug(circleci_client, project):
    assert circleci_client.generate_project_slug(project) == "github/foobar/hello-world"


def test_get_all_pipelines_single_page(
    requests_mock, circleci_client, project_slug, pipeline_page_2, pipeline
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/project/{project_slug}/pipeline?page-token=",
        json=pipeline_page_2,
    )
    pipelines = circleci_client.get_all_pipelines(project_slug)
    assert len(pipelines) == 2
    assert pipelines[0] == pipeline


def test_get_all_pipelines_auth_error(
    requests_mock, circleci_client, project_slug, pipeline_page_2, pipeline
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/project/{project_slug}/pipeline?page-token=",
        json=pipeline_page_2,
        status_code=403,
    )
    with pytest.raises(circleci.CircleCIAuthError) as err:
        circleci_client.get_all_pipelines(project_slug)
    assert (
        str(err.value)
        == "Authentication error hitting project/github/foobar/hello-world/pipeline, do you have permission?"
    )


def test_get_all_pipelines_error(
    requests_mock, circleci_client, project_slug, pipeline_page_2, pipeline
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/project/{project_slug}/pipeline?page-token=",
        text="Something went wrong",
        status_code=500,
    )
    with pytest.raises(circleci.CircleCIRequestError) as err:
        circleci_client.get_all_pipelines(project_slug)
    assert str(err.value) == "Error 500: Something went wrong"


def test_get_all_pipelines_multi_page(
    requests_mock,
    circleci_client,
    project_slug,
    pipeline_page_1,
    pipeline_page_2,
    pipeline,
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/project/{project_slug}/pipeline?page-token=",
        json=pipeline_page_1,
    )
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/project/{project_slug}/pipeline?page-token=page-2",
        json=pipeline_page_2,
    )
    pipelines = circleci_client.get_all_pipelines(project_slug)
    assert len(pipelines) == 4
    assert pipelines[0] == pipeline


def test_get_workflows_for_pipeline(
    requests_mock, circleci_client, pipeline_id, workflows_page_1, workflow
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/pipeline/{pipeline_id}/workflow?page-token=",
        json=workflows_page_1,
    )
    workflows = circleci_client.get_workflows_for_pipeline(pipeline_id)
    assert len(workflows) == 2
    assert workflows[0] == workflow


def test_get_jobs_for_workflow(
    requests_mock, circleci_client, workflow_id, jobs_page_1, job
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/workflow/{workflow_id}/job?page-token=",
        json=jobs_page_1,
    )
    jobs = circleci_client.get_jobs_for_workflow(workflow_id)
    assert len(jobs) == 2
    assert jobs[0] == job


def test_get_previous_completed_status(
    requests_mock,
    circleci_client,
    filtered_pipelines,
    workflow_name,
    workflows_page_1,
    pipeline_id,
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/pipeline/{pipeline_id}/workflow?page-token=",
        json=workflows_page_1,
    )
    status = circleci_client.get_previous_completed_status(
        filtered_pipelines["master"], workflow_name
    )
    assert status == circleci.STATUS_SUCCESS


def test_get_previous_completed_status_unknown(
    requests_mock,
    circleci_client,
    filtered_pipelines,
    workflow_name,
    workflows_unknown,
    pipeline_id,
):
    requests_mock.get(
        f"{circleci_client.api_url}/api/v2/pipeline/{pipeline_id}/workflow?page-token=",
        json=workflows_unknown,
    )
    status = circleci_client.get_previous_completed_status(
        filtered_pipelines["master"], workflow_name
    )
    assert status == circleci.STATUS_UNKNOWN


def test_filter_pipeline_per_branch(pipelines, pipeline):
    filtered_pipelines = circleci.filter_pipeline_per_branch(pipelines)
    assert len(filtered_pipelines) == 1
    assert len(filtered_pipelines["master"]) == 2
    assert filtered_pipelines["master"][0] == pipeline


def test_get_latest_pipeline_per_branch(pipelines, pipeline_id):
    latest_pipelines = circleci.get_latest_pipeline_per_branch(pipelines)
    assert len(latest_pipelines) == 1
    assert latest_pipelines["master"] == pipeline_id


def test_workflow_status_completed(circleci_client, workflow2, filtered_pipelines):
    status = circleci_client.workflow_status(workflow2, filtered_pipelines["master"])
    assert status == circleci.STATUS_SUCCESS


def test_workflow_status_not_completed(
    mocker, circleci_client, workflow, filtered_pipelines
):
    mocker.patch.object(
        circleci_client,
        "get_previous_completed_status",
        return_value=circleci.STATUS_SUCCESS,
    )
    status = circleci_client.workflow_status(workflow, filtered_pipelines["master"])
    assert status == f"{circleci.STATUS_CANCELED} {circleci.STATUS_SUCCESS}"
    circleci_client.get_previous_completed_status.assert_called_once()


def test_create_dashboard_monitor(project, workflow):
    dashboard_monitor = circleci.create_dashboard_monitor(
        project, workflow, "master", circleci.STATUS_SUCCESS, "https://foobar.com"
    )
    assert dashboard_monitor == {
        "name": "foobar/hello-world",
        "workflow": "deploy",
        "branch": "master",
        "status": circleci.STATUS_SUCCESS,
        "link": "https://foobar.com",
    }


def test_sort_dashboard_data(dashboard_data, presorted_dashboard_data):
    sorted_dashboard_data = circleci.sort_dashboard_data(dashboard_data)
    assert len(sorted_dashboard_data) == len(dashboard_data)
    assert sorted_dashboard_data == presorted_dashboard_data


def test_get_dashboard_data(
    mocker,
    circleci_client,
    dashboard_projects,
    dashboard_pipeline_side_effect,
    dashboard_workflow_side_effect,
):
    mocker.patch.object(
        circleci_client, "get_all_projects", return_value=dashboard_projects
    )
    mocker.patch.object(
        circleci_client, "get_all_pipelines", side_effect=dashboard_pipeline_side_effect
    )
    mocker.patch.object(
        circleci_client,
        "get_workflows_for_pipeline",
        side_effect=dashboard_workflow_side_effect,
    )
    dashboard_data = circleci.get_dashboard_data(circleci_client)
    assert len(dashboard_data) == 4
    assert dashboard_data == [
        {
            "branch": "dev",
            "link": "https://app.circleci.com/github/foobar/example/pipelines/5/workflows/workflow_6",
            "name": "foobar/example",
            "status": "success",
            "workflow": "test",
        },
        {
            "branch": "master",
            "link": "https://app.circleci.com/github/foobar/example/pipelines/6/workflows/workflow_6",
            "name": "foobar/example",
            "status": "success",
            "workflow": "test",
        },
        {
            "branch": "dev",
            "link": "https://app.circleci.com/github/foobar/hello-world/pipelines/2/workflows/workflow_33",
            "name": "foobar/hello-world",
            "status": "failed",
            "workflow": "deploy",
        },
        {
            "branch": "master",
            "link": "https://app.circleci.com/github/foobar/hello-world/pipelines/4/workflows/workflow_4",
            "name": "foobar/hello-world",
            "status": "canceled success",
            "workflow": "deploy",
        },
    ]


def test_get_dashboard_data(
    mocker,
    circleci_client,
    dashboard_projects,
    dashboard_pipeline_side_effect,
    dashboard_workflow_side_effect,
):
    mocker.patch.object(
        circleci_client, "get_all_projects", return_value=dashboard_projects
    )
    mocker.patch.object(
        circleci_client, "get_all_pipelines", side_effect=dashboard_pipeline_side_effect
    )
    mocker.patch.object(
        circleci_client,
        "get_workflows_for_pipeline",
        side_effect=dashboard_workflow_side_effect,
    )
    dashboard_data = circleci.get_dashboard_data(
        circleci_client, {"foobar/example": None}
    )
    assert len(dashboard_data) == 2
    assert dashboard_data == [
        {
            "branch": "dev",
            "link": "https://app.circleci.com/github/foobar/example/pipelines/5/workflows/workflow_6",
            "name": "foobar/example",
            "status": "success",
            "workflow": "test",
        },
        {
            "branch": "master",
            "link": "https://app.circleci.com/github/foobar/example/pipelines/6/workflows/workflow_6",
            "name": "foobar/example",
            "status": "success",
            "workflow": "test",
        },
    ]
