# irohasu_jaspr

A new Jaspr project

## Running the project

Copy the environment template and set the backend domain:

```sh
cp .env.example .env
```

Run your project with the environment file:

```sh
jaspr serve --dart-define-from-file=.env
```

The development server will be available on `http://localhost:8080`.

`API_BASE_URL` is the internal URL used by the Jaspr server to call Gin.
`PUBLIC_API_BASE_URL` is the browser-visible URL used by search and article
assets. Use `/` behind the repository's Nginx proxy.

## Building the project

Build your project using:

```sh
jaspr build --dart-define-from-file=.env
```

The server executable and web assets will be located inside `build/jaspr/`.
