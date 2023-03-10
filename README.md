# Golang + Echo + Gorm + Nodemon + PostgreSQL + Swagger

Requirement :
- Go 1.16
- Node.js
- PostgreSQL

# Start π

## Install Modules

```bash
go mod download && go mod tidy && go mod verify
```

If the message below was shown, do the next step.
```
go: finding module for package github.com/forkyid/go-rest-api/docs
github.com/forkyid/go-rest-api/src/route imports
        github.com/forkyid/go-rest-api/docs: no matching versions for query "latest"
```

## Swagger Installation and Swag Initialization

```bash
go install github.com/swaggo/swag/cmd/swag@v1.6.7
```

```bash
swag init -g src/main.go
```

```bash
go mod tidy
```

## Install Nodemon

```bash
npm install
```

or

```bash
npm install -g nodemon
```

## Running the Server

### Go run + Nodemon

```bash
nodemon --exec go run src/main.go --signal SIGTERM
```

Swagger API Documentation URL:
```url
http://localhost:5000/swagger/index.html#/
```

![1](/images/1.png)
![2](/images/2.png)

### Docker

```bash
docker-compose up --build
```

## Repository Structure

```bash
.
βββ .github
β   βββ PULL_REQUEST_TEMPLATE.md
βββ database-migrations
β   βββexamples
β   βββREADME.md
βββ src
β   βββ connection
β   β   βββ connection.go
β   βββ constant
β   β   βββ constant.go
β   βββ controller
β   β   βββ v1
β   β       βββ account
β   β       β   βββ account.go
β   β       βββ auth
β   β       β   βββ auth.go
β   β       βββ location
β   β           βββ location.go
β   βββ http
β   β   βββ account.go
β   β   βββ auth.go
β   β   βββ location.go
β   βββ model
β   β   βββ account.go
β   β   βββ location.go
β   βββ pkg
β   β   βββ bcrypt
β   β   β   βββ bcrypt.go
β   β   βββ jwt
β   β       βββ jwt.go
β   βββ repository
β   β   βββ v1
β   β       βββ account
β   β       β   βββ account.go
β   β       βββ location
β   β           βββ location.go
β   βββ routes
β   β   βββ main.go
β   βββ service
β   β   βββ v1
β   β       βββ account
β   β       β   βββ account.go
β   β       βββ location
β   β           βββ location.go
β   βββ main.go
βββ .env.example
βββ .gitignore
βββ go.mod
βββ LICENSE
βββ README.md
```
