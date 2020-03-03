from pytest import fixture
from circleci import CircleCI


@fixture
def circleci_client():
    return CircleCI("foo")


@fixture
def project():
    return {
        "irc_server": None,
        "ssh_keys": [],
        "branches": {
            "master": {
                "latest_workflows": {
                    "test": {
                        "status": "success",
                        "created_at": "2020-02-19T10: 49: 27.966Z",
                        "id": "1b1e7a9f-eaca-4a37-8366-4d6db2a48a93",
                    }
                },
                "pusher_logins": ["foobar"],
                "running_builds": [],
                "recent_builds": [
                    {
                        "outcome": "success",
                        "status": "success",
                        "build_num": 1,
                        "vcs_revision": "b02ff726d3c43552ee3ae1f62e000b86da7e99fa",
                        "pushed_at": "2020-02-19T10: 49: 27.000Z",
                        "is_workflow_job": True,
                        "is_2_0_job": True,
                        "added_at": "2020-02-19T10: 52: 44.156Z",
                    }
                ],
            },
            "last_success": {
                "outcome": "success",
                "status": "success",
                "build_num": 1,
                "vcs_revision": "b02ff726d3c43552ee3ae1f62e000b86da7e99fa",
                "pushed_at": "2020-02-19T10: 49: 27.000Z",
                "is_workflow_job": True,
                "is_2_0_job": True,
                "added_at": "2020-02-19T10: 52: 44.156Z",
            },
            "is_using_workflows": True,
        },
        "irc_keyword": None,
        "oss": False,
        "slack_channel": None,
        "reponame": "hello-world",
        "dependencies": "",
        "aws": {"keypair": None},
        "slack_webhook_url": None,
        "irc_channel": None,
        "parallel": 1,
        "slack_integration_access_token": None,
        "username": "foobar",
        "slack_integration_team": None,
        "slack_integration_channel": None,
        "heroku_deploy_user": None,
        "irc_username": None,
        "slack_notify_prefs": None,
        "slack_subdomain": None,
        "has_usable_key": True,
        "setup": "",
        "vcs_type": "github",
        "feature_flags": {
            "trusty-beta": False,
            "builds-service": True,
            "osx": False,
            "set-github-status": True,
            "build-prs-only": False,
            "forks-receive-secret-env-vars": True,
            "fleet": None,
            "build-fork-prs": False,
            "autocancel-builds": False,
        },
        "irc_password": None,
        "compile": "",
        "slack_integration_notify_prefs": None,
        "slack_integration_webhook_url": None,
        "irc_notify_prefs": None,
        "slack_integration_team_id": None,
        "extra": "",
        "jira": None,
        "slack_integration_channel_id": None,
        "language": "Go",
        "flowdock_api_token": None,
        "slack_channel_override": None,
        "vcs_url": "https://github.com/foobar/hello-world",
        "following": True,
        "default_branch": "master",
        "slack_api_token": None,
        "test": "",
    }


@fixture
def projects(project):
    return [project]


@fixture
def project_slug(circleci_client, project):
    return circleci_client.generate_project_slug(project)


@fixture
def pipeline_id():
    return "b8a20fa9-1e61-4bcf-9353-de9a812377e6"


@fixture
def prev_pipeline_id():
    return "afc20fa9-1e61-4bcf-9353-de9a812374b9"


@fixture
def prev_pipeline_id_2():
    return "b1d20fa9-1e61-4bcf-9353-de9a81237c36"


@fixture
def pipeline(pipeline_id):
    return {
        "id": pipeline_id,
        "errors": [],
        "project_slug": "gh/foobar/hello-world",
        "updated_at": "2020-02-19T10:19:23.471Z",
        "number": 18,
        "state": "created",
        "created_at": "2020-02-19T10:19:23.471Z",
        "trigger": {
            "received_at": "2020-02-19T10:19:23.436Z",
            "type": "webhook",
            "actor": {
                "login": "foobar",
                "avatar_url": "https://foobar.com/u/123456?v=4",
            },
        },
        "vcs": {
            "origin_repository_url": "https://github.com/foobar/hello-world",
            "target_repository_url": "https://github.com/foobar/hello-world",
            "revision": "5db088471dd324e94cf6f7af8084d2ebd7109f69",
            "provider_name": "GitHub",
            "branch": "master",
        },
    }


@fixture
def pipeline_prev(prev_pipeline_id):
    return {
        "id": prev_pipeline_id,
        "errors": [],
        "project_slug": "gh/foobar/hello-world",
        "updated_at": "2020-02-19T10:19:23.471Z",
        "number": 18,
        "state": "created",
        "created_at": "2020-02-19T10:19:23.471Z",
        "trigger": {
            "received_at": "2020-02-19T10:19:23.436Z",
            "type": "webhook",
            "actor": {
                "login": "foobar",
                "avatar_url": "https://foobar.com/u/123456?v=4",
            },
        },
        "vcs": {
            "origin_repository_url": "https://github.com/foobar/hello-world",
            "target_repository_url": "https://github.com/foobar/hello-world",
            "revision": "5db088471dd324e94cf6f7af8084d2ebd7109f69",
            "provider_name": "GitHub",
            "branch": "master",
        },
    }


@fixture
def pipeline_prev2(prev_pipeline_id_2):
    return {
        "id": prev_pipeline_id_2,
        "errors": [],
        "project_slug": "gh/foobar/hello-world",
        "updated_at": "2020-02-19T10:19:23.471Z",
        "number": 18,
        "state": "created",
        "created_at": "2020-02-19T10:19:23.471Z",
        "trigger": {
            "received_at": "2020-02-19T10:19:23.436Z",
            "type": "webhook",
            "actor": {
                "login": "foobar",
                "avatar_url": "https://foobar.com/u/123456?v=4",
            },
        },
        "vcs": {
            "origin_repository_url": "https://github.com/foobar/hello-world",
            "target_repository_url": "https://github.com/foobar/hello-world",
            "revision": "5db088471dd324e94cf6f7af8084d2ebd7109f69",
            "provider_name": "GitHub",
            "branch": "master",
        },
    }


@fixture
def pipeline_broken():
    return {
        "id": "54f18b24-866d-4798-8d86-451dfbb95a35",
        "errors": [
            {"type": "config", "message": "|   |     null"},
            {"type": "config", "message": "|   |   INPUT:"},
            {"type": "config", "message": "|   |     type: object"},
            {"type": "config", "message": "|   |   SCHEMA:"},
            {
                "type": "config",
                "message": "|   2. [#/jobs/lint/steps/1] expected type: Mapping, found: String",
            },
            {"type": "config", "message": "|   |     - add_ssh_keys"},
            {"type": "config", "message": "|   |     - setup_remote_docker"},
            {"type": "config", "message": "|   |     - checkout"},
            {"type": "config", "message": "|   |     enum:"},
            {
                "type": "config",
                "message": "|   |   Steps without arguments can be called as strings",
            },
            {
                "type": "config",
                "message": "|   1. [#/jobs/lint/steps/1] Input not a valid enum value",
            },
            {
                "type": "config",
                "message": "1. [#/jobs/lint/steps/1] 0 subschemas matched instead of one",
            },
            {
                "type": "config",
                "message": "[#/jobs/lint] only 1 subschema matches out of 2",
            },
            {"type": "config", "message": "ERROR IN CONFIG FILE:"},
        ],
        "project_slug": "gh/foobar/hello-world",
        "updated_at": "2020-02-18T14:48:47.776Z",
        "number": 13,
        "state": "errored",
        "created_at": "2020-02-18T14:48:47.776Z",
        "trigger": {
            "received_at": "2020-02-18T14:48:47.742Z",
            "type": "webhook",
            "actor": {
                "login": "foobar",
                "avatar_url": "https://foobar.com/u/123456?v=4",
            },
        },
        "vcs": {
            "origin_repository_url": "https://github.com/foobar/hello-world",
            "target_repository_url": "https://github.com/foobar/hello-world",
            "revision": "208f4fa5c3d5a624621ce743066d1fa9f7861157",
            "provider_name": "GitHub",
            "commit": {"body": "", "subject": "Update config.yml"},
            "branch": "master",
        },
    }


