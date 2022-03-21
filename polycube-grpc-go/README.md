# Polycube gRPC GO library (WORK IN PROGRESS)

This is the polycube-grpc-go library that allows you to write a Polycube service in GO.


## Use it/Work on it

Create a folder with the name of the project/service, for example `pcn-helloworldgo`
```
mkdir pcn-helloworldgo
cd pcn-helloworldgo
```
Create a file `serviceName.go`, like `helloworldgo.go`, then:

```
go mod init pcn-helloworldgo
```
Create the file where to write the control plane for example `helloworldgo.go` and import the library inside
```
pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
```


NOTE: go.mod and go.sum could refer to an older version of the library than the current one, in which case the go build of the example would give various errors (in particular invalid version: unknown revision), therefore, delete go.mod and go.sum and execute the following commands:
Note of the Note: I am looking for a simpler solution

```

go mod init helloworldgo
go mod tidy
go build helloworldgo.go

```

NOTE: Only for the case this library is still part of a private repository

```
git config --global url.git@github.com:.insteadOf https://github.com/
```

```
export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
```

Finall  again:
```
go build helloworldgo.go
```

# Update the library

If you want to make any changes to the proto file, to add a new feature, you need to regenerate the files. For any doubt, here is the [official documentation](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
1. scaricare protobuf e grpc-go
```
go get -d -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -d -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
2. You may need to export the PATH
```
export PATH="$PATH:$PATH:$(go env GOPATH)/bin"
```
3. Inside the `polycube-grpc-go` directory
```
protoc --proto_path=../protos --go_out=commons --go_opt=paths=source_relative --go-grpc_out=commons --go-grpc_opt=paths=source_relative ../protos/polycube_grpc_commons.proto
```
Thanks to this command, the files will be regenerated and placed inside the commons folder, using the `.proto` file located in the` protos` directory


Finally, in order to use the local library, you need to modify the go.mod file of the service you are working on by adding this:

```
replace "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" v0.0.0-20220204154707-a45291043d50 => "local_path_to_polycube_grpc_go_library"
```


# OSSERVATIONS

* In this first version of the library, it is possible to create a cube with default configuration parameters. For example in the case of helloworld:

```
polycubectl helloworldgo h1 loglevel=trace span=true type=XDP_SKB shadow=true
```
But you can't set the action directly on cube creation, so you can't do something like that:
```
polycubectl helloworldgo h1 action=slowpath
```
This is to make the developer's life easier and not have to deal with configurations, json and more.



### Osservations regarding the gopath

1. watch out for gopath (from me it's netgroup/go)
2. follow here https://grpc.io/docs/languages/go/quickstart/ and do two go installs for protoc-gen-go and protoc-gen-grpc
3. go into polycube-grpc-go and issue the following command
```
protoc --proto_path=../protos --go_out=commons --go_opt=paths=source_relative --go-grpc_out=commons --go-grpc_opt=paths=source_relative ../protos/polycube_grpc_commons.proto
```

