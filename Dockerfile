# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:latest AS build-stage

ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /llm_mesh

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /llm_mesh /llm_mesh

EXPOSE 5984

USER nonroot:nonroot

ENTRYPOINT ["/llm_mesh", "-f", "/conf/config.yaml"]
