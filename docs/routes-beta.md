# [BETA] Routes

Inspr allows communication via routes, that is, via http calls between nodes. In a simple and quick way, it is possible to register endpoints and send requests from one dApp to another dApp. This feature is still in the beta phase, but the main features are already up and running. This step of development of the routes feature includes the essential ---- of the routes structure. In this beta, you can register handler functions in pre-determined paths, that should be intantiate when you install the dapp on Inspr. When you're writing the dApp yaml, you can include the endpoints you need for each node you want. An declared endpoint will be available for all the other nodes inside the same parent dApp. An endpoint to be available doesn't mean it exists in the node, so you still need to create it. For that, Inspr have a route client library, that handle paths and send request in a simple way. Lets walkthrough an example to clarify to make it clear.

## Declaring the endpoints in the yaml

So the first thing you need to do is decide what endpoints your apps will have. In this example, we will have two nodes dapps that will comunicate between them. Those nodes will have a parent dapp, called `router`. Inside of it we declare both the `api` and the `client` dapps. Those dapps (the `api` and the `client`) are node dapps, therefore it can't have any children. The idea here is that the `api` will register some endpoints, while the `client` will automatically send requests to those endpoints, and then receiving a response from the `api`. For this example, the `api` will declare three endpoints:

- `add`, that adds two numbers received in the request body.
- `sub`, that substract two numbers received in the request body.
- `mul`, that multiply two numbers received in the request body (as you already guess it).

As the client will only send requests, no endpoint will be declared for it. Simple enough. The final yaml will look like this:

```yaml
apiVersion: v1
kind: dapp

meta:
  name: router
spec:
  apps:
    client:
      meta:
        name: client
        parent: router
      spec:
        node:
          spec:
            replicas: 1
            image: <client_app_image>
    api:
      meta:
        name: api
        parent: router
      spec:
        node:
          spec:
            replicas: 1
            image: <api_app_image>
            endpoints:
              - add
              - sub
              - mul

```
Note that the images fields of both the `client` and the `api` need to be filled, as is the code that will run once the dApp is up in the cluster. In the next section, we'll write this code using Inspr's client library, that can register/listen specific endpoints and send requests to another node.

## Using Inspr's Route Client

The Inspr Client makes available two functions related to the routes feature:

```go
func SendRequest(ctx context.Context, nodeName, path, method string, body interface{}, responsePtr interface{})
```
The SendRequest function will send a http request to the node passed as parameter. The path that will be reached is also passed as parameter and, if the path exists, the function will return the response in the responsePtr parameter. In our example, we created models to define the request and response structures. The response one is used in the responsePtr, and the request as the body:

`models.go`
```go
type Request struct {
	Op1 int `json:"op1"`
	Op2 int `json:"op2"`
}

type Response struct {
	Result int `json:"result"`
}
```

`client.go`
```go
var resp model.Response

req := model.Request{
    Op1: 1, // Any number
    Op2: 2, // Any number
}

// ... 

client := insprclient.NewAppClient()
err = client.SendRequest(ctx, "api", "add", http.MethodPost, req, &resp)
```

```go
func HandleRoute(path string, handler func(w http.ResponseWriter, r *http.Request))
```
The HandleRoute function will register the given function in the desired path. So, if you want to register a route in the path `add` (just like we want in our example), we pass "add" as the first parameter of the function and the function that adds the two numbers as the second parameter. Using it in our `api` app will look like this:

`api.go`
```go
func addHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Request
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			rest.ERROR(w, ierrors.New("wrong request body type").BadRequest())
			return
		}
		resp := model.Response{
			Result: data.Op1 + data.Op2,
		}
		fmt.Printf("Add Result: %v\n", resp.Result)
		rest.JSON(w, http.StatusOK, resp)
	}
}

// ...

client := insprclient.NewAppClient()
client.HandleRoute("add", addHandler())
```

Doing this with every route we need will make our dApp work as wanted. You can check the full example [here].

## Conclusion

With this simple steps you can add routes handlers to your dapps on Inspr. You just need to declare the endpoints in the yaml file, and use the sendRequest and handleRoute functions (provided by the Inspr Client) to send http requests to a desired node and to register functions that will run as soon as the path is reached.