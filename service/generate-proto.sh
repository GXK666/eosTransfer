#!/bin/bash

set -x

PROTOS="./*/*.proto"

INCLUDES="\
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/gogo/protobuf/protobuf \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  -I.
"

# https://github.com/golang/protobuf/issues/39
# https://github.com/square/goprotowrap
# TODO: protowrap cannot notify any plugin error output
protowrap --print_structure \
  $INCLUDES \
  --gogoslick_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. \
  --grpc-gateway_out=logtostderr=true,request_context=true:. \
  --swagger_out=logtostderr=true:. \
  $PROTOS

# protoc \
#  $INCLUDES \
#  --gogoslick_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. \
#  --grpc-gateway_out=logtostderr=true,request_context=true:. \
#  --swagger_out=allow_merge=true,merge_file_name=rpc,logtostderr=true:. \
#  $PROTOS
