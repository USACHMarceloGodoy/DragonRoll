# syntax=docker/dockerfile:1.7
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Cache de dependencias.
COPY go.mod go.sum* ./
RUN go mod download

# Código.
COPY . .

# Binario estático (sin CGO) → corre en distroless/scratch.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/dragonroll ./

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/dragonroll /app/dragonroll
# Seeds (skills universales + ejemplos homebrew).
COPY --from=builder /src/data /app/data

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/dragonroll"]
