FROM golang:1.15-alpine as builder

LABEL maintainer="reoxey"

WORKDIR /micro/shipping

COPY . .

RUN CGO_ENABLED=0 go test . && go build -o app .

FROM alpine

COPY --from=builder /micro/shipping/app .

EXPOSE 8001

CMD ["/app"]
