# whip-go

whip-go is a simple WHIP (https://datatracker.ietf.org/doc/draft-ietf-wish-whip/) client implementation in go using the WebRTC [Pion libraries](https://github.com/pion).
It includes a WHIPClient class and a simple command line client supporrting screensharing to a WHIP ingestion endpoint.

It has been tested with [janus-gateway](https://github.com/meetecho/janus-gateway) with [whip-simple-server](https://github.com/lminiero/simple-whip-server) and also with mediasoup, millicast and cloudflare implementations.

## Installation

```
go build
```

## Running

```
./whip-go -v VIDEO_SOURCE -a AUDIO_SOURCE -vc VIDEO_CODEC -t TOKEN WHIP_ENDPOINT_URL
```

The supported sources are either "screen" for screensharing, "test" for a test patter or the name of a file (f.e. "/dev/stdin") to ready raw YUV420 samples from.

The supported video codecs are vp8 and h264.

For more information and additional configuration run:

```
./whip-go -h
```

## Cross-Compiling

We will use the docker build for this. Getting `zig` to work as a cross compiler is too complicated (codec libraries & `pkg-config`) and not worth the effort.

```shell
docker build --tag whip-go .
docker buildx build --platform linux/amd64 --tag whip-go .
docker run --rm -it whip-go sh

# in another terminal
docker cp __container_id__:/opt/whip-go/whip-go ./whip-go.linux

# stop container & clean up
docker rmi whip-go
```

```shell
# zig approach, which we abandoned
PKG_CONFIG_PATH="$HOMEBREW_PREFIX/lib/pkgconfig" \
CGO_ENABLED=1 \
GOOS=linux GOARCH=amd64 \
CC="zig cc -target x86_64-linux-gnu" \
CXX="zig c++ -target x86_64-linux-gnu" \
go build -o whip-go.linux
```
