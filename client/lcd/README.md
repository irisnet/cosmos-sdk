# Cosmos-LCD(light-client daemon) REST-server with swagger-ui

## Getting Start
* If you have gaiad running on your local machine, and its listened port is 26657, then you can start Cosmos-LCD just with the following command:
```
gaiacli advanced rest-server-swagger --chain-id {your chain id}
```
* Open this uri with your explorer:
```
http://localhost:1317/swagger/index.html

```
## More Options

| Parameter       | Type      | Default                 | Required | Description                                          |
| -----------     | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| home            | string    | "$HOME/.gaiacli"        | false    | directory for save checkpoints and validator sets    |
| chain-id        | string    | null                    | true     | chain id of the full node to connect                 |
| node-list       | URL       | "tcp://localhost:26657" | false    | addresses of the full node to connect                |
| laddr           | URL       | "tcp://localhost:1317"  | false    | address to run the rest server on                    |
| trust-node      | bool      | "false"                 | false    | Whether this LCD is connected to a trusted full node |
| swagger-host-ip | string    | "localhost"             | false    | The IP of the server which Cosmos-LCD is running on  |
| modules         | string    | "general,key,token"     | false    | enabled modules.                                     |

* You can specify more full node URIs with this command:
```
gaiacli advanced rest-server-swagger --chain-id {your chain id} --node-list tcp://10.10.10.10:26657,tcp://20.20.20.20:26657
```

* If you want to run Cosmos-LCD on remote server, you can start it like this:
```
gaiacli advanced rest-server-swagger --chain-id {your chain id} --laddr 0.0.0.0:1317 --swagger-host-ip {your server ip}
```