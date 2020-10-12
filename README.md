# _Liège_

[![build](https://img.shields.io/github/workflow/status/GaelGirodon/liege/CI?style=flat-square)](https://github.com/GaelGirodon/liege/actions)
[![release](https://img.shields.io/github/v/release/GaelGirodon/liege?style=flat-square)](https://github.com/GaelGirodon/liege/releases)
[![license](https://img.shields.io/github/license/GaelGirodon/liege?color=informational&style=flat-square)](https://github.com/GaelGirodon/liege/blob/master/LICENSE)
[![docker image size](https://img.shields.io/docker/image-size/gaelgirodon/liege?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com/r/gaelgirodon/liege)

HTTP stub server from plain files.

## About

_Liège_ is a static file server for tests with some special features
to easily setup a stub server without requiring too much configuration:

- Handle all HTTP methods
  (static file servers usually handle `GET`/`HEAD` methods only)
- Allow customizing the response using the file name
  (no configuration file required)
- Send the request body back in an HTTP header
- Handle index files

## Usage

```shell
liege [flags] <root-dir>
```

### Arguments

| Argument     | Description                        | Environment variable |
| ------------ | ---------------------------------- | -------------------- |
| `<root-dir>` | Path to the server root directory  | `LIEGE_ROOT`         |
| `-p <port>`  | Port to listen on (default `3000`) | `LIEGE_PORT`         |
| `-h`         | Show the help message and exit     |

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
 |-- products/
 |    |-- index.json
 |    |-- 1.json
 |-- admin/
 |    |-- index__403 (empty file)
```

The server will handle the following routes:

| Method | Path                                                       | Response | Served file           |
| ------ | ---------------------------------------------------------- | -------- | --------------------- |
| `*`    | `/products`<br>`/products/index`<br>`/products/index.json` | `200`    | `products/index.json` |
| `*`    | `/products/1`<br>`/products/1.json`                        | `200`    | `products/1.json`     |
| `*`    | `/admin`<br>`/admin/index`                                 | `403`    |                       |

And send the following response:

- **Status code**: `200` (default) or a custom code
- **Headers**:
  - `Content-Type`: determined from file content and extension,
    e.g. `application/json; charset=utf-8`
  - `X-Request-Body`: base64 encoded request body (only if body size <= 4 KB)
- **Body**: stub file contents

The response can be customized using the following file name syntax:

```text
<name>[__<code>][.<ext>]
```

| Param  | Description                                      |
| ------ | ------------------------------------------------ |
| `name` | File name, used in the URL path                  |
| `code` | Custom HTTP response status code (default `200`) |
| `ext`  | File extension, helps to determine content type  |

On start-up, the server loads stub files in memory and build routes.
To reload stub files from the root directory and update routes, call the
`refresh` endpoint.

### Management endpoints

The server provides the following management endpoints:

| Method | Path              | Response | Description          |
| ------ | ----------------- | -------- | -------------------- |
| `GET`  | `/_liege/config`  | `200`    | Get configuration    |
| `PUT`  | `/_liege/config`  | `204`    | Update configuration |
| `POST` | `/_liege/refresh` | `204`    | Reload stub files    |
| `GET`  | `/_liege/routes`  | `200`    | Get available routes |

## License

_Liège_ is licensed under the GNU General Public License.
