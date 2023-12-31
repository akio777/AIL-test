# AIL-test

### About this repo

At the outset, the repository was designed with two architectural plans for these services: plan A and plan B.

For plan A, the plan was to have two subsidiary services:
1. **fetcher** - Responsible for reading on-chain Ethereum data from each UniSwapv3 pool, converting decimal data, calculating what is needed for APY computation, and storing this information in a PostgreSQL database.
2. **api** - This service will encompass the endpoints GET `/apy`, POST `/pool`, and DELETE `/pool`. Each endpoint serves the following purpose:
   2.1 `/apy` - Takes a pool address from the query string to retrieve data from the database and calculate the 24-hour APY.
   2.2 POST, DELETE - These are APIs for CREATING and DELETING data in the database.

Plan B is similar to plan A but with a key difference in planA point 1, instead of reading data directly from the rpc chain, it changes to read data from a graphql endpoint at `https://api.thegraph.com/`.

During initial implementation, there were issues with calculating decimals and the APY formula, which led to the adoption of plan B. Data from thegraph is provided in a more summarized form, making it preferable.

The workflow of these two services is as follows:
- The **fetcher** is tasked with pulling different pool data, sourcing addresses from the database table `pool_address` and using go routines to pull historical data in 30-day increments (to accumulate data for APY calculation), storing the relevant information in the `pool_state` table.
- The **api** (describing only the GET /apy part) will query `pool_state` using the pool address from the query string, requesting data from the past 365 days (364 days before the request date to the request date itself).


### 'localhost:3007' for testing /apy /pool


### FYI : after docker-compose up , please waiting 1-2 minute for setup and fetch data

Steps for setting up docker-compose to run the container repository:
### The `.env` file will include the following
#### .env `api` service
```
API_DB_HOST=0.0.0.0
API_DB_PORT=35437
API_DB_USER=postgres
API_DB_PASSWORD=postgres
API_DB_NAME=ail
API_DB_DB=ail
API_DB_SSL_ENABLE=disable

API_NAME=ail-test
API_PORT=3007
API_HOST=0.0.0.0
API_RPC_URL=https://eth.llamarpc.com
```
#### .env `fetch` service
```
API_DB_HOST=0.0.0.0
API_DB_PORT=35437
API_DB_USER=postgres
API_DB_PASSWORD=postgres
API_DB_NAME=ail
API_DB_DB=ail
API_DB_SSL_ENABLE=disable

API_NAME=ail-test
API_PORT=3008
API_HOST=localhost
API_RPC_URL=https://eth.llamarpc.com
API_SCHEDULE_FETCH_POOL=*/15 * * * * *
API_GRAPHQL_URL=https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3
API_GRAPHQL_READ_FIRST=30
```

### The `docker-compose.yml` will be set up as follows
```
version: '3'
services:
  ali-db:
      image: xakiox/ail-test-db:latest
      command: ["postgres", "-c", "max_connections=2048", "-c", "shared_buffers=1GB"]
      restart: unless-stopped
      environment:
        POSTGRES_USER: postgres
        POSTGRES_PASSWORD: postgres
        POSTGRES_DB: ail
      ports:
        - "35437:5432"
      networks:
        - ail-network
  fetch:
      image: xakiox/ail-test-fetch:latest
      restart: unless-stopped
      depends_on:
        - ali-db
      env_file:
        - ./cmd/fetch/.env
      environment:
        - API_DB_HOST=host.docker.internal
  api:
      image: xakiox/ail-test-api:latest
      restart: unless-stopped
      depends_on:
        - ali-db
        - fetch
      env_file:
        - ./cmd/api/.env
      environment:
        - API_DB_HOST=host.docker.internal
      ports:
        - "3007:3007"
networks:
  ail-network:
    external: false

volumes:
  ail-db:
```
