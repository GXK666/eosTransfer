# Download and install protoc
PROTOC_VERSION=3.5.1 # you should replace it with the new version number
wget https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip
unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip bin/protoc -d /usr/local

# Install plugins or tools
go get -v github.com/gogo/protobuf/protoc-gen-gogoslick
go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -v github.com/square/goprotowrap/cmd/protowrap
