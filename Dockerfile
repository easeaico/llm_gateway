FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go build .

EXPOSE 4000
ENTRYPOINT ["./llmmesh"]
