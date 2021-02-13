FROM golang:1.15-alpine as builder

LABEL maintainer="reoxey"

WORKDIR /micro/product

COPY . .

RUN CGO_ENABLED=0 go test . && go build -o app .

FROM alpine

COPY --from=builder /micro/product/app .

EXPOSE 8002

CMD ["/app"]