FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o track .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /build/track .
COPY --from=builder /build/config.yaml .
COPY --from=builder /build/templates ./templates
COPY --from=builder /build/static ./static

RUN mkdir -p /app/data

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/track.db

CMD ["./track"]
