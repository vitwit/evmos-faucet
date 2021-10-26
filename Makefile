
clean:
	rm -rf build

build:
	make clean
	mkdir -p build
	go build -ldflags="-X 'github.com/cosmos/cosmos-sdk/version.Name=sim' -X 'github.com/cosmos/cosmos-sdk/version.AppName=simd'" \
		-o ./build/faucet main.go