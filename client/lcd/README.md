# Cosmos-LCD(light-client daemon) REST-server with swagger-ui

## Getting Start
* Execute the following command:
```
export GIN_MODE=release && gaiacli advanced rest-server-swagger --chain-id test-chain-RtAS0K
```
* Open this uri with your explorer:
```
http://localhost:1317/swagger/index.html
```
## More Options

* You can specify more full node URIs and localhost listen address with this command:
```
gaiacli advanced rest-server-swagger --chain-id test-chain-RtAS0K --laddr localhost:8080 --node-list tcp://10.10.10.10:46657,tcp://20.20.20.20:46657
```
* Then the swagger-ui URI will be:
```
http://localhost:8080/swagger/index.html
```