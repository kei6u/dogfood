FROM golang:1.17 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/gateway cmd/gateway
COPY driver driver
COPY proto proto
COPY pkg pkg

WORKDIR /workspace/cmd/gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/cmd/gateway/app .
USER 65532:65532

ENTRYPOINT ["/app"]
