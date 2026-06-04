FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/memere-api ./cmd/api

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /app/memere-api /app/memere-api

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/memere-api"]
