BASE_DIR=/etc/cochain/eosTransfer
COMMON_NAME=127.0.0.1:8123 # Caution: grpc client will validate this common name

mkdir -p $BASE_DIR/certs/
openssl genrsa -out ${BASE_DIR}/certs/server.key 2048
openssl req -new -x509 -key ${BASE_DIR}/certs/server.key -out ${BASE_DIR}/certs/server.pem -days 3650 -subj "/CN=${COMMON_NAME}/O=Cochain Technology/C=CN"
