# Environment Variables
The project multiple containers rely on environment variables stored in `.env` file located in the root folder.
Copy the given template file with the command `cp .env.sample .env` in the root folder, and define the values of each variable according to your setup.

For security reasons, please set obfuscate values for usernames, passwords and secret key variables.

This app uses [Keycloak](https://www.keycloak.org/) to handle OAuth2 user authentication. (more information in [Authentication](docs/authentication.md)), 
an so many enviromental variables of this project are set for a Keycloak setup.
Please define the variables of your need if you a different authentication service.

## Definitions

`API_PORT`
Port number where the app runs. Default: `3000`

`KEYCLOAK_URL` 
Keycloak service endpoint. Default: `http://localhost:9083`

`KEYCLOAK_REALM`
App users keycloak realm. Used for initial keycloak configuration, if changed, keycloak container and volume must be rebuild. Default: `lacpass`

`KEYCLOAK_CLIENT_ID`
App users keycloak client id. User for initial keycloak configuration, if changed, keycloak container and volume must be rebuild. Default: `app`

`KEYCLOAK_ADMIN_CLIENT_SECRET`
Secret key of the `admin-cli` keycloak client of your instance. No default value since it must be set to the one generated in your Keycloak instance during setup.

`KEYCLOAK_HOSTNAME`
Keycloak service hostname used for flows like password recovery. Default: `http://keycloak.lacpass.create.cl`

`KC_BOOSTRAP_ADMIN_USERNAME`
Keycloak admin console username. Default: `admin`

`KC_BOOTSTRAP_ADMIN_PASSWORD`
Keyclaok admin console password. Default: `admin`

`KEYCLOAK_DEFAULT_USER`
Testing user username. For testing porposes, the app creates a first mock user instance in the keycloak service. Default: `test`

`KEYCLOAK_DEFAULT_USER_PASSWORD`
Testing user password. Default: `test`

`API_SWAGGER`
Enable `/swagger/index.html` endpoint. Boolean value, can be `true` or `false`. Default: `false`

`POSTGRES_USER`
Postgres database default role name. Default: `postgres`

`POSTGRES_PASSWORD`
Postgres database default role password. Default: `postgres`

`FHIR_BASE_URL`
Fhir server endpoint for FHIR IPS managment. Default: `http://lacpass.create.cl:8080`

`VHL_BASE_URL`
VHL server endpoint for VHL QR generation, validation and retrieve. Default: `http://lacpass.create.cl:8182`

`ICVP_VALIDATOR_URL`
Endpoint to validate the QR content of ICVPs not linked to an IPS. Default: `http://lacpass.create.cl:7089`
