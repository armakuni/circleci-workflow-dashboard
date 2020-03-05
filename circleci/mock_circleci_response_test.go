package circleci_test

const project_resp = `[
	{
		"vcs_type": "github",
		"username": "foobar",
		"reponame": "example",
		"branches": {
			"master": {}
		}
	},
	{
		"vcs_type": "github",
		"username": "fwibble",
		"reponame": "fizzbuzz",
		"branches": {
			"master": {}
		}
	}
]`

const invalid_items_resp = `{
	"next_page_token": null,
	"items": {}
}`

const pipeline_resp_single_page = `{
	"next_page_token": null,
	"items": [
		{
			"id": "1",
			"vcs": {
				"branch": "master"
			}
		},
		{
			"id": "2",
			"vcs": {
				"branch": "develop"
			}
		}
	]
}`

const pipeline_resp_1 = `{
	"next_page_token": "page2",
	"items": [
		{
			"id": "1",
			"vcs": {
				"branch": "master"
			}
		},
		{
			"id": "2",
			"vcs": {
				"branch": "develop"
			}
		}
	]
}`

const pipeline_resp_2 = `{
	"next_page_token": "page3",
	"items": [
		{
			"id": "3",
			"vcs": {
				"branch": "master"
			}
		},
		{
			"id": "4",
			"vcs": {
				"branch": "develop"
			}
		}
	]
}`

const pipeline_resp_3 = `{
	"next_page_token": null
}`

const workflows_resp_single_page = `{
	"next_page_token": null,
	"items": [
		{
			"ID": "1",
			"Name": "workflow1",
			"Status": "success"
		},
		{
			"ID": "2",
			"Name": "workflow2",
			"Status": "failed"
		}
	]
}`
const workflows_resp_1 = `{
	"next_page_token": "page2",
	"items": [
		{
			"ID": "1",
			"Name": "workflow1",
			"Status": "success"
		},
		{
			"ID": "2",
			"Name": "workflow2",
			"Status": "failed"
		}
	]
}`

const workflows_resp_2 = `{
	"next_page_token": "page3",
	"items": [
		{
			"ID": "3",
			"Name": "workflow3",
			"Status": "running"
		},
		{
			"ID": "4",
			"Name": "workflow4",
			"Status": "failed"
		}
	]
}`

const workflows_resp_3 = `{
	"next_page_token": null
}`

const jobs_resp_single_page = `{
	"next_page_token": null,
	"items": [
		{
			"ID": "1"
		},
		{
			"ID": "2"
		}
	]
}`
const jobs_resp_1 = `{
	"next_page_token": "page2",
	"items": [
		{
			"ID": "1"
		},
		{
			"ID": "2"
		}
	]
}`

const jobs_resp_2 = `{
	"next_page_token": "page3",
	"items": [
		{
			"ID": "3"
		},
		{
			"ID": "4"
		}
	]
}`

const jobs_resp_3 = `{
	"next_page_token": null
}`

const workflow_resp_previous_state_1 = `{
	"next_page_token": null,
	"items": [
		{
			"ID": "1",
			"Name": "previous_status",
			"Status": "caceled"
		}
	]
}`

const workflow_resp_previous_state_2 = `{
	"next_page_token": null,
	"items": [
		{
			"ID": "2",
			"Name": "Build Error",
			"Status": "caceled"
		}
	]
}`

const workflow_resp_previous_state_3 = `{
	"next_page_token": null,
	"items": [
		{
			"ID": "3",
			"Name": "previous_status",
			"Status": "success"
		}
	]
}`
