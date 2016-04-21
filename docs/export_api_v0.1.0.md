Export API
==========

This document defined the API used for exporting data from Raft

## Enabling Exporting ##
By default the exporting function is off, to turn it on `pfsd` must be run with
`-enable-export` flag.

The default port for the exporter is `10100`, and it can be changed with
`-export-port <port>`

## API Specification ##
| Name | Description | Available Options |
| ---- | ----------- | --------------- |
| type | Defines the type of the message | `state`<br> `nodechange`<br>`event`
| data | The data changes depending on what type the message is | [See Below](#Data)

The first message send is of type `state`, it is send on connecting to the server

### Data ###
| Type | Content |
| ---- | ------- |
| `state` | `nodes` - Array of [nodes](#Node) |
| `nodechange` | `action` - Can be of type `add`, `delete`, `update`<br>`node` - Single [node](#Node) to which the change applies. |
| `event` | `source` - Node UUID which is the source of the event<br>`target` - Node UUID which is the target of the event<br>`details` - Specific details of the event, declared internally

### Node ###
| Field Name | Description |
| ---------- | ----------- |
| `uuid` | Unique identifier for the node |
| `commonName` | Display Name of the node |
| `addr` | Address of the node, in a `host:port` format
| `state` | State of the node, can contain multiple states, with <code>&#124;</code> as the delimiter. It's the server's responsibility to make sure there are no race conditions in states. The available states are as follows: `current`, `leader`, `candidate`, `follower` (or blank), `inactive`

## Examples ##
```json
{
  "type":"state",
  "data": {
    "nodes": [
      {
        "uuid":"1234-abcd-5678-efgh",
        "commonName":"node-1",
        "addr":"10.0.0.1:67890",
        "state":"current|leader"
      }, {
        "uuid":"9012-ijkl-3456-mnop",
        "commonName":"node-2",
        "addr":"10.0.0.2:78901",
        "state":"inactive"
      }
    ]
  }
}
```

```json
{
  "type":"nodechange",
  "data": {
    "action":"update",
    "node":{
      "uuid":"9012-ijkl-3456-mnop",
      "commonName":"node-2",
      "addr":"10.0.0.2:78901",
      "state":"follower"
    }
  }
}
```

```json
{
  "type":"event",
  "data": {
    "source":"9012-ijkl-3456-mnop",
    "target":"1234-abcd-5678-efgh",
    "details": "write-request"
  }
}
```
