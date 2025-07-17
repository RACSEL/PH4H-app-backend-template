#!/bin/bash

source ./.env

if ! command -v jq >/dev/null 2>&1
then
    echo "JQ could not be found"
    exit 1
fi

TOKEN_ENDPOINT="$KEYCLOAK_URL/realms/$KEYCLOAK_REALM/protocol/openid-connect/token"

get_access_token() {
  RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "client_id=$KEYCLOAK_CLIENT_ID" \
    -d "username=$KEYCLOAK_DEFAULT_USER" \
    -d "password=$KEYCLOAK_DEFAULT_USER_PASSWORD" \
    -d "scope=openid" \
    -d "grant_type=password")

  if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to Keycloak."
    exit 1
  fi

  # Extract the access token from the JSON response using jq
  ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r .access_token)
  REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r .refresh_token)
  mkdir -p ./tmp
  echo -n "$REFRESH_TOKEN" > ./tmp/refresh_token

  # Check if an access token was returned
  if [[ -z "${ACCESS_TOKEN}" || "$ACCESS_TOKEN" == null ]]; then
    echo "Error: Failed to obtain access token. Check your credentials and client configuration."
    echo "Response from Keycloak: $RESPONSE"
    exit 1
  fi

  echo "Successfully logged in!"
  echo "Access Token: $ACCESS_TOKEN"
  exit 0
}

get_refresh_token() {
  REFRESH_TOKEN=$(cat ./tmp/refresh_token)
  if [[ -z "${REFRESH_TOKEN}" || "$REFRESH_TOKEN" == null ]]; then
    echo "Could not find refresh token, getting a new access token"
    get_access_token
  fi
  RESPONSE=$(curl -s -X POST "$TOKEN_ENDPOINT" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      -d "client_id=$KEYCLOAK_CLIENT_ID" \
      -d "refresh_token=$REFRESH_TOKEN" \
      -d "grant_type=refresh_token")

  if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to Keycloak."
    exit 1
  fi

  # Extract the access token from the JSON response using jq
  ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r .access_token)

  # Check if an access token was returned
  if [[ -z "${ACCESS_TOKEN}" || "$ACCESS_TOKEN" == null ]]; then
    echo "Error: Failed to refresh token. Check your credentials and client configuration."
    echo "Response from Keycloak: $RESPONSE"
    exit 1
  fi

  echo "Successfully refreshed token!"
  echo "Access Token: $ACCESS_TOKEN"
  exit 0
}

logout() {
  LOGOUT_ENDPOINT="$KEYCLOAK_URL/realms/$KEYCLOAK_REALM/protocol/openid-connect/logout"
  REFRESH_TOKEN=$(cat ./tmp/refresh_token)
  if [[ -z "${REFRESH_TOKEN}" || "$REFRESH_TOKEN" == null ]]; then
    echo "Could not find refresh token, cannot logout"
    exit 1
  fi
  RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "client_id=$KEYCLOAK_CLIENT" \
        -d "refresh_token=$REFRESH_TOKEN" \
        "$LOGOUT_ENDPOINT")
  if [ "$RESPONSE" -eq 204 ]; then
      echo "Success: Logout successful. The refresh token has been invalidated."
      exit 0
  elif [ "$RESPONSE" -eq 400 ]; then
      echo "Error: Bad Request (400). The refresh token was likely invalid or already revoked."
  elif [ "$RESPONSE" -eq 401 ]; then
      echo "Error: Unauthorized (401). Check if your KEYCLOAK_CLIENT_SECRET is correct."
  else
      echo "Error: An unexpected error occurred. Keycloak responded with HTTP status $RESPONSE."
  fi
  exit 1
}

case $1 in
  access-token)
    get_access_token
  ;;
  refresh-token)
    get_refresh_token
  ;;
  logout)
    logout
  ;;
esac

