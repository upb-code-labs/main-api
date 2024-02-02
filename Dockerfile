# -- Stage 1 --
FROM docker.io/golang:1.21.6-alpine3.19 AS builder

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
FROM docker.io/alpine:3.19.1 AS runner

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