FROM golang:1.15-alpine as builder

LABEL maintainer="reoxey"

WORKDIR /micro/order

COPY . .

RUN CGO_ENABLED=0 go test . && go build -o app .

FROM scratch

WORKDIR /micro

COPY --from=builder /micro/order .

CMD ["./app"]
