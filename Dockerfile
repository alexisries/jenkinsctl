FROM golang:1.16-alpine as build-env

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY pkg/ ./pkg/

RUN apk add upx

# Build without CGO and DWARF debugging informations (reduce the binary size & execution time)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o jenkinsctl .

# Binary compression (reduce size)
RUN upx --best --lzma jenkinsctl

FROM scratch

COPY --from=build-env /app/jenkinsctl /usr/bin/

WORKDIR /app

ENTRYPOINT [ "jenkinsctl" ]
