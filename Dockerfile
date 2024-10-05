# syntax=docker/dockerfile:1.7-labs
#############################################
#                Preparer go                #
#############################################
FROM golang:1.22-alpine AS preparer-go

# install mockery
RUN go install github.com/vektra/mockery/v2@v2.42.2

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

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
FROM node:20-alpine AS builder-web

WORKDIR /app/build
COPY ./frontend/package.json ./frontend/yarn.lock ./
RUN yarn --frozen-lockfile

COPY ./frontend ./
COPY --from=swagger-client-builder /app/build/swagger-client ./src/swagger-client
RUN yarn build

#############################################
#                Builder go                 #
#############################################
FROM preparer-go AS builder-go

COPY --from=builder-web /app/build/dist ./frontend/dist

ARG APP_VERSION="v0.0.0"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X main.version=$APP_VERSION" -o /app/build/app

#############################################
#               Runtime image               #
#############################################
FROM scratch AS release

ENV PORT=3000
EXPOSE 3000

COPY --from=builder-go /app/build/app /app

ENTRYPOINT ["/app"]
