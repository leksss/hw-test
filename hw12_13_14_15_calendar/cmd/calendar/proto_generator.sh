#!/bin/bash

rm -rf ../../pb/event
mkdir -p ../../pb/event
protoc -I ../../api/googleapis -I ../../api/event \
  --proto_path=../../api/event \
  --go_out=../../pb/event \
  --go-grpc_out=../../pb/event \
  --grpc-gateway_out=../../pb/event \
  ../../api/event/*.proto
