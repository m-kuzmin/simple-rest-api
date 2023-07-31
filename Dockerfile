ARG go_version=1.20
ARG base_image=alpine:latest

FROM golang:$go_version AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM $base_image

COPY --from=build /app/server /usr/bin/server
CMD ["server"]
