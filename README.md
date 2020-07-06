## opa-server

Start the OPA server:
```bash
$ docker-compose up -d
```

In the `query` textarea, add this:
```
isAlice := input.name == "alice"
isBob := input.name == "bob"
```

In the `input` textarea, add this:
```
{
  "name": "alice"
}
```

Click submit:
```json
{
  "result": [
    {
      "isAlice": true,
      "isBob": false
    }
  ]
}
```

## Loading bundles

We set the docker volumes to load bundles from the `example` namespace:

```bash
http://localhost:8181/v1/data/example/greeting
http://localhost:8181/v1/data/example/math
```

If we perform the following `query`:
```
math := data.example.math
greeting := data.example.greeting
```
You should get the output:
```json
{
  "result": [
    {
      "greeting": "hello from container \"e7fb3cfe71b1\"!",
      "math": 1
    }
  ]
}
```

## Example api.rego

`query` textarea:
```bash
api := data.example.allow
is_get := data.example.is_get
```

`input` textarea (NOTE: If we change path `john` to `jane`, it does not produce any more output. Not sure if it is by design or it's a bug):
```bash
{
  "method": "GET",
  "path": ["users", "john"],
  "user": "john"
}
```

`output` textarea:
```
{
  "result": [
    {
      "api": true,
      "is_get": true
    }
  ]
}
```

## Example external data

See `example/external.rego`. In this policy, we are actually making a REST call to an external API `/admins` endpoint to check the name.


We can choose to either run both the OPA server and API server locally, or run them in containers. Note that if we run the OPA server in container, and the go server locally, in order for the OPA server to call the go server (which is running in the host), it needs to call `http://host.docker.internal:8080/admins` instead (A simple acronym is **HDI**).

For our `query`:
```
name := data.example.admin == true
```

For out `input`:
```json
{
  "name": "alice"
}
```

The output would be:
```json
{
  "result": [
    {
      "name": true
    }
  ]
}
```

If we change the name to `john` instead, the output would be false.
