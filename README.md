# Turbonomic-Go-Client

Use this simple GoLang Library to access the Turbonomic API. It currently supports authentication by using a username and password.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.23.7

## Authenticating to the Turbonomic API

This library creates a client object by passing in a `ClientParameters` struct with the following parameters:

- Hostname
- Username
- Password

Create a client with these parameters, similar to the following example:

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

You can then use this client to call other methods to interact with the Turbonomic API.

If your server has a self-signed certificate, you can skip SSL validation by also passing in the `Skipverify` parameter in the `ClientParameters` struct:

```
newClientOpts := ClientParameters{
    Hostname: "TurboHostname",
    Username: "TurboUsername",
    Password: "TurboPassword",
    Skipverify: true
    }
```

## Searching for an entity by name

To search for an entity by name, pass a `SearchRequest` struct to the `SearchEntityByName` method:

```
    searchReq := SearchRequest{
        Name:            "VM_Test_1",
        EntityType:      "VirtualMachine",
        EnvironmentType: "ONPREM",
        CaseSensitive:   true,
    }

    entityName, err := c.SearchEntityByName(searchReq)
```

## Retrieving an entity by UUID

To retrieve entity data based on its UUID, pass a `EntityRequest` struct to the `GetEntity` method:

```
    entityReq := EntityRequest{Uuid: "123456789"}

    entityName, err := c.GetEntity(entityReq)
```

## Retrieving actions based on request parameters and entity UUID

To retrieve action data, pass an `ActionsRequest` struct to the `GetActions` method:

```
    actionReq := ActionsRequest{
        Uuid: "123456789",
        ActionState: []string{"READY", "ACCEPTED", "QUEUED", "IN_PROGRESS"},
        ActionType: []string{"RESIZE"},
    }

    actions, err := GetActionsByUUID(actionReq)
```
