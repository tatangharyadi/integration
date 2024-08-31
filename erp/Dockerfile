
## Build
FROM golang:1.23-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o /integration-erp

## Deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /integration-erp /integration-erp

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/integration-erp"]
