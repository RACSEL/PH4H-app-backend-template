# IPS Lacpass Backend

Backend for the IPS Lacpass App. A unified health app for Connectathon users to view, merge, and share IPS data
securely cross-border. This project provides a backend system composed of the IPS Lacpass API for handling business logic and
a Keycloak server for authentication and authorization. The entire stack is containerized using Docker and can be easily
managed with Docker Compose.

## Table of Contents

- [Project Overview](#project-overview)
- [Components](#components)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Environment Variables](/docs/environment.md)
  - [Running the Application](#running-the-application)
- [Accessing the Services](#accessing-the-services)
  - [Keycloak Admin Console](#keycloak-admin-console)
  - [Golang API](#golang-api)
- [Stopping the Application](#stopping-the-application)
- [Multiple Nodes Support](#multiple-nodes-support)
- [OpenAPI Documentation](#openapi-documentation)

## Project Overview

The architecture of this backend system is designed to separate concerns between the application's business logic and user authentication.

## Components

- **IPS Lacpass API**: A lightweight, high-performance API that implements the core features of your application.
  It is protected and requires a valid JWT from an Authorization server to be accessed.
- **Authorization**: Identity and Access Management solution. It handles user registration, login, and token issuance.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Go 1.24](https://go.dev/dl/) `Only if you plan to ran the API without docker`
- [JQ](https://jqlang.org/download/) `If you plan on running shell scripts`

## Getting Started

### Configuration

- [Authentication](/docs/authentication.md)
- IPS Lacpass API (WIP)

### Running the Application

First make a local copy for the sample environment variables file:

```bash
~ cp .env.sample .env
```

Edit the `.env` file to match you server configuration (More information about [Enviroment Variables](/docs/environment.md)).
Then, open a terminal in the root directory of the project and run the following command:

```bash
~ docker compose up
[+] up 26/26
 ✔ Image haravich/fake-smtp-server:20250615 Pulled
 ✔ Image keycloak/keycloak:26.2.5           Pulled
 ✔ Image postgres:17.5-alpine               Pulled
[+] Building 1.2s (24/24) FINISHED
...
 ✔ Image haravich/fake-smtp-server:20250615 Pulled
 ✔ Image keycloak/keycloak:26.2.5           Pulled
 ✔ Image postgres:17.5-alpine               Pulled
 ✔ Image docker-lacpass-backend             Built
 ✔ Network docker_auth                      Created
 ✔ Network docker_backend                   Created
 ✔ Container mailcatcher                    Created
 ✔ Container auth-db                        Created
 ✔ Container auth                           Created
 ✔ Container lacpass-backend                Created
Attaching to auth, auth-db, lacpass-backend, mailcatcher
Container auth-db Waiting
auth-db  |
auth-db  | PostgreSQL Database directory appears to contain a database; Skipping initialization
auth-db  |
auth-db  | 2026-05-03 06:33:54.841 UTC [1] LOG:  starting PostgreSQL 17.5 on x86_64-pc-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
auth-db  | 2026-05-03 06:33:54.841 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
auth-db  | 2026-05-03 06:33:54.841 UTC [1] LOG:  listening on IPv6 address "::", port 5432
auth-db  | 2026-05-03 06:33:54.853 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
auth-db  | 2026-05-03 06:33:54.864 UTC [29] LOG:  database system was shut down at 2026-05-03 06:33:29 UTC
auth-db  | 2026-05-03 06:33:54.871 UTC [1] LOG:  database system is ready to accept connections
mailcatcher  | Starting MailCatcher v0.10.0
mailcatcher  | ==> smtp://0.0.0.0:1025
mailcatcher  | ==> http://0.0.0.0:1080
```

This command will:

- Build the Docker image for the IPS Lacpass API.
- Pull the official Docker images for Keycloak and Postgres.
- Create and start the containers for all three services.
- Attach your terminal to the logs of all running containers.

If you dont want to attach to the running containers,
you can add `--detach` at the end of the command. To stop the container just do `CTRL+C` and then make sure to take
down the containers with:

```bash
~ docker compose down
  [+] down 6/6
   ✔ Container lacpass-backend Removed
   ✔ Container mailcatcher     Removed
   ✔ Container auth            Removed
   ✔ Container auth-db         Removed
   ✔ Network docker_auth       Removed
   ✔ Network docker_backend    Removed
```

### Setup Keycloak

After all services are running, you need to setup keycloak to have the correct configurations for the backend to authenticate and create users. Before continuing, please follow the instructions [here](/docs/keycloak-setup.md).

## Accessing the Services

### Keycloak Admin Console

Once the services are running, you can access the Keycloak Admin Console to configure realms, clients, and users.

1.  Open your web browser and navigate to `http://localhost:9083`.
2.  You will be redirected to the Keycloak landing page. Click on the **Administration Console** link.
3.  Log in with the admin credentials provided in your `.env` file (`KC_BOOTSTRAP_ADMIN_USERNAME` and `KC_BOOTSTRAP_ADMIN_PASSWORD`).

### IPS Lacpass API

IPS Lacpass API will be accessible at `http://localhost:9081`. You can use a tool like `curl` or Postman to interact
with your API endpoints. Remember that your API endpoints will be protected by Keycloak, so you will need to obtain a
valid JWT from Keycloak to make successful requests. There is a [helper script](./scripts/auth.sh), where you can request
a token using. For it to work, you need to set up a user and add that credentials in your `.env` file.
The steps are details in our [authentication guide](/docs/authentication.md). After you will be able to run:

```bash
~ sh scripts/auth.sh access-token
Successfully logged in!
Access Token: XXXXX.VVVV.BBBB
```

If the token expires you can refresh it with:

```bash
~ sh scripts/auth.sh refresh-token
Successfully refreshed token!
Access Token: XXXXX.VVVV.BBBB
```

And to logout you can do:

```bash
~ sh scripts/auth.sh logout
Success: Logout successful. The refresh token has been invalidated.
```

## Stopping the Application

To stop and remove the containers, network, and volumes, press `Ctrl+C` in the terminal where `docker-compose` is running, and then run the following command:

```bash
~ docker compose down
[+] down 6/6
 ✔ Container mailcatcher     Removed
 ✔ Container lacpass-backend Removed
 ✔ Container auth            Removed
 ✔ Container auth-db         Removed
 ✔ Network docker_backend    Removed
 ✔ Network docker_auth       Removed
```

## Multiple Nodes Support

The backend supports connecting to multiple national nodes within a single instance. This feature is optional and disabled by default.

### Enabling Multiple Nodes

Set the following environment variable in your `.env` file:

```env
USE_MULTIPLE_NODES=1
```

### Configuration

Create a file named `node-services.json` in the root directory of the project. This file defines the services for each node:

```json
[
    {
        "id": "lacpass",
        "name": "lacpass",
        "FHIR_BASE_URL": "http://lacpass.create.cl:8080",
        "FHIR_MEDIATOR_BASE_URL": "http://lacpass.create.cl:3000",
        "VHL_BASE_URL": "http://lacpass.create.cl:8182",
        "ICVP_VALIDATOR_URL": "http://lacpass.create.cl:7100"
    },
    {
        "id": "itb",
        "name": "ITB",
        "FHIR_BASE_URL": "http://itb.racsel.cl/fhir",
        "FHIR_MEDIATOR_BASE_URL": "http://itb.racsel.cl/mediator",
        "VHL_BASE_URL": "http://itb.racsel.cl/vhl",
        "ICVP_VALIDATOR_URL": "http://itb.racsel.cl/validator"
    }
]
```

### Usage

When this feature is enabled, all API endpoints accept an optional `Node-Name` request header. This header receives the node `id` and determines which node configuration to use for that specific request. The `name` is only for display.

Example:
`GET /ips` with header `Node-Name: itb`

If `Node-Name` is omitted or the id is not found in the configuration, the backend will fall back to the default services defined by the standard environment variables (`FHIR_BASE_URL`, etc.).

# OpenAPI Documentation

The OpenAPI specification is defined in the project [api.yaml](/api/openapi/api.yaml). Our docker compose creates a
container that serve it. You can run only swagger doing:

```bash
~ docker compose up swagger-ui
```

Then, the API docs can be seen in `http://localhost:9999`. Or whatever port you defined in the docker compose.
