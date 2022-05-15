# _Liège_

[![release](https://img.shields.io/github/v/release/GaelGirodon/liege?style=flat-square)](https://github.com/GaelGirodon/liege/releases)
[![license](https://img.shields.io/github/license/GaelGirodon/liege?color=informational&style=flat-square)](https://github.com/GaelGirodon/liege/blob/master/LICENSE)
[![build](https://img.shields.io/gitlab/pipeline/GaelGirodon/liege/master?style=flat-square)](https://gitlab.com/GaelGirodon/liege/-/pipelines/latest)
[![coverage](https://img.shields.io/gitlab/coverage/GaelGirodon/liege/master?style=flat-square)](https://gitlab.com/GaelGirodon/liege/-/pipelines/latest)
[![docker](https://img.shields.io/docker/image-size/gaelgirodon/liege?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com/r/gaelgirodon/liege)

HTTP stub server from plain files

## About

_Liège_ is a static file server for tests with some special features
to easily setup a stub server without requiring too much configuration:

- Handle all HTTP methods
  (static file servers usually handle `GET`/`HEAD` methods only)
- Allow customizing request routing and response using the file name without
  requiring a configuration file
- Send the request body back in an HTTP header
- Handle index files

## Usage

```shell
liege [flags] <root-dir>
```

### Arguments

| Argument     | Description                          | Environment variable | Configuration |
| ------------ | ------------------------------------ | -------------------- | ------------- |
| `<root-dir>` | Path to the server root directory    | `LIEGE_ROOT`         | `root`        |
| `-p <port>`  | Port to listen on (default `3000`)   | `LIEGE_PORT`         |
| `-c <cert>`  | Path to the TLS certificate PEM file | `LIEGE_CERT`         |
| `-k <key>`   | Path to the TLS private key PEM file | `LIEGE_KEY`          |
| `-l <lat>`   | Simulated response latency in ms     | `LIEGE_LATENCY`      | `latency`     |
| `-v`         | Print the version number and exit    |
| `-h`         | Print the help message and exit      |

### Example

```shell
$ liege ./data/
_________ __   _________________
________ / /  /  _/ __/ ___/ __/
_______ / /___/ // _// (_ / _/
______ /____/___/___/\___/___/

HTTP server started on port 3000
```

### Configuration

Given the following file tree:

```text
data/ -> server root directory
 |-- items/
 |    |-- index.json
 |    |-- index__qs.json
 |    |-- 1__GET.json
 |-- admin/
 |    |-- index__403_l50 (empty file)
```

The server will handle the following routes:

| Method | Paths                                             | Query   | Response | Served file            | Latency |
| ------ | ------------------------------------------------- | ------- | -------- | ---------------------- | ------- |
| `*`    | `/items`<br>`/items/index`<br>`/items/index.json` |         | `200`    | `items/index.json`     | ~ 0 ms  |
| `*`    | `/items`<br>`/items/index`<br>`/items/index.json` | `s[=*]` | `200`    | `items/index__qs.json` | ~ 0 ms  |
| `GET`  | `/items/1`<br>`/items/1.json`                     |         | `200`    | `items/1__GET.json`    | ~ 0 ms  |
| `*`    | `/admin`<br>`/admin/index`                        |         | `403`    |                        | ~ 50 ms |

And send the following response with an optional additional latency:

- **Status code**: `200` (default) or a custom code
- **Headers**:
  - `Content-Type`: determined from file content and extension,
    e.g. `application/json; charset=utf-8`
  - `X-Request-Body`: base64 encoded request body (only if body size <= 4 KB)
- **Body**: stub file contents

Routing and response can be customized using the following file name syntax:

```text
<path>[__<options>][.<ext>]
```

| Param     | Description                                     |
| --------- | ----------------------------------------------- |
| `path`    | File name, used in the URL path                 |
| `options` | Additional routing and response configuration   |
| `ext`     | File extension, helps to determine content type |

`__<options>` can be used to further customize request routing and response by
appending a list of options, prefixed by `__` and separated by `_`, at the end
of the file name:

| Syntax           | Description                      | Default | Examples        |
| ---------------- | -------------------------------- | ------- | --------------- |
| `<method>`       | HTTP method                      | `*`     | `GET`           |
| `q<key>[=<val>]` | Required query parameter(s)      |         | `qerror=1`      |
| `<code>`         | Custom HTTP response status code | `200`   | `401`           |
| `l<x>[-<y>]`     | Simulated response latency in ms | `0`     | `l40`, `l50-90` |

For example, the content of a file named `page__GET_qsearch_403_l250` will be
sent with a `403` status code and at least 250 ms latency only for `GET`
requests on `/page` URL with a `search` query parameter.

The latency can be constant (`<x>`) or random between a range (`<x>-<y>`) and
can be defined globally (using the CLI or the environment variable) and at the
route level using the file name. The same syntax (`<x>[-<y>]`) is used in both
cases. Latency defined at the file level overrides globally defined latency
unless the latter is set to `-1` which totally disables latency.

On start-up, the server loads stub files in memory and build routes. To reload
stub files from the root directory and update routes, call the
`refresh` endpoint.

### Management endpoints

The server provides the following management endpoints:

| Method | Path              | Response | Description          |
| ------ | ----------------- | -------- | -------------------- |
| `GET`  | `/_liege/config`  | `200`    | Get configuration    |
| `PUT`  | `/_liege/config`  | `204`    | Update configuration |
| `POST` | `/_liege/refresh` | `204`    | Reload stub files    |
| `GET`  | `/_liege/routes`  | `200`    | Get available routes |

### TLS setup

Generate a self-signed X.509 TLS certificate or obtain a certificate from a CA,
and start the server with `-c` and `-k` flags:

```shell
$ liege -c cert.pem -k key.pem ./data/
```

## License

_Liège_ is licensed under the GNU General Public License.
