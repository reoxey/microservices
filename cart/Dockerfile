FROM golang:1.15-alpine as builder

LABEL maintainer="reoxey"

WORKDIR /micro/cart

COPY . .

RUN CGO_ENABLED=0 go test . && go build -o app .

FROM alpine

COPY --from=builder /micro/cart/app .

EXPOSE 8003

CMD ["/app"]