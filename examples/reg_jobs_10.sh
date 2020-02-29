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

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello6 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast6
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello7 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast7
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello8 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast8
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello9 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast9
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello10 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast10

