# CircleCI Dashboard

Gives a build radiator view to a CircleCI instance.

The built in tools for CircleCI are okay if you want to see the history of everything ever. but that's not helpful for reacting to issues and failures.

This dashboard is designed to be displayed loud and proud to show everyone what urgent issues have been found in the builds. The quick-look nature enables teams to focus on the important stuff quickly.

## Installation

This project uses [Poetry](https://python-poetry.org/) to manage dependencies. Make sure this is [installed](https://python-poetry.org/docs/#installation)

To install the dependencies for this project, run poetry
```bash
$ poetry install
```

## Setup

Running this dashboard requires an API token from CircleCI. Create a new one in your [User Account](https://circleci.com/account/api)

Export this as an environment variable before starting the dasboard

```bash
$ export CIRCLECI_TOKEN=<Personal API Token>
```

## Running

With the dependencies installed and the API Token available, start the dashboard via poetry

```bash
$ poetry run python app.py
```
