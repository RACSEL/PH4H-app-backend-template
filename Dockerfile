# Swag stage to generate swagger docs
FROM golang:1.24-alpine AS docs

RUN apk add --no-cache git

WORKDIR /build

RUN go install github.com/swaggo/swag/cmd/swag@latest

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . /context
RUN cp -R /context/cmd/api/. ./ \
    && cp -R /context/internal ./internal \
    && cp -R /context/pkg ./pkg \
    && if [ -f /context/node-services.json ]; then \
        cp /context/node-services.json ./node-services.json; \
    else \
        printf '[]\n' > ./node-services.json; \
    fi

RUN swag fmt
RUN swag init -o pkg/docs

# Go stage to build the API
FROM golang:1.24 AS builder

WORKDIR /build
COPY --from=docs /build .

RUN CGO_ENABLED=0 GOOS=linux go build -o api .
EXPOSE 8080
EXPOSE 3000

FROM scratch
# Include trusted CA roots so HTTPS requests from the scratch image can verify TLS certificates.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/api .
COPY --from=builder /build/node-services.json .
ENTRYPOINT ["./api"]
