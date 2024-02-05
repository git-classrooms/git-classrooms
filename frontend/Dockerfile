#############################################
# Development Image with nodejs
#############################################
FROM node:lts-alpine as builder

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash yarn

WORKDIR /app

USER node
CMD yarn && yarn dev
