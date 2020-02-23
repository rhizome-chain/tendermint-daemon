#!/bin/bash

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/veryfast

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X DELETE http://localhost:7777/v1/daemon/job/fast

curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X DELETE http://localhost:7777/v1/daemon/job/hi1
curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X DELETE http://localhost:7777/v1/daemon/job/hi2
curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X DELETE http://localhost:7777/v1/daemon/job/hi3
curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X DELETE http://localhost:7777/v1/daemon/job/hi4
curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X DELETE http://localhost:7777/v1/daemon/job/hi5
