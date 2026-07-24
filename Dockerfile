# Stage 1: Compile tdlib

FROM alpine:latest AS tdlib

RUN apk add --no-cache alpine-sdk linux-headers git zlib-dev openssl-dev gperf cmake

WORKDIR /usr/src/

RUN git clone --depth 1 https://github.com/tdlib/td.git

WORKDIR /usr/src/td/build/

RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local ..
RUN cmake --build . --target install


# Stage 2: Compile Go project

FROM golang:1.26.2-alpine AS go

RUN apk add --no-cache pkgconf gcc musl-dev

COPY --from=tdlib /usr/local/ /opt/tdlib/

WORKDIR /usr/src/

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY bridge ./bridge/
COPY tdlib ./tdlib/

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/opt/tdlib/include/"
ENV CGO_LDFLAGS="-L/opt/tdlib/lib/"

RUN go build -o ./release/app .


# Stage 3: Release stage

FROM alpine:latest AS release

RUN apk add --no-cache ca-certificates libstdc++

COPY --from=tdlib /usr/local/lib/libtdjson.so* /usr/local/lib/
COPY --from=go /usr/src/release/app /usr/local/bin/app

ENV LD_LIBRARY_PATH=/usr/local/lib/

RUN adduser -D -u 10001 app
USER app

ENTRYPOINT ["/usr/local/bin/app"]
