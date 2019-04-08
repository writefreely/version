# version

Version is a simple API for finding the latest release version of a project on GitHub. It's built for projects that are distributed through GitHub's _release_ functionality, and meant to be run on a project's own domain, e.g. [version.writefreely.org](https://version.writefreely.org).

## Features

* Lightweight wrapper on top of GitHub's API
* Resilient / still functional when GitHub is down
* Useful public API for varying use cases
* Minimal logging

## Getting Started

With Go installed, run these commands:

```text
go get github.com/writefreely/version/cmd/version

export VER_ORG=writeas
export VER_REPO=writefreely
export VER_PORT=8080
version
```

Then open your browser to http://localhost:8080

## API

### GET /

Returns version number / tag as plain text, e.g. `v0.8.1`

#### Parameters

| Parameter | Type | Use |
| --- | --- | --- |
| `v` | String | Supply a version number to get a response of whether or not it's the current one |
