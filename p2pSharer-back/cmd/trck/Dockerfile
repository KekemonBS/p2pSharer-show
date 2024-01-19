FROM alpinelinux/golang:latest
ENV REDISURI="redis://:testPass@cache:6739/"
WORKDIR /opt
COPY . .

USER root
RUN go mod tidy
RUN go build -mod vendor ./cmd/trck/...
EXPOSE 8080

CMD ["./trck"]
