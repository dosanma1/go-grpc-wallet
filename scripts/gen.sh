#!/bin/bash

protoc --proto_path=./pkg/proto --go_out=paths=source_relative:pkg/pb/ --go-grpc_out=paths=source_relative:pkg/pb/ pkg/proto/*.proto