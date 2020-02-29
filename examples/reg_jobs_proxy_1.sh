#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello1" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy1