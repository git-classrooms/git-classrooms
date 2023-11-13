#############################################
# Builder stage image
#############################################
FROM golang:1.21-alpine as builder

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash  nodejs yarn

WORKDIR /app

COPY ./ ./
## Backend
RUN go mod tidy
RUN go build -o app ./main.go


## Frontend
RUN cd frontend && yarn
RUN cd frontend && yarn build


#############################################
# Runtime image
#############################################
FROM alpine:3.18 as release

RUN mkdir -p /app/public
WORKDIR /app

COPY --from=builder /app/app /app
COPY --from=builder /app/frontend/dist /app/public
EXPOSE 3000

ENTRYPOINT /app/app
