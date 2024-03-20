# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:latest AS build-stage

ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o llm_gateway

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app
COPY --from=build-stage /build/llm_gateway llm_gateway
EXPOSE 8055
USER nonroot:nonroot
ENTRYPOINT ["/app/llm_gateway", "-f", "/conf/config.yaml"]
