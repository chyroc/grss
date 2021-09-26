#!/bin/bash

set -e

curl -X POST https://api.github.com/repos/chyroc/rss/dispatches \
-H 'Accept: application/vnd.github.everest-preview+json' \
-H "Authorization: token $GITHUB_TOKEN" \
--data '{"event_type": "opened"}'