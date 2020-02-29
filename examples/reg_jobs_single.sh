#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast
