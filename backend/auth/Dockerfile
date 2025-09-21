FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy go.work และ go.sum
COPY go.work go.work.sum ./

# copy go.mod/go.sum ของทุก module
COPY auth/go.mod auth/go.sum ./auth/
COPY crud/go.mod crud/go.sum ./crud/
COPY database/go.mod database/go.sum ./database/

# copy source code ทั้งหมด (auth, crud, database)
COPY auth ./auth
COPY crud ./crud
COPY database ./database

# build service
WORKDIR /app/auth/cmd/server
RUN go build -o /auth_service .

# stage 2: runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /auth_service .

EXPOSE 8081
CMD ["./auth_service"]
