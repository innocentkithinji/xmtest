# XM Company API Test

This is an service that creates and updates company data via http REST api. It follows the
clean architecture and uses the echo framework.

## Running the API

### Running in docker

#### Prerequisites

You need the following to be installed

- Docker
- Docker-Compose

Run the following command

 ```bash
  docker-compose up --build -d
 ```

The API can be accessed via localhost:8080.

### Running locally

You need atleast `go:1.19` installed

- first install dependencies needed
    ```sh
       go mod download
    ```
- Setup the needed infrastructure

  You need to have a running **mongodb** server and a **kafka** client running.

- Build the server
    ```shell
       go build main -o 
    ``` 
- Run the service
    ```shell
     chmod +x ./main
     ./main
    ```

## Env Variables

|variable| default                                                               |
|---|-----------------------------------------------------------------------|
| XMC_MONGODB_URI | "mongodb://root:toor123@localhost:27017/?retryWrites=true&w=majority" |
| XMC_DB_NAME | xmtest                                                                |
| XMC_COMPANY_COLLECTION | company                                                               |
| XMC_USERS_COLLECTION | users                                                                 |
| XMC_SIGNING_SECRET | secret                                                                |
| XMC_KAFKA_URI | localhost:9092                                                        |
| COMPANY_EVENT_TOPIC | company                                                               |
