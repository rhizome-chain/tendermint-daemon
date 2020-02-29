#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_proxy1



