# Setup project with docker compose
## Requirements:
**Make sure you have docker-compose installed on your machine.**

## Docker compose Setup Information

* Backend server is listening at
    * http://localhost:3000
* Frontend server is listening at
    * http://localhost:5173
* Mail server is listening at
    * http://localhost:8025/

## Step-by-Step Guide to launch the docker application

### Clone the project
1. Open your Terminal and clone the project:
```
git clone https://gitlab.hs-flensburg.de/fb3-masterprojekt-gitlab-classroom/gitlab-classroom.git
```
* When prompted, enter your HS-Flensburg Gitlab credentials.
2. Navigate to the project directory:
```
cd gitlab-classroom
```

### Set up your Enviromets
1. Copy the default Enviromets:
```
cp .env.example .env
```
2. Make your changes in your .env file.
* If running the project locally, adjust the public URL to localhost:
```
PUBLIC_URL=http://localhost:5173
```

### Add the application to your Gitlab Instance
Since we use Gitlab as an OAuth provider, add this application in your Gitlab.
1. Open Gitlab in your browser and navigate to Edit profile.
2. Under Applications:
* Click on "Add new application."
    * Name: e.g.: Gitlab Classroom
    * Redirect URI: The Callback URI for backend and frontend, e.g.:
      http://localhost:3000/auth/gitlab/callback
      http://localhost:5173/auth/gitlab/callback
    * Uncheck "Confidential."
    * Needed Scopes: "api"
    * Save the application and copy the displayed Application ID (AUTH_CLIENT_ID) and Secret (AUTH_CLIENT_SECRET) to your local .env file.

### Install front end dependencies
1. Navigate to the frontend directory:
```
cd frontend
```
2. Install dependencies with:
```
yarn
```
3. Navigate to the application directory:
```
cd ..
```

### Start the Application via Docker Compose
```
docker-compose up -d
```
### For encrypted connections, generate a self-signed certificate
```
openssl req -x509 -newkey rsa:4096 -nodes -keyout .docker/mail/privkey.pem -out .docker/mail/cert.pem -sha256 -days 3650
```

### Restart Containers
```
docker-compose restart
```

### To exit the application call
```
docker-compose down
```
