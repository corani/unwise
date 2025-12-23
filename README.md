# A partial reimplementation of the Readwise API in Go

[![codecov](https://codecov.io/gh/corani/unwise/graph/badge.svg?token=2FGRUHZ1B5)](https://codecov.io/gh/corani/unwise)

Rather than manually exporting my book notes from Moon+ Reader and importing them into Obsidian,
I'd like to automate the process. Readwise has an [API](https://readwise.io/api_deets) that can be
used for this purpose, that is supported by both Moon+ Reader and Obsidian, however I'd like to
keep my data private. So this is a partial reimplementation of the Readwise API in Go, supporting
just the endpoints I need.

1. Only a few of the endpoints are implemented.
2. Only a few of the fields are supported.
3. Only category "book" is supported.

## Moon+ Reader

Open any book and go to "Bookmarks". Click on the settings icon at the bottom and enable "Share new
highlights and notes to Readwise automatically". In the settings, enter the Token and server URL.

## Obsidian

I'm using my own fork of the the [Obsidian Readwise (Community Plugin)](https://github.com/renehernandez/obsidian-readwise)
to import my highlights into Obsidian: [corani/obsidian-readwise](https://github.com/corani/obsidian-readwise). This allows
to set the API server location in the settings and supports the "chapter" property for highlights.

You can install the plugin via [BRAT](https://tfthacker.com/BRAT).

## Endpoints

| Method | Path                 | Description                   | Used by      |
| ------ | -------------------- | ----------------------------- | ------------ |
| GET    | `/api/v2/auth`       | Validate authentication token |              |
| POST   | `/api/v2/highlights` | Create a highlight            | Moon+ Reader |
| GET    | `/api/v2/highlights` | Get all highlights            | Obsidian     |
| GET    | `/api/v2/books`      | Get all highlights            | Obsidian     |

## Running locally

Run the following command:

```sh
$ ./build.sh -b
[INFO] Building unwise version dev/371e928027a8c9f03dccf5b59acd07640c52e4ea
[CMD ] go build -o bin/unwise ./cmd/unwise/
[TIME] took 0m0.186s

$ ./bin/unwise

16:13:41 INFO generated new token token=68ef286d-c71a-4225-aba4-1e4cd6633fc4

 ┌───────────────────────────────────────────────────┐
 │                   Fiber v2.51.0                   │
 │               http://127.0.0.1:3123               │
 │       (bound on host 0.0.0.0 and port 3123)       │
 │                                                   │
 │ Handlers ............ 10  Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............. 28732 │
 └───────────────────────────────────────────────────┘
```

### Configuration

The following environment variables can be used to configure the server (you can also add them to
a `.env` file):

| Variable    | Description                   | Default     |
| ----------- | ----------------------------- | ----------- |
| `LOGLEVEL`  | Log Level                     | `info`      |
| `REST_ADDR` | Address to listen on          | `:3123`     |
| `REST_PATH` | Base path to listen on        | `/api/v2`   |
| `DATA_PATH` | Path to store data            | `/tmp`      |
| `TOKEN`     | Authentication token          | (generated) |

Note: if you don't provide a `TOKEN`, the application will generate one and print it to the
console during startup.

## Docker

```sh
$ docker run --rm -it -p 3123:3123  \
    -e TOKEN=my-token               \
    ghcr.io/corani/unwise:latest
```

Or using docker compose:

```sh
$ docker-compose -f docker/docker-compose.yml up
...
```

## Traefik

I'm running the app behind Traefik, so I can use Let's Encrypt for SSL. Here is an example:

```yaml
services:
  unwise:
    image: "ghcr.io/corani/unwise:latest"
    container_name: "unwise"
    user: "${MY_UID}:${MY_GID}"
    env_file:
      - "./unwise/unwise.env"
    volumes:
      - "./unwise/data:/data"
    restart: "unless-stopped"
    networks:
      - "proxy"
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=proxy"
      # Replace with your own domain name.
      - "traefik.http.routers.unwise.rule=Host(`unwise.example.com`)"
      # The 'websecure' entryPoint is basically your HTTPS entrypoint.
      - "traefik.http.routers.unwise.entrypoints=websecure"
      - "traefik.http.routers.unwise.service=unwise"
      - "traefik.http.services.unwise.loadbalancer.server.port=3123"
      - "traefik.http.routers.unwise.tls=true"
      # Replace the string 'letsencrypt' with your own certificate resolver
      - "traefik.http.routers.unwise.tls.certresolver=letsencrypt"
      - "traefik.http.routers.unwise.middlewares=unwisecors"
      # The part needed for CORS to work on Traefik 2.x starts here
      - "traefik.http.middlewares.unwisecors.headers.accesscontrolallowmethods=GET,PUT,POST,HEAD,DELETE"
      - "traefik.http.middlewares.unwisecors.headers.accesscontrolallowheaders=accept,authorization,content-type,origin,referer"
      - "traefik.http.middlewares.unwisecors.headers.accesscontrolalloworiginlist=app://obsidian.md,capacitor://localhost,http://localhost"
      - "traefik.http.middlewares.unwisecors.headers.accesscontrolmaxage=3600"
      - "traefik.http.middlewares.unwisecors.headers.addvaryheader=true"
      - "traefik.http.middlewares.unwisecors.headers.accessControlAllowCredentials=true"
```

## Persistence

The application uses a local SQLite database (`/data/unwise.db`) to persist the data.
