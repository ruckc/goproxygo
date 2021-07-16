# Usage

This example runs a webserver at `http://localhost:8080` with /api/ mapped to localhost:8081 and / mapped to localhost:8082.  The first matching regular expression gets the request.
```
goproxygo /api/.*:http://localhost:8081 /.*:http://localhost:8082
```

This example runs a webserver at `http://0.0.0.0:8080` with /api/ mapped to localhost:8081 and / mapped to localhost:8082.  The first matching regular expression gets the request.
```
goproxygo --host 0.0.0.0 /api/.*:http://localhost:8081 /.*:http://localhost:8082
```

This example runs a webserver at `http://0.0.0.0:8000` with /api/ mapped to localhost:8081 and / mapped to localhost:8082.  The first matching regular expression gets the request.
```
goproxygo --host 0.0.0.0 --port 8000 /api/.*:http://localhost:8081 /.*:http://localhost:8082
```

# Build Instructions

In order to build, using [godel](https://github.com/palantir/godel).  The binary is built in `out/build/goproxgo/unspecified/linux-amd64/goproxygo`
```bash
./godelw build
```
