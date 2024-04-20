FROM golang:1.22.1 as healthcheck-builder

WORKDIR /build

COPY healthcheck/go.mod healthcheck/go.sum ./

RUN go mod download

COPY healthcheck ./

RUN GCO_ENABLED=0 GOOS=linux go build -o /healthcheck

#### MAIN APP ####

FROM golang:1.22.1 as builder

WORKDIR /build

# Required dependency to build swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY go.sum go.mod ./

RUN go mod download

COPY main.go ./
COPY src src

# Create docs directory and documentation
RUN swag init

RUN GOOS=linux go build -o /turbex-api-server

#### STAGE 2 ####

FROM gcr.io/distroless/base-debian12 as production

WORKDIR /

COPY --from=healthcheck-builder /healthcheck /
COPY --from=builder /turbex-api-server /

USER nonroot:nonroot

# Security by default
ENV GIN_MODE=release

HEALTHCHECK --interval=5s --timeout=5s CMD ["/healthcheck"]
EXPOSE 8000

ENTRYPOINT ["/turbex-api-server"]

