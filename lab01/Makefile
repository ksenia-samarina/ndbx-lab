CUR_DIR=$(shell pwd)
BIN_DIR=${CUR_DIR}/bin
PACKAGE=samarina/ndbx/cmd/app

bin-dir:
	mkdir -p ${BIN_DIR}

build: bin-dir
	go build -o ${BIN_DIR}/app ${PACKAGE}