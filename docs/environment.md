# Environment

This document describes the required environment variables to run the micro-service.

| Name                          | Description                                                                                                    | Example                                                         | Mandatory |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- | --------- |
| `DB_CONNECTION_STRING`        | Connection string to the database.                                                                             | `postgres://user:password@domain:port/database?sslmode=disable` | Yes       |
| `RABBIT_MQ_CONNECTION_STRING` | Connection string to RabbitMQ.                                                                                 | `amqp://username:password@address:port/`                        | Yes       |
| `WEB_CLIENT_URL`              | URL of the web client, used to configure CORS.                                                                 | `http://domain:5173`                                            | Yes       |
| `JWT_SECRET`                  | Secret used to sign JWT tokens.                                                                                | `secret`                                                        | Yes       |
| `JWT_EXPIRATION_HOURS`        | Expiration time of JWT tokens.                                                                                 | `6`                                                             | No        |
| `DB_MIGRATIONS_PATH`          | Path to the database migrations.                                                                               | `file://path`                                                   | No        |
| `ARCHIVES_MAX_SIZE_KB`        | Maximum size of the archives in KB.                                                                            | `1024`                                                          | No        |
| `EXEC_ENVIRONMENT`            | Environment where the micro-service is running. It is just being used to close `SSE` connections during tests. | `testing`                                                       | No        |
