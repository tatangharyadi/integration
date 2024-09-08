
## Build
FROM golang:1.23-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o /integration-database

## Deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /integration-database /integration-database

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/integration-database"]
