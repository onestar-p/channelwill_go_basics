function genProto {
    PROTO_NAME=$1
    PROTO_PATH=./${PROTO_NAME}
    if [ $2 ]; then
        GO_OUT_PATH="${PROTO_PATH}/gen/${2}"
    else
        GO_OUT_PATH=${PROTO_PATH}/gen
    fi
    mkdir -p $GO_OUT_PATH
    protoc -I=${PROTO_PATH} --go_out=plugins=grpc,paths=source_relative:$GO_OUT_PATH ${PROTO_NAME}.proto
    protoc -I=${PROTO_PATH} --grpc-gateway_out=paths=source_relative,grpc_api_configuration=${PROTO_PATH}/${PROTO_NAME}.yaml:$GO_OUT_PATH ${PROTO_NAME}.proto
}

# genProto service v1
genProto etranslate v1
genProto auth v1
