#############################################
# Development Image with nodejs
#############################################
FROM node:lts-alpine as builder

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash yarn

USER 1000
WORKDIR /app

CMD yarn && yarn dev
