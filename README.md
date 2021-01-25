# resource-tool

submission for postman hackathon

## How to run it

cd api-mgr/cmd

```go
    go run main.go
```

```
    POSTMAN_TOKEN=`provide your postman token`
    URI=https://api.getpostman.com
```

simple tool for visualizing what endpoints services are referencing within your application.
more documentation on: (postmanager)[https://github.com/mar-tina/postmanager/edit/main/README.md]

When calling the mgr.New() function in postmanager pass in the endpoint that the resource tool will be running on. In this instance: "http://localhost:9999"

