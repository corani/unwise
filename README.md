# A partial reimplementation of the Readwise API in Go

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

I'm using the [Obsidian Readwise (Community Plugin)](https://github.com/renehernandez/obsidian-readwise) 
to import my highlights into Obsidian. To allow for setting a custom API server location, use the 
following patch: [Add setting for API server location](https://github.com/algocentric/obsidian-readwise/commit/f2da99bd9d387536171a1ed37217c5548b236ee4)

## Running locally

Run the following command:

```sh
$ go run ./cmd/unwise/

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

## Configuration 

The following environment variables can be used to configure the server (you can also add them to
a `.env` file): 

| Variable    | Description                   | Default     |
| ----------- | ----------------------------- | ----------- |
| `LOGLEVEL`  | Log Level                     | `info`      |
| `REST_ADDR` | Address to listen on          | `:3123`     |
| `REST_PATH` | Base path to listen on        | `/api/v2`   |
| `TOKEN`     | Authentication token          | (generated) |

## Endpoints

| Method | Path                 | Description                   | Used by      |
| ------ | -------------------- | ----------------------------- | ------------ |
| GET    | `/api/v2/auth`       | Validate authentication token |              |
| POST   | `/api/v2/highlights` | Create a highlight            | Moon+ Reader |
| GET    | `/api/v2/highlights` | Get all highlights            | Obsidian     |
| GET    | `/api/v2/books`      | Get all highlights            | Obsidian     |

## Docker 

TODO
