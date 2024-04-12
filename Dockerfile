#############################################
# Builder web
#############################################
FROM node:20-alpine as builder-web

WORKDIR /app/build
COPY ./frontend/package.json ./frontend/yarn.lock ./
RUN yarn --pure-lockfile

COPY ./frontend ./
RUN yarn build

#############################################
# Builder go
#############################################
ARG APP_VERSION="v0.0.0"
FROM golang:1.22-alpine as builder-go

# install mockery
RUN go install github.com/vektra/mockery/v2@v2.42.2

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app/build

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$APP_VERSION" -o /app/build/app

#############################################
# Runtime image
#############################################
FROM alpine:3.18 as release

ENV FRONTEND_PATH=/public
ENV PORT=3000
EXPOSE 3000

RUN adduser -D gorunner

USER gorunner

WORKDIR /

COPY --chown=gorunner:gorunner --from=builder-go /app/build/app /app
COPY --chown=gorunner:gorunner --from=builder-web /app/build/dist /public

ENTRYPOINT /app
