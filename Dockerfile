# -- Stage 1 --
FROM golang:1.21.0-alpine3.18 AS builder

# Install upx
WORKDIR /source
RUN apk --no-cache add git upx

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY src/ ./src/
RUN go build -o dist/ ./...
RUN upx dist/*

# -- Stage 2 --
FROM alpine:3.18 AS runner

# Add non-root user
RUN adduser -D -h /opt/codelabs -s /sbin/nologin codelabs
WORKDIR /opt/codelabs
USER codelabs

# Copy files
COPY --from=builder /source/dist/src .
COPY sql/ ./sql/

# Run
EXPOSE 8080
ENV DB_CONNECTION_STRING "postgres://postgres:postgres@localhost:5432/codelabs?sslmode=disable"
ENV DB_MIGRATIONS_PATH "file:///opt/codelabs/sql/migrations"
ENV JWT_SECRET "default"
ENTRYPOINT ["/opt/codelabs/src"]