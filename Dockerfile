ARG go_version=1.20
ARG base_image=alpine:latest

FROM golang:$go_version AS build
WORKDIR /app

# Prepare dependecies
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download

# Build the app
COPY . .
RUN swag init -o api/swaggerui --ot go
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Copy the build application to a clean image
FROM $base_image
EXPOSE 8000

COPY --from=build /app/server /usr/bin/server
COPY db/migrations/ /migrations
CMD ["server"]
