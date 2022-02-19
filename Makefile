export GO111MODULE=on

ARMCC:=/usr/local/gcc-linaro-aarch64-linux-gnu/bin/aarch64-linux-gnu-gcc

all:
	go build -o service-iot main.go

init:
	go env -w GOPROXY=https://goproxy.cn,direct
	go mod tidy
	go mod vendor

arm64:
	# cd aarch64/ && sudo bash build.sh build
	# sudo chmod a+r main
	LD_PRELOAD="" \
	CC=${ARMCC}  \
	GOOS=linux  \
	GOARCH=arm64  \
	GOARM=7  \
	CGO_ENABLED=1 \
	go build -o service-iot main.go 

arm64_m2:
	LD_PRELOAD="" \
	CC=${ARMCC}  \
	GOOS=linux  \
	GOARCH=arm64  \
	GOARM=7  \
	CGO_ENABLED=1 \
	go build -o service-iot main.go 

r:all
	./service-iot