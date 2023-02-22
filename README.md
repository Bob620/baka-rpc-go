# Baka-RPC-Go
Baka-RPC implements the [JSON-RPC 2.0](https://www.jsonrpc.org/specification) specifications with a small amount of
custom Baka-RPC to Baka_RPC functionality for improving compression and data transfer when possible. A lot more can be
done.

The registration of new methods can be a bit clunky visually, but is fairly easy to use.

### Client Creation and Connection
Creating a new BakaRPC communication channel is easy. 
```go
// Default method can be passed existing chans, or nil
rpcClient := rpc.CreateBakaRpc(nil, nil)

// Helper chan creators exist for Gorilla Websockets and regular Reader/Writer streams
streamClient := rpc.CreateBakaRpc(rpc.MakeReaderChan(streamIn), rpc.MakeWriterChan(streamOut))

// Sometimes you want to make the RPC client before the channels have been finalized, or change the channels later.
//   Once the client is made you can tell it to use the new channels.
rpcClient.UseChannels(rpc.MakeReaderChan(streamIn), rpc.MakeWriterChan(streamOut))

rpcClient.HandleDisconnect(func(uuid *Interface{}) {
	// One of the channels disconnected.
})
```

* Note: AddChannels allows for use of multiple channels and will send on all of them when making a call. Experimental.

### Method Registration and Calling
It is generally *Important* to register basic methods before establishing connections. While the client will hold off
until the other side confirms connection, if the other side has requests queued for delivery those will be sent asap.

This means that the requests could be processed before your code has the chance to register the corresponding methods.

Registering methods is handled strictly in order to guarantee compliance of input given to the method.

#### Parameters
`parameters.Param` is the interface for all parameters implemented.

`parameters` includes several common basic variable types already implemented, `GenericParam` can be used to send any
`json.RawMessage` data.

Parameters are made up of a Name, Default Value, and Required status.


`Name` is the name of the parameter in the method.

`Default` specifies the value.

`Required` tells Registration that it is required for the method to be called.

```go
parameters.GenericParam{
	Name:     "a param",
	Default:  json.RawMessage("{\"test\": \"object after Marshal\"}"),
	Required: false,
}
```

#### Registering
This next example creates a method named "Method Name", requiring a single string parameter named "a param" that echos
back to the caller.

* Note: This can be improved using Reflection to remove the manual casting step.

```go
rpcClient.RegisterMethod(
	"Method Name",
	[]parameters.Param{
		&parameters.StringParam{Name: "a param", Required: true},
	}, func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error) {
		// Unfortunately we know the type but still have to cast it.
		stringToEcho, _ := params["a param"].(*parameters.StringParam).GetString()
		
		log.Print("Received: ", stringToEcho)
		return json.Marshal(stringToEcho)
	})

rpcClient.DeregisterMethod("Method Name")
```

#### Calling Methods
Methods, as specified in the JSON-RPC spec, must have an ordered and by-name system for calling.

* Note: UUID is only useful with multiple channels, `nil` send on all channels.


`rpcClient.CallMethodByName(UUID, MethodName, ...&parameters.GenercParam{})`

`rpcClient.CallMethodByPosition(UUID, MethodName, ...&parameters.GenericParam{})`

`rpcClient.CallMethodWithNone(UUID, MethodName)`


These call the requested method if possible and params are valid, waiting for the remote client to respond.
The return value will be a `json.RawMessage` for client-side handling as the Method makes no guarantee of the response.

* Note: Using Reflection might be able to fix this but only for Baka-RPC communication
* Note: Best practice is to include `Name`, even in Positional arguments, but it's not required.

```go
rawData, resErr := rpcClient.Client.CallMethodByPosition(nil, "GetItem", &parameters.StringParam{Name: "itemID", Default: "123"})
if resErr != nil {
    return nil, errors.New(resErr.Message)
}
```

#### Notifying Methods
Notifying Methods are a method of Asynchronously calling a method and not caring about any return value. They work the
same as Call equivalents; However, they do not wait for the Method to finish nor provide a return value. 


`rpcClient.NotifyMethodByName(UUID, MethodName, ...&parameters.GenercParam{})`

`rpcClient.NotifyMethodByPosition(UUID, MethodName, ...&parameters.GenericParam{})`

`rpcClient.NoptifyMethodWithNone(UUID, MethodName)`