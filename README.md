# Gitlab-Classrooms

## Project Structure 

The project has two parts:

- Golang with fiber as backend=> `/`
- React with Typescript and vite as Frontend => `/frontend/`

The frontend proxies the requests for the path `/api/*` to the backend server.

## Development

For development we use the git flow branching model for simplicity. 

### Setup

Copy the `.env.example` file and make your changes:

```
cp .env.example .env
```

#### OAuth with Gitlab
1. We use Gitlab as an OAuth provider, so you have to add this application in your Gitlab.
   * The Redirect URI is for example: https//staging.hs-flensburg.dev/api/auth/gitlab/callback
   * Uncheck Confidential
   * Needed Scopes: api
2. Click on Save application and copy the shown Application ID and Secret to your local .env file

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

### Mail

For local development we use [mailpit](https://mailpit.axllent.org/), running on [localhost:8025](http://localhost:8025).
**This requires the docker setup.**

## Environments

We have to environments:

- Staging: `https//staging.hs-flensburg.dev`
- Production: `TBD`
