# syntax=docker/dockerfile:1.7-labs
#############################################
#                Preparer go                #
#############################################
FROM golang:1.25-alpine3.23 AS preparer-go

RUN apk add --no-cache make git

WORKDIR /app/build
RUN mkdir -p ./frontend/dist && touch ./frontend/dist/robots.txt

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY --exclude=frontend ./ ./
RUN go generate

#############################################
#              Swagger client               #
#############################################
FROM swaggerapi/swagger-codegen-cli-v3 AS swagger-client-builder

WORKDIR /app/build

COPY --from=preparer-go /app/build/docs ./docs

RUN java -jar /opt/swagger-codegen-cli/swagger-codegen-cli.jar generate -i ./docs/swagger.yaml -l typescript-axios -o ./swagger-client
RUN cd swagger-client && \
    rm -rf .swagger-codegen && \
    rm .gitignore .swagger-codegen-ignore .npmignore git_push.sh package.json README.md tsconfig.json && \
    for file in apis/*.ts; \
    do sed -i '1s;^;\/\/ @ts-nocheck\n;' $file; \
    done

#############################################
#                Builder web                #
#############################################
FROM node:26-alpine AS builder-web

WORKDIR /app/build
COPY ./frontend/package.json ./frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY ./frontend ./
COPY --from=swagger-client-builder /app/build/swagger-client ./src/swagger-client
RUN pnpm build

#############################################
#                Builder go                 #
#############################################
FROM preparer-go AS builder-go

ARG APP_VERSION="v0.0.0-dev"
ARG APP_GIT_COMMIT="unknown"
ARG APP_GIT_BRANCH="develop"
ARG APP_GIT_REPOSITORY="https://github.com/git-classrooms/git-classrooms"
ARG APP_BUILD_TIME="unknown"

COPY --from=builder-web /app/build/dist ./frontend/dist
RUN make CI=true build \
    APP_VERSION=${APP_VERSION} \
    APP_GIT_COMMIT=${APP_GIT_COMMIT} \
    APP_GIT_BRANCH=${APP_GIT_BRANCH} \
    APP_GIT_REPOSITORY=${APP_GIT_REPOSITORY} \
    APP_BUILD_TIME=${APP_BUILD_TIME}


#############################################
#               Runtime image               #
#############################################
FROM alpine:3.23 AS release

RUN apk add --no-cache tzdata

ENV PORT=3000
EXPOSE 3000

RUN adduser -D gorunner

USER gorunner

WORKDIR /app

COPY --chown=gorunner:gorunner --from=builder-go /app/build/bin/git-classrooms /app/git-classrooms

ENTRYPOINT ["/app/git-classrooms"]
