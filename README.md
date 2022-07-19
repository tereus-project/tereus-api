# tereus-api

## Requirements

- Docker
- Docker Compose
- Git
- Node.js v16.x

## Setting up the API and storage services

The `docker-compose.yml` file inside the repo contains the following services:

- `api`: The API service
- `postgres`: The PostgreSQL service for database
- `minio`: The MinIO service for object storage
- `nsqd`: The NSQ service for queuing

Setting up the API with the services is done by running the following commands:

```sh
git clone git@github.com:tereus-project/tereus-api.git
cd tereus-api
cp env.example .env
docker-compose up -d
```

The docker-compose services, including the API, use environment variables to configure themselves. The `.env.example` file contains the environment variables that you need to set for a local environment, but you can modify the file to your own needs.

The docker containers should look like this:

```sh
~/tereus/tereus-api ‹main› » docker-compose ps
         Name                        Command                  State                                               Ports
----------------------------------------------------------------------------------------------------------------------------------------------------------------
tereus-api_api_1          /go/bin/air                      Up             0.0.0.0:1323->1323/tcp
tereus-api_minio_1        /usr/bin/docker-entrypoint ...   Up (healthy)   0.0.0.0:9000->9000/tcp, 0.0.0.0:9001->9001/tcp
tereus-api_nsqadmin_1     /nsqadmin --lookupd-http-a ...   Up             4150/tcp, 4151/tcp, 4160/tcp, 4161/tcp, 4170/tcp, 0.0.0.0:4171->4171/tcp
tereus-api_nsqd_1         /nsqd --lookupd-tcp-addres ...   Up             0.0.0.0:4150->4150/tcp, 0.0.0.0:4151->4151/tcp, 4160/tcp, 4161/tcp, 4170/tcp, 4171/tcp
tereus-api_nsqlookupd_1   /nsqlookupd                      Up             4150/tcp, 4151/tcp, 0.0.0.0:4160->4160/tcp, 0.0.0.0:4161->4161/tcp, 4170/tcp, 4171/tcp
tereus-api_postgres_1     docker-entrypoint.sh postgres    Up             0.0.0.0:5432->5432/tcp
```
