#############################################
# Preparer go
#############################################
FROM golang:1.22-alpine as preparer-go

WORKDIR /app/build

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN go generate

#############################################
# Swagger client
#############################################
FROM swaggerapi/swagger-codegen-cli-v3 as swagger-client-builder

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
# Builder web
#############################################
FROM node:20-alpine as builder-web

WORKDIR /app/build
COPY ./frontend/package.json ./frontend/yarn.lock ./
RUN yarn --frozen-lockfile

COPY ./frontend ./
COPY --from=swagger-client-builder /app/build/swagger-client ./src/swagger-client
RUN yarn build

#############################################
# Builder go
#############################################
ARG APP_VERSION="v0.0.0"
FROM preparer-go as builder-go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$APP_VERSION -s -w" -o /app/build/app

#############################################
# Runtime image
#############################################
FROM scratch

ENV FRONTEND_PATH=/public
ENV PORT=3000
EXPOSE 3000

COPY --from=builder-go /app/build/app /app
COPY --from=builder-web /app/build/dist /public

ENTRYPOINT ["/app"]
