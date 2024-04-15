FROM golang:1.20 as builder

RUN mkdir /build
WORKDIR /build

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/scanner

# Stage 2
FROM alpine:3.16

ENV TZ=Etc/UTC
RUN apk add --no-cache tzdata
RUN mkdir -p /app

COPY --from=builder /build/main /app/

EXPOSE 8080

WORKDIR /app

CMD ["./main"]