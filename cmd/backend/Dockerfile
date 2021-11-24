FROM golang:1.17 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/backend cmd/backend
COPY driver driver
COPY grpc/backend grpc/backend
COPY proto proto
COPY pkg pkg

WORKDIR /workspace/cmd/backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/cmd/backend/app .
USER 65532:65532

ENTRYPOINT ["/app"]
