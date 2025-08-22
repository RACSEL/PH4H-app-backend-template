# Swag stage to generate swagger docs
FROM golang:1.24-alpine AS docs

RUN apk add --no-cache git

WORKDIR /build

RUN go install github.com/swaggo/swag/cmd/swag@latest

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/api/ ./
COPY ./internal ./internal
COPY ./pkg      ./pkg

RUN swag fmt
RUN swag init -o pkg/docs

# Go stage to build the API
FROM golang:1.24 as builder

WORKDIR /build
COPY --from=docs /build .

RUN CGO_ENABLED=0 GOOS=linux go build -o api .
EXPOSE 8080
EXPOSE 3000

FROM scratch
COPY --from=builder /build/api .
ENTRYPOINT ["./api"]
