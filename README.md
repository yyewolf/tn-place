<p align="center">
  <p align="center">
    <a href="https://github.com/yyewolf/tn-place/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/yyewolf/tn-place.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/yyewolf/tn-place"><img alt="Travis" src="https://travis-ci.org/yyewolf/tn-place.svg?branch=master"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://coveralls.io/github/yyewolf/tn-place?branch=master"><img alt="Coveralls" src="https://coveralls.io/repos/github/yyewolf/tn-place/badge.svg?branch=master"></a>
    <a href="https://codeclimate.com/repos/59bede4e2bfc96025600026b/feed"><img alt="Code Climate Issue Count" src="https://codeclimate.com/repos/59bede4e2bfc96025600026b/badges/d8e88772201d137ea8b7/issue_count.svg"></a>
    <a href="https://goreportcard.com/report/github.com/yyewolf/tn-place"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/yyewolf/tn-place"></a>
    <a href="https://godoc.org/github.com/yyewolf/tn-place"><img src="https://godoc.org/github.com/yyewolf/tn-place?status.svg" alt="GoDoc"></a>
  </p>
</p>

# Golang r/place with OAuth

This project is born to be the clone to rule them all, it's meant to be reusable and fast to deploy.

## Authentication

As for authentication, I went with Google because we have an Education Workspace (rip) at my school.

# Installation

# Prerequisite

- A valid google secret JSON file stored in `back/google.json` for Google Authentication. It needs to have the redirect URI to local as well if you're running this locally.

## Running as-is

To run the project as is, with Docker, simply do the following :

```
touch back/log.txt && touch back/place.png && touch back/place.json
docker compose build
docker compose up
```
