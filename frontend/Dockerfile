#############################################
# Development Image with nodejs
#############################################
FROM node:lts-alpine as builder

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash yarn

USER node
WORKDIR /app
COPY ./package.json ./yarn.lock ./
RUN yarn --frozen-lockfile

COPY ./ ./
CMD yarn dev
