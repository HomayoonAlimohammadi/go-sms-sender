# SMS Sender with Golang

This app serves you as an SMS Sender to easily send any message
from your own number to other peers. Also you can search through records of the
messaged sent from a given number

## How to use

First of all you need to have `.env` and `config.yaml` files for the
configuration and environment varialbles.

- Both files should be located in the root of the project (where this `README.md` exists).

Here is a template for both the `config.yaml` and `.env` file:

- .cobra.yaml

```yaml
SmsSender:
  ApiKey: "<kavenegar_api_key>"
  Port: "8000"
  Host: localhost
  GracefulTimeout: 15000000000

PostgresDB:
  Host: "postgres"
  Port: "5432"
  User: "postgres"
  Password: "postgres"
  DbName: "smssender"
  SslMode: "disable"
  MigrationsPath: "file://internal/database/migrations"
```

- .env

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=smssender
POSTGRES_SSL_MODE=disable
```

If you want to use this without `Dockerfile` or `docker-compose`, don't forget to change the `postgres.host` to `localhost` or your own postgres remote host.

Run the application by:

```bash
make serve
```

This should take care of downloading the dependencies and starting the project

Or if you decided to run the project in a containerized environment, run:

```bash
docker compose up --build
```

Which should run the application on `port 8000` alongside a postgres container.

## Enjoy!
