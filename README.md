# PH4H App Backend Template

Template for backend for the IPS PH4H App. A unified health app for Connectathon users to view, merge, and share IPS data
securely cross-border. This project provides a backend system composed of the IPS PH4H API for handling business logic and
a Keycloak server for authentication and authorization. The entire stack is containerized using Docker and can be easily
managed with Docker Compose.

## Table of Contents

- [Project Overview](#project-overview)
- [Components](#components)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Accessing the Services](#accessing-the-services)
  - [Keycloak Admin Console](#keycloak-admin-console)
  - [Golang API](#golang-api)
- [Stopping the Application](#stopping-the-application)
- [Swagger Documentation](#swagger-documentation)

## Project Overview

The architecture of this backend system is designed to separate concerns between the application's business logic and user authentication.

## Components

- **IPS PH4H API**: A lightweight, high-performance API that implements the core features of your application.
  It is protected and requires a valid JWT from an Authorization server to be accessed.
- **Authorization**: Identity and Access Management solution. It handles user registration, login, and token issuance.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Go 1.24](https://go.dev/dl/) `Only if you plan to ran the API without docker`

## Getting Started

### Configuration

- [Authentication](/docs/authentication.md)
- IPS PH4H API (WIP)

### ⚠️ Complete the not implemented calls

This template backend application includes some functionalities that are intentionally left unimplemented. These are meant to be completed by the participant to assess their knowledge of the FHIR/VHL standards.

To complete these functionalities, search the code for the `TODO: To be implemented by the participant` message. You will find 3 missing functionaties for:

- Implementing an ITI-67 call
- Implementing an ITI-68 call
- Implementing a QR code generation using VHL

After completing these steps, you can proceed to the next section on running the application.


### Running the Application

Open a terminal in the root directory of the project and run the following command:

```bash
 docker compose --file=./docker/compose.yaml up
```

This command will:

- Build the Docker image for the IPS PH4H API.
- Pull the official Docker images for Keycloak and Postgres.
- Create and start the containers for all three services.
- Attach your terminal to the logs of all running containers.

### Setup Keycloak 

After all services are running, you need to setup keycloak to have the correct configurations for the backend to authenticate and create users. Before contuining please follow the instructions [here](/docs/keycloak-setup.md).

## Accessing the Services

### Keycloak Admin Console

Once the services are running, you can access the Keycloak Admin Console to configure realms, clients, and users.

1.  Open your web browser and navigate to `http://localhost:9083`.
2.  You will be redirected to the Keycloak landing page. Click on the **Administration Console** link.
3.  Log in with the admin credentials provided in your [configuration](/docs/authentication.md)

### IPS PH4H API

IPS PH4H API will be accessible at `http://localhost:9081`. You can use a tool like `curl` or Postman to interact 
with your API endpoints. Remember that your API endpoints will be protected by Keycloak, so you will need to obtain a 
valid JWT from Keycloak to make successful requests. There is a [helper script](./scripts/auth.sh), where you can request 
a token using:

```bash
sh scripts/auth.sh access-token
```

If the token expires you can refresh it with:

```bash
sh scripts/auth.sh refresh-token
```

And to logout you can do:

```bash
sh scripts/auth.sh logout
```

## Stopping the Application

To stop and remove the containers, network, and volumes, press `Ctrl+C` in the terminal where `docker-compose` is running, and then run the following command:

```bash
docker-compose down
```

# Swagger Documentation

To activate Swagger make sure to set `API_SWAGGER=true` in your `.env` file and the `docker/compose.yaml`.

Then, the API docs can be seen in `http://localhost:9081/swagger/index.html`.
