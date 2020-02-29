#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"0.1s","greet":"hello1 0.1s" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello_slow
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello_slow" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_slow_proxy