#!/bin/bash

rm -rf ../../proto/protobuf
mkdir -p ../../proto/protobuf
protoc -I ../../proto/googleapis -I ../../proto/event \
  --proto_path=../../proto/event \
  --go_out=../../proto/protobuf \
  --go-grpc_out=../../proto/protobuf \
  --grpc-gateway_out=../../proto/protobuf \
  ../../proto/event/*.proto
