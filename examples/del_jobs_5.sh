#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast1
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast2
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast3
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast4
curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast5