@fixture
def pipelines(pipeline, pipeline_broken):
    return [pipeline, pipeline_broken]


@fixture
def pipeline_page_2(pipelines):
    return {"next_page_token": None, "items": pipelines}


@fixture
def pipeline_page_1(pipelines):
    return {"next_page_token": "page-2", "items": pipelines}


@fixture
def workflow_id():
    return "7d1c893f-982f-45a3-9aec-627345cece6d"


@fixture
def workflow_name():
    return "deploy"


@fixture
def workflow(workflow_id, workflow_name, pipeline_id):
    return {
        "stopped_at": "2020-02-19T16:31:51Z",
        "name": workflow_name,
        "project_slug": "gh/foobar/hello-world",
        "pipeline_number": 28,
        "status": "canceled",
        "id": "7d1c893f-982f-45a3-9aec-627345cece6d",
        "created_at": "2020-02-19T16:29:52Z",
        "pipeline_id": pipeline_id,
    }


@fixture
def workflow2(pipeline_id):
    return {
        "stopped_at": "2020-02-19T15:52:04Z",
        "name": "Build Error",
        "project_slug": "gh/foobar/hello-world",
        "pipeline_number": 28,
        "status": "errored",
        "id": "fc17ac52-9558-4aad-8d98-1234c48f50as",
        "created_at": "2020-02-19T15:48:49Z",
        "pipeline_id": pipeline_id,
    }


@fixture
def workflow3(workflow_name, pipeline_id):
    return {
        "stopped_at": "2020-02-19T15:52:04Z",
        "name": workflow_name,
        "project_slug": "gh/foobar/hello-world",
        "pipeline_number": 28,
        "status": "success",
        "id": "fc17ac52-9558-4aad-8d98-1234c48f5055",
        "created_at": "2020-02-19T15:48:49Z",
        "pipeline_id": pipeline_id,
    }


@fixture
def workflow4(workflow_name, prev_pipeline_id_2):
    return {
        "stopped_at": "2020-02-19T15:52:04Z",
        "name": workflow_name,
        "project_slug": "gh/foobar/hello-world",
        "pipeline_number": 28,
        "status": "success",
        "id": "fc17ac52-9558-4aad-8d98-1234c48f5055",
        "created_at": "2020-02-19T15:48:49Z",
        "pipeline_id": prev_pipeline_id_2,
    }


@fixture
def workflows_page_1(workflow, workflow3):
    return {"next_page_token": None, "items": [workflow, workflow3]}


@fixture
def workflows_page_prev_1(workflow):
    return {"next_page_token": None, "items": [workflow]}


@fixture
def workflows_page_prev_2(workflow2):
    return {"next_page_token": None, "items": [workflow2]}


@fixture
def workflows_page_prev_3(workflow4):
    return {"next_page_token": None, "items": [workflow4]}


@fixture
def workflows_unknown(workflow):
    return {"next_page_token": None, "items": [workflow]}


@fixture
def job_id():
    return "a60fb799-b010-4fec-a899-6325de9d21b3"


@fixture
def job():
    return {
        "dependencies": [],
        "job_number": 111,
        "id": "a60fb799-b010-4fec-a899-6325de9d21b3",
        "started_at": "2020-02-19T16:30:02Z",
        "name": "test",
        "project_slug": "gh/foobar/hello-world",
        "status": "success",
        "type": "build",
        "stopped_at": "2020-02-19T16:30:36Z",
    }


