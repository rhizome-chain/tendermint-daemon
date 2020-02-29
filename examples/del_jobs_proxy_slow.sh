#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_slow
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/hello_slow_proxy



