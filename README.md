# resource-tool

Resource-tool is a library that wraps basic functions over mux router to ease generation of fields
and collections on postman.

submission for postman hackathon

## How to run it

```
    cd api-mgr/cmd
    go run main.go
```

### .env in api-mgr/cmd/

```
    POSTMAN_TOKEN=`provide your postman token`
    URI=https://api.getpostman.com
```

simple tool for visualizing what endpoints services are referencing within your application.

### TO Access UI

```
    cd web
    npm install
    npm run dev
```

## How to use
### Step One : Instantiate a new manager 

When calling the mgr.New() function in postmanager pass in the endpoint that the resource tool will be running on. In this instance: "http://localhost:9999"


```go
    mgr := apimgr.New(resource, "service-b")
```

This returns a new manager instance with the instantiated name of the collection or ID to make it easier 
to index the Endpoints that this particular service will have access to when they call mgr.Use function later on.
The instance keeps track of a lot of things. The resource in this instance in the endpoint that is running the resource-tool
API that keeps track of all changes and interacts with postman.

### Step Two: 

```go
    mgr.NewEnv("service-b").EnvVar("service_b", "http://localhost:7767")
```

Create  a new postman Environment that is going to hold the variables needed for this service on Postman. The base url for 
instance and any other variables you could store on postman.

### Step Three

```go
    c := mgr.NewCollection("service-b", apimgr.DefaultCollOpts())
```

Create a collection that is going to hold all the endpoints that you define on the returned c variable. The new collection returns a struct that has inherited functions in the *mux.Router Library. Makes doing the below possible.

### Step Three 

Define your routes and there respective Handler functions that will be called when the router matches it. This HandlerFunc function provided is a light wrapper over mux.Router.HandlerFunc() to allow passing a lot more varibales than would be allowed on
the latter. 

```go
    c.HandleFunc("health-check", "{{service_b}}", "/health", "POST", HealthHandler, TestPayload{})
	c.HandleFunc("test", "{{service-b-endpoint}}", "/test", "POST", TestHandler, TestNested{}).Header("Authorization", "Bearer {{token_special}}")
```


### Step Four 

Bundle all the requests under one folder and call create which runs through all the requests in the particular collection and
returns a proper formatted collection .

```go
    requests := c.Requests("samples")
	err = mgr.Create(&requests)
	if err != nil {
		log.Printf("req creation failed %s", err)
    }
```

### Step Five

Call the APIs or endpoints you would like access to using the mgr.Use function. Current implementation references the server for each request . Plan to implement an inmemory cache and use an event store to propagate changes to listening clients. i.e how spring cloud config propagates changes

```go
	_, err = mgr.Use("service-a.service_a", "bye-endpoint")
	if err != nil {
		log.Printf("resource is unavailable %s", err)
	}

	_, err = mgr.Use("service-c.service_c", "c-endpoint")
	if err != nil {
		log.Printf("resource is unavailable %s", err)
	}
```

The first parameter is a "." separated string of the environment you would like to access and the collection you would want access to. The second parameter is the name of the request as is viewed on postman . 


The resource tool comes with a small ui that allows you to update env variables and to view the relationships and dependencies 
in your services in realtime and the services they are accessing. 

More examples in the examples folder
