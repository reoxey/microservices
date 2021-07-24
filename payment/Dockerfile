FROM golang:1.16-alpine as builder

LABEL maintainer="reoxey"

WORKDIR /micro/payment

COPY . .

RUN CGO_ENABLED=0 go test . && go build -o app .

FROM alpine

COPY --from=builder /micro/payment/app .

EXPOSE 8009

CMD ["/app"]