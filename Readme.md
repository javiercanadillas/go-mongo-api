
## Introduction

This is a sample MongoDB web API written with Go and [Gin-Gonic](https://github.com/gin-gonic/gin). The app is using a `library` database in a MongoDB instance with a particular collection already imported.

The app tries first to obtain the mongoDB connection string from Google Cloud Secret Manager, and if it fails, it'll then look for the environment variable `MONGODB_URI`.

## Installation

- Clone this repo
- Setup a Mongo Database instance. Create a database called `library` and import the `extras/books.json` collection there.
- Create a MongoDB user in the `library` database with `readWrite` permissions
- Store the Mongo DB connection string in Secret Manager
- Configure a Service Account in the service where the application is running and grant it the `roles.SecretAccessor` IAM role.
- Generate an image
    - `docker build -t mongoApp .`
    - `docker run -t mongoApp`

## TODO

- Automation for MongoDB GCE instance creation
- Automation for Service Accounts and Secret Manager setup