@fixture
def job2():
    return {
        "dependencies": ["a60fb799-b010-4fec-a899-6325de9d21b3"],
        "job_number": 112,
        "id": "c819c8ab-a43c-43da-ae9c-74266ad74ba1",
        "started_at": "2020-02-19T16:30:41Z",
        "name": "build",
        "project_slug": "gh/foobar/hello-world",
        "status": "success",
        "type": "build",
        "stopped_at": "2020-02-19T16:31:20Z",
    }


@fixture
def jobs_page_1(job, job2):
    return {"next_page_token": None, "items": [job, job2]}


@fixture
def filtered_pipelines(pipeline):
    return {"master": [pipeline]}


@fixture
def filtered_pipelines_prev(pipeline, pipeline_prev, pipeline_prev2):
    return {"master": [pipeline, pipeline_prev, pipeline_prev2]}


@fixture
def dashboard_data():
    return [
        {
            "name": "zzz",
            "workflow": "foo",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "foo",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "foo",
            "branch": "active",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "bar",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "aaaa",
            "workflow": "bar",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
    ]


@fixture
def presorted_dashboard_data():
    return [
        {
            "name": "aaaa",
            "workflow": "bar",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "bar",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "foo",
            "branch": "active",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "abc/def",
            "workflow": "foo",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
        {
            "name": "zzz",
            "workflow": "foo",
            "branch": "master",
            "status": "success",
            "link": "https://foobar.com",
        },
    ]


@fixture
def dashboard_projects():
    return [
        {"reponame": "hello-world", "username": "foobar", "vcs_type": "github"},
        {"reponame": "example", "username": "foobar", "vcs_type": "github"},
        {"reponame": "example", "username": "foobar", "vcs_type": "github"},
    ]


@fixture
def dashboard_pipelines_hello_world():
    return [
        {"id": "4", "vcs": {"branch": "master"}},
        {"id": "3", "vcs": {"branch": "master"}},
        {"id": "2", "vcs": {"branch": "dev"}},
        {"id": "1", "vcs": {"branch": "master"}},
    ]


@fixture
def dashboard_pipelines_example():
    return [
        {"id": "6", "vcs": {"branch": "master"}},
        {"id": "5", "vcs": {"branch": "dev"}},
    ]


@fixture
def dashboard_workflows_4():
    return [
        {
            "name": "deploy",
            "status": "canceled",
            "id": "workflow_4",
            "pipeline_id": "4",
        },
        {"name": "deploy", "status": "success", "id": "workflow_3", "pipeline_id": "4"},
    ]


@fixture
def dashboard_workflows_3():
    return [
        {"name": "deploy", "status": "failed", "id": "workflow_33", "pipeline_id": "3"},
        {"name": "deploy", "status": "failed", "id": "workflow_32", "pipeline_id": "3"},
    ]


@fixture
def dashboard_workflows_6():
    return [
        {"name": "test", "status": "success", "id": "workflow_6", "pipeline_id": "6"}
    ]


@fixture
def dashboard_workflows_5():
    return [
        {"name": "test", "status": "success", "id": "workflow_6", "pipeline_id": "5"}
    ]


@fixture
def dashboard_pipeline_side_effect(
    dashboard_pipelines_hello_world, dashboard_pipelines_example
):
    def impl(project_slug):
        if project_slug == "github/foobar/hello-world":
            return dashboard_pipelines_hello_world
        return dashboard_pipelines_example

    return impl


@fixture
def dashboard_workflow_side_effect(
    dashboard_workflows_6,
    dashboard_workflows_5,
    dashboard_workflows_4,
    dashboard_workflows_3,
):
    def impl(pipeline_id):
        switcher = {
            "6": dashboard_workflows_6,
            "5": dashboard_workflows_5,
            "4": dashboard_workflows_4,
            "3": dashboard_workflows_3,
        }
        return switcher.get(pipeline_id, dashboard_workflows_3)

    return impl
