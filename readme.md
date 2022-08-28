## Proxy service

### Usage

`docker compose up -d` - command to start application in docker space

OR

`go build -o {output place} cmd/proxy/proxy.go` - to build project in local machine

You can configure application with env variables:

* `APP_PORT` - http listen port. Default is `":8080"`
* `SRV_READ_TIMEOUT` - http server read timeout. Default is 5 second
* `CLIENT_TIMEOUT` - http client timeout. Default is 15 second
* `SHUTDOWN_TIMEOUT` - shutdown timeout. Used for graceful shutdown. Default is 30 second.

### Application

Http server have one `"/"` endpoint to proxy. Example `http://localhost:8080/`.
THe input/output data is must be JSON. 

Input JSON Format
```json
{
  "url": "{resourse url}",
  "method": "{http verb}",
  "headers": {
    "User-Agent": "{map of request headers}"
  },
  "body": "{base64 encoded body to request}"
}
```

Output JSON Format
```json

{
  "id": "{generated id to request in GUID format}",
  "status": "{response status}",
  "headers": {
    "Accept-Ch": "{map of response headers}"
  },
  "length": "{length of body}",
  "body": "{base64 encoded response body}"
}
```
