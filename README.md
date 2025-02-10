# Turbonomic-Go-Client

This simple GoLang Library allows users to access Turbonomic's API. It currently supports authentication using
a username and password.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.22.2

## Authenticating

The library creates a client object by passing in the _ClientParameters_ struct with the following parameters:

- Hostname
- Username
- Password

You then create a new client as seen below:

```
newClientOpts := ClientParameters{
    Hostname: "TurboHostname",
    Username: "TurboUsername",
    Password: "TurboPassword",
    }

turboClient, err := NewClient(&newClientOpts)
if err != nil {
    panic(err)
}
```

This client can then be used to call other methods to interact with the Turbonomic API.

Note: If you have a server with a self signed certificate, you can also skip ssl validation by passing in the _Skipverify_ parameter:

```
newClientOpts := ClientParameters{
    Hostname: "TurboHostname",
    Username: "TurboUsername",
    Password: "TurboPassword",
    Skipverify: true
    }
```

## Searching for a Entity by Name

To search for an entity by name you use need to pass a _SearchRequest_ struct to the _SearchEntityByName_ method:

```
    searchReq := SearchRequest{
        Name:            "VM_Test_1",
        EntityType:      "VirtualMachine",
        EnvironmentType: "ONPREM",
        CaseSensitive:   true,
    }

    entityName, err := c.SearchEntityByName(searchReq)
```

## Retrieving an Entity By UUID

To retrieve entity data based on UUID, you pass a _EntityRequest_ struct to the _GetEntity_ method:

```
    entityReq := EntityRequest{Uuid: "123456789"}

    entityName, err := c.GetEntity(entityReq)
```

## Retrieve Actions based on request parameters and Entity UUID

To retrieve action data, you pass a _ActionsRequest_ struct to the _GetActions_ method:

```
    actionReq := ActionsRequest{
        Uuid: "123456789",
        ActionState: []string{"READY", "ACCEPTED", "QUEUED", "IN_PROGRESS"},
        ActionType: []string{"RESIZE"},
    }

    actions, err := GetActionsByUUID(actionReq)
```
