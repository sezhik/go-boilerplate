FROM golang:1.25-alpine as base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@v1.63.4
WORKDIR /app
COPY . .
CMD ["air"]

FROM base as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o service ./cmd/service

FROM alpine:latest as prod
WORKDIR app
COPY --from=builder /app/service .
ENTRYPOINT ["/app/service"]
