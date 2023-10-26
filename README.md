# Gitlab-Classrooms

## Project Structure 

The project has two parts:

- Golang with fiber as backend=> `/`
- React with Typescript and vite as Frontend => `/frontend/`

The frontend proxies the requests for the path `/api/*` to the backend server.

## Development

For development we use the git flow branching model for simplicity. 

### Setup without docker

Install air [cosmtrek/air](https://github.com/cosmtrek/air) and run the following:

```
./script/start
```

- Frontend Dev server is listening at (localhost:5173)[http://localhost:5173]
- Backend server is listening at (localhost:3000)[http://localhost:3000]

**NOTE:** You need to setup a postgres db on your machine.

### Setup with docker

**Make sure you have docker-compose installed on your machine.**

```
docker-compose up
```

## Environments

We have to environments:

- Staging: `https//staging.hs-flensburg.dev`
- Production: `TBD`
