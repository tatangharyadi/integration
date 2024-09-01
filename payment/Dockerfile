
## Build
FROM golang:1.23-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o /integration-payment

## Deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /integration-payment /integration-payment

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/integration-payment"]
