# General organization contribution guidelines

Please, read the [General Contribution Guidelines](https://github.com/upb-code-labs/docs/blob/main/CONTRIBUTING.md) before contributing.

# Gateway repository contribution guidelines

Please, read the following guidelines before contributing to this repository and remember to update the [Bruno collection](./bruno/), [Insomnia collection](./insomnia/),[OpenAPI Specification](./openapi/spec.openapi.yaml) and [Tests](../__tests__/) when contributing to this repository.

## Project structure / architecture

The gateway service's architecture is based on the Hexagonal Architecture and the Vertical Slice Architecture in order to make it easier to maintain, change and scale.

Note that, for each entity in the REST API, there is a folder with the following structure:

- `domain`:

  - `definitions`: Contains interfaces defining contracts for the repositories and services needed by the entity.
  - `dtos`: Contains the data transfer objects used to transfer data between the different layers.
  - `entities`: Contains the entity's definition.
  - `errors`: Contains the custom domain errors used by the entity.

- `application`: Contains the application logic (use cases) of the entity.

- `infrastructure`:

  - `http`: Contains the implementation of the REST API endpoints (routes and controllers).
  - `implementations`: Contains the implementation of the repositories and services needed by the entity.
  - `requests`: Contains the request objects used to parse the requests received form the clients.
  - `responses`: Contains the response objects used to send responses to the clients.

The above structure allows us to easily change the implementation of the repositories and services (Lets say, from a Postgres database to a MongoDB database) without having to change the application logic.

Furthermore, thanks to the Vertical Slice Architecture, we can extract the entity's folder (And its dependencies) and move it to a different repository, making it easier to scale the system.

## Local development

The following dependencies are required to run the gateway service locally:

- [Go 1.21.5](https://golang.org/doc/install)
- [Podman](https://podman.io/getting-started/installation) (To build and test the container image)
- [Podman Compose](https://github.com/containers/podman-compose)
- [Bruno](https://www.usebruno.com/) (To test the REST API)
- [Insomnia](https://insomnia.rest/) (To test the endpoints that require sending files through multipart/form-data)
- [GNU Make](https://www.gnu.org/software/make/)

Please, note that `Podman` and `Podman Compose` are a drop-in replacement for `Docker` and `Docker Compose` respectively, so you can use the latter if you prefer.

Additionally, you may want to install the following dependencies to make your life easier:

- [Air](https://github.com/cosmtrek/air) (for live reloading)

## Running the gateway service locally

First, you need to run the `docker-compose.yaml` file located in the root of the repository in order to start the dependencies of the gateway service (Postgres, RabbitMQ and the microservices):

```bash
podman-compose up -d
```

Then, you can run the gateway service with the following command:

```bash
air
```

This will start the gateway service and will watch for changes in the source code and restart the service automatically.
