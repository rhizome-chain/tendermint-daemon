#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"10ms","greet":"hello 10ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/fast

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"100ms","greet":"hello 100ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi1

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"300ms","greet":"hello 300ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi2

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"0.5s","greet":"hello 0.5s" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi3

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"0.8s","greet":"hello 0.8s" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi4

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1.5s","greet":"hello 1.5s" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi5
