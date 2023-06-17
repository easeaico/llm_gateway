FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go build .

EXPOSE 5984
ENTRYPOINT ["./llm_mesh"]
