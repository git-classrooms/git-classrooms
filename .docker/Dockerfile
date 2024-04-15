FROM cosmtrek/air:latest

RUN adduser air

USER air

WORKDIR /app

RUN go install github.com/vektra/mockery/v2@v2.42.2
RUN go install github.com/swaggo/swag/cmd/swag@latest
