#!/bin/bash

set -e

curl -X POST https://api.github.com/repos/chyroc/rss/dispatches \
-H 'Accept: application/vnd.github.everest-preview+json' \
-H "Authorization: token $DISPATCH_ACTION_TOKEN" \
--data '{"event_type": "opened"}'