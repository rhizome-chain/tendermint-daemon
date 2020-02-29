#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello2 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast2
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello3 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast3
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello4 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast4
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello5 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast5
