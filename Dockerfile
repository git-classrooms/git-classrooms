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
FROM golang:1.21-alpine as builder-go

# RUN apk --no-cache --update-cache --available upgrade \
#    && apk add git bash  nodejs yarn

WORKDIR /app/build

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

#############################################
# Runtime image
#############################################
FROM alpine:3.18 as release

ENV FRONTEND_PATH=/public
ENV TemplateFilePath=/templates/template.html
ENV PORT=3000
EXPOSE 3000

RUN adduser -D gorunner

USER gorunner

WORKDIR /

COPY --chown=gorunner:gorunner --from=builder-go /app /app
COPY --chown=gorunner:gorunner --from=builder-go /app/build/repository/mail/template.html /templates/template.html
COPY --chown=gorunner:gorunner --from=builder-web /app/build/dist /public

ENTRYPOINT /app
