# Gitlab-Classrooms

## Project Structure 

The project has two parts:

- Golang with fiber as backend=> `/`
- React with Typescript and vite as Frontend => `/frontend/`

The frontend proxies the requests for the path `/api/*` to the backend server.

## Development

For development, we use the git flow branching model for simplicity.

### Setup
The Setup for development is documented in the the following file
[dev_setup.md](dev_setup.md)

#### Code generation

To generate up to date mock files and database code you can use the command `go generate` in projects root dir.

### Postman testing

1. Login via gitlab in the browser
2. Copy the session_id cookie from the browser to your postman environment
3. Copy the csrf_ cookie from the browser to your postman environment
4. Add the following header to your `POST|PUT|PATCH|DELETE` requests:
    - `X-CSRF-Token: {{csrf_}}`
 

## Environments

We have to environments:

- Staging: `https//staging.hs-flensburg.dev`
- Production: `TBD`
