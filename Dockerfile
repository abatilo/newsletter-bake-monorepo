# syntax=docker/dockerfile:1.4
FROM node:18.10.0-alpine as node-modules

WORKDIR /app/web
COPY ./web/package.json ./web/package-lock.json ./
COPY ./web/app1/package.json ./app1/
COPY ./web/app2/package.json ./app2/
RUN npm ci --no-audit

FROM node-modules as node-builder
COPY ./web ./

ARG app
RUN npm run build --workspace ${app}

FROM golang:1.19.0 as go-modules

# Install dependencies
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

FROM go-modules as go-builder

COPY ./internal ./internal
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/...

FROM gcr.io/distroless/static-debian11:nonroot
ENV STATIC_ASSETS_PATH /static

ARG app
COPY --from=node-builder /app/web/${app}/build /static
COPY --from=go-builder /go/bin/${app} /usr/local/bin/entrypoint
ENTRYPOINT ["entrypoint"]
