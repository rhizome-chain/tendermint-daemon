#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello1" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy1


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello2
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello3" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy2


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello3
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello3" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy3


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello4
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello4" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy4


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"interval":"1ms","greet":"hello1 1ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hello5
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data '{"source":"hello5" }' http://localhost:7777/v1/daemon/job/add/factory/hello-proxy-worker/jobid/hello_proxy5

