# Client Controller

The client controller is the entry point for making changes to the Inspr tree structure. By using the client, it is possible to create, modify and delete dApps, Channels and Types.

## Instantiating a new Client

The structure of the `client` basically needs a `request.Client` that defines where the requests will be sent. So the first step is to instantiate a new `request.Client`. For this, it is possible to use the builder defined in the `request` package:

```go
rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
```
Where `url` is the route to which Inspr Daemon is listening.

Then, use the `NewControllerClient` function defined in the `client` package passing the created `request.Client` as a parameter to instantiate a new client controller:

```go
client = client.NewControllerClient(rc)
```

## Using the Client Controller

To use the client, call the respective function to the type of structure you want to manipulate followed by the operation function that must be done. For example, to create a new dApp `HelloWorldApp` inside the root dApp, just do:

```go
resp, err := client.Apps().Create(context.Background(), "", &meta.App{
    Meta: meta.Metadata{
        Name: "HelloWorldApp",
    },
    Spec: meta.AppSpec{},
}, dryRun)
```

In the example above, the function for creating a dApp receives:
*  A [go context](https://golang.org/pkg/context/) (context.Background())
*  The path in which that dApp will be created("")
*  The dApp itself (&meta.App{...})
*  The dryRun flag, which is present in all methods other than `get`, is a bool that indicates whether the modifications should really be applied to the structure or if they are simply used to visualize which changes would be made.

Similarly, to create a Channel called "NewChannel" within the `HelloWorldApp` dApp that was just created, do:

```go
resp, err := client.Channels().Create(context.Background(), "HelloWorldApp", &meta.Channel{
    Meta: meta.Metadata{
        Name: "NewChannel",
    },
    Spec: meta.ChannelSpec{
        Type: "TypeHello",
    },
}, dryRun)
```
Remember that in the case above, the Type `TypeHello` must exist within `HelloWorldApp`.

## Apps

### func \(\*AppClient) Get

```go
func (ac *AppClient) Get(ctx context.Context, context string) (*meta.App, error)
```
`Get` gets information from a dApp inside the Insprd. The `context string` refers to the dApp itself, represented with a dot separated query, such as **app1.app2**.  
So to get a dApp inside `app1` that is called `app2` you would call ac.Get(context.Background(), "app1.app2").

### func \(\*AppClient) Create

```go
func (ac *AppClient) Create(ctx context.Context, context string, app *meta.App, dryRun bool) (diff.Changelog, error)
```
`Create` creates a dApp inside the Insprd. The `context string` refers to the parent dApp where the actual dApp will be instantiated, represented with a dot separated query, such as **app1.app2**. The information of the dApp, such as name and other metadata, will be extracted from the definition of the dApp itself.   
So to create a dApp inside `app1` with the name app2 you would call `ac.Create(context.Background(), "app1", &meta.App{...}, false)`.

### func \(\*AppClient) Update

```go
func (ac *AppClient) Update(ctx context.Context, context string, app *meta.App, dryRun bool) (diff.Changelog, error)
```
`Update` updates a dApp inside the Insprd. If the dApp doesn't exist, it will return a error. The `context string` refers to the parent dApp where the actual dApp will be instantiated, represented with a dot separated query, such as **app1.app2**. The information of the dApp, such as name and other metadata, will be extracted from the definition of the dApp itself.   
So to update a dApp inside `app1` with the name `app2` you would call `ac.Update(context.Background(), "app1", &meta.App{...}, false)`.

### func \(\*AppClient) Delete

```go
func (ac *AppClient) Delete(ctx context.Context, context string, dryRun bool) (diff.Changelog, error)
```
`Delete` deletes a dApp inside the Insprd. The `context string` refers to the dApp itself, represented with a dot separated query, such as **app1.app2**.  
So to delete a dApp inside `app1` with the name `app2` you would call `ac.Delete(context.Background(), "app1.app2")`.

## Channels

### func \(\*ChannelClient) Get

```go
func (cc *ChannelClient) Get(ctx context.Context, context string, name string) (*meta.Channel, error)
```
`Get` gets a Channel from the Insprd. The `context string` refers to the parent dApp of the given Channel, represented with a dot separated query, such as app1.app2. The name is the name of the Channel.  
So to search for a Channel inside `app1` with the name `channel1` you would call `cc.Get(context.Background(), "app1", "channel1")`.

### func \(\*ChannelClient) Create

```go
func (cc *ChannelClient) Create(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
```
`Create` creates a Channel inside the Insprd. The `context string` refers to the parent dApp of the given Channel, represented with a dot separated query, such as **app1.app2**. The Channel information such as its name will be extracted from the given Channel's metadata.  
So to create a Channel inside `app1` with the name `channel1` you would call `cc.Create(context.Background(), "app1", &meta.Channel{...})`.

### func \(\*ChannelClient) Update

```go
func (cc *ChannelClient) Update(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
```
`Update` updates a Channel inside the Insprd. The `context string` refers to the parent dApp of the given Channel, represented with a dot separated query, such as **app1.app2**. The Channel information such as its name will be extracted from the given Channel's metadata.  
So to update a Channel inside `app1` with the name `channel1` you would call `cc.Update(context.Background(), "app1", &meta.Channel{...})`.

### func \(\*ChannelClient) Delete

```go
func (cc *ChannelClient) Delete(ctx context.Context, context string, name string, dryRun bool) (diff.Changelog, error)
```
`Delete` deletes a Channel inside the Insprd. The `context string` refers to the parent dApp of the given Channel, represented with a dot separated query, such as **app1.app2**. The name is the name of the Channel to be deleted.  
So to delete a Channel inside `app1` with the name `channel1` you would call `cc.Delete(context.Background(), "app1", "channel1")`.

## Types

### func \(\*TypeClient) Get

```go
func (ctc *TypeClient) Get(ctx context.Context, context string, name string) (*meta.Type, error)
```
`Get` gets a Type from the Insprd. The `context string` refers to the parent dApp of the given Type, represented with a dot separated query, such as **app1.app2**. The name is the name of the Type.  
So to search for a Type inside `app1` with the name `Type1` you would call ctc.Get(context.Background(), "app1", "Type1").

### func \(\*TypeClient) Create

```go
func (ctc *TypeClient) Create(ctx context.Context, context string, ch *meta.Type, dryRun bool) (diff.Changelog, error)
```
`Create` creates a Type inside the Insprd. The `context string` refers to the parent dApp of the given Type, represented with a dot separated query, such as **app1.app2**. The Type information such as its name will be extracted from the given Type's metadata.  
So to create a Type inside `app1` with the name `Type1` you would call `ctc.Create(context.Background(), "app1", &meta.Type{...})`.

### func \(\*TypeClient) Update

```go
func (ctc *TypeClient) Update(ctx context.Context, context string, ch *meta.Type, dryRun bool) (diff.Changelog, error)
```
`Update` updates a Type inside the Insprd. The `context string` refers to the parent dApp of the given Type, represented with a dot separated query, such as **app1.app2**. The Type information such as its name will be extracted from the given Type's metadata.  
So to update a Type inside `app1` with the name `Type1` you would call ` ctc.Create(context.Background(), "app1", &meta.Type{...})`.

### func \(\*TypeClient) Delete

```go
func (ctc *TypeClient) Delete(ctx context.Context, context string, name string, dryRun bool) (diff.Changelog, error)
```
`Delete` deletes a Type inside the Insprd. The `context string` refers to the parent dApp of the given Type, represented with a dot separated query, such as **app1.app2**. The name is the name of the Type to be deleted.   
So to delete a Type inside `app1` with the name `Type1` you would call `ctc.Delete(context.Background(), "app1", "Type1")`.