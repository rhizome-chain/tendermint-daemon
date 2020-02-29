#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy1


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello2
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy2


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello3
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy3


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello4
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy4


curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello5
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy5



