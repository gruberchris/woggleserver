FROM golang:1.20-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o woggleserver .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app/ /app/
EXPOSE 3000
ENTRYPOINT ["./woggleserver"]
