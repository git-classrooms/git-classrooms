# Setup project with docker compose

## Windows WSL
If you want to use Docker Desktop you need activate wsl on your machine.
1. Open PowerShell as Administrator and run:
```powershell
wsl --install [-d Ubuntu-22.04]
```
2. [Install Docker Desktop](https://docs.docker.com/desktop/install/windows-install/)
3. (Optinal)Enable WSL Integration for Ubuntu-22.04 in Docker Desktop under `Settings -> Resources -> WSL Integration`

### Develop in WSL (recommended)
1. Open PowerShell and run:
```powershell
wsl
```
2. check if you are in the right WSL distro
```bash
cat /etc/os-release
```
3. Check if the docker command is available
```bash
docker run --rm hello-world
```
#### Connect your IDE to the WSL distro

#### Visual Studio Code
1. Install the [WSL extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl)
2. Open View -> Command Palette and exec command: `WSL: Connect to WSL`
3. Open your desired folder in the WSL distro

#### GoLand | Idea
1. You need to already cloned this repo in your WSL distro
```bash
wsl
cd /home/<your-username>
git clone https://gitlab.hs-flensburg.de/fb3-masterprojekt-gitlab-classroom/gitlab-classroom.git
cd gitlab-classroom
```
2. Open the cloned folder in your JetBrains IDE like `\\wsl$\Ubuntu-22.04\home\<your-username>\gitlab-classroom`
3. Go to the top right corner and click on settings and check the following things
    - <Go | Languages & Frameworks -> Go> -> GOROOT like `\\wsl$\Ubuntu-22.04\usr\local\go`
    - <Go | Languages & Frameworks -> Go> -> GOPATH like `\\wsl$\Ubuntu-22.04\home\<your-username>\go`
    - Tools -> Terminal -> Starting directory like `\\wsl$\Ubuntu-22.04\home\<your-username>\gitlab-classroom`
    - Tools -> Terminal -> Shell path like `wsl.exe --distribution Ubuntu-22.04`
    - Languages & Frameworks -> Node.js -> Node interpreter like `\\wsl$\Ubuntu-22.04\...`
    - Languages & Frameworks -> Node.js -> Package manager like `yarn \\wsl$\Ubuntu-22.04\...`
    - Languages & Frameworks -> Typescript -> like `\\wsl$\Ubuntu-22.04\...`


## Requirements:
**Make sure you have docker-compose installed on your machine.**
- [Docker Desktop on Mac and Windows](https://docs.docker.com/desktop/install/)
- Docker on linux


## Docker compose Setup Information

* Backend server is listening at
    * http://localhost:3000
* Frontend server is listening at
    * http://localhost:5173
* Mail server is listening at
    * http://localhost:8025/
* Postgres server is listening at
    * localhost:5432

## Step-by-Step Guide to launch the docker application

### Clone the project
1. Open your Terminal and clone the project:
```
git clone https://github.com/git-classrooms/git-classrooms.git
```
* When prompted, enter your HS-Flensburg Gitlab credentials.
2. Navigate to the project directory:
```
cd git-classrooms
```

### Set up your Environment
1. Copy the default Environment:
```
cp .env.example .env
```
2. Make your changes in your .env file.
* If running the project locally, adjust the public URL to localhost:
```
PUBLIC_URL=http://localhost:5173
POSTGRES_HOST=<localhost(local) | postgres(docker)>
SMTP_HOST=<localhost(local) | mail(docker)>
```

### Add the application to your Gitlab Instance
Since we use Gitlab as an OAuth provider, add this application in your Gitlab.
1. Open Gitlab in your browser and navigate to Edit profile.
2. Under Applications:
* Click on "Add new application."
    * Name: e.g.: GitClassrooms
    * Redirect URI: The Callback URI for backend and frontend, e.g.:
      http://localhost:3000/api/v1/auth/gitlab/callback
      http://localhost:5173/api/v1/auth/gitlab/callback
    * Check "Confidential."
    * Needed Scopes: "api"
    * Save the application and copy the displayed Application ID (AUTH_CLIENT_ID) and Secret (AUTH_CLIENT_SECRET) to your local .env file.

### For the mail-server, generate a self-signed certificate
```
openssl req -x509 -newkey rsa:4096 -nodes -keyout .docker/mail/privkey.pem -out .docker/mail/cert.pem -sha256 -days 3650
```

### Start the application

#### Local Development

Run the app directly on your machine and not within docker

##### Prerequisites
- Node.js
- Yarn
- Docker
- Docker Compose
- go

##### Install air and other tools for hot reloading of the backend

```bash
go install github.com/cosmtrek/air@latest
go install github.com/vektra/mockery/v2@v2.42.2
go install github.com/swaggo/swag/cmd/swag@latest
```

#### Start the Application via Script
##### Linux | Mac | WSL

```bash
./script/dev.sh
or
zsh ./script/dev.sh
or
bash ./script/dev.sh
```

##### Windows
```poweshell
.\script\dev.ps1
```
