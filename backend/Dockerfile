FROM golang:1.22.1 as builder

WORKDIR /build

# Required dependency to build swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY go.sum go.mod .

RUN go mod download

COPY main.go .
COPY src src

# Create docs directory and documentation
RUN swag init

RUN GCO_ENABLED=0 GOOS=linux go build -o /turbex-api-server

#### STAGE 2 ####

# TODO: Investigate why we need busybox instead of scratch
FROM busybox as production

WORKDIR /
USER 1728:1728

COPY --from=builder /turbex-api-server /

# Security by default
ENV GIN_MODE=release

EXPOSE 8000

ENTRYPOINT ["/turbex-api-server"]

