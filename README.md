# Gitlab-Classrooms

## Project Structure 

The project has two parts:

- Golang with fiber as backend=> `/`
- React with Typescript and vite as Frontend => `/frontend/`

The frontend proxies the requests for the path `/api/*` to the backend server.

## Development

For development, we use the git flow branching model for simplicity.

### Setup
The Setup for development is documented in the following file
[dev_setup.md](docs/dev_setup.md)


#### Code Generation

The following file-types in our repo are auto-generated:
- mocks
- gorm-gen queries
- swagger documentation
- swagger-client

The first three types are generated when running `go generate` which is automatically executed before each build when running the app through air.

To generate the frontend `swagger-client` after some changes to the docs the following script needs to be executed:

```bash
# Linux / Mac / WSL
./script/swagger-codegen.sh
# Windows
./script/swagger-codegen.ps1
```

### Run Staging|Production image locally

To test the completly build image locally you can use the following command:

```bash
docker compose -f docker-compose.dev.yml up --build
```

And access the application via `http://localhost:3000`

### Swagger-UI testing (recommended)

1. Run the app described in [Dev-Setup](docs/dev_setup.md#start-the-application)
2. Add [Swag](https://github.com/swaggo/swag/?tab=readme-ov-file#api-operation) comments directly above your endpoint
3. Visit the page `http://localhost:5173/docs.html`
4. Click on the `Sign-In`-Button on the top-right
5. Get your csrf-token
6. Test your endpoints
7. Refresh page after swag-comments were edited

### Postman testing (not recommended)

1. Login via gitlab in the browser
2. Copy the session_id cookie from the browser to your postman environment
3. Copy the csrf_ cookie from the browser to your postman environment
4. Add the following header to your `POST|PUT|PATCH|DELETE` requests:
    - `X-CSRF-Token: {{csrf_}}`

### Writing Tests

...

To run the test simply exec the following script:
```bash
# Linux / Mac / WSL
./script/test_backend.sh one|air
./script/test_frontend.sh
# Windows
./script/test_backend.ps1 one|air
./script/test_frontend.ps1
```

## Environments

We have to environments:

- Staging: `https//staging.hs-flensburg.dev`
- Production: `TBD`
