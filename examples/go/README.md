# How the examples work (pcn-helloworldgo case)

To run an example just run the following commands:
```
cd examples/go/pcn-helloworldgo
go run helloword.go 
```
There may be problems with go.sum and go.mod, in that case just follow the steps below:
```
cd examples/go/pcn-helloworldgo
rm go.sum
rm go.mod
go mod init helloworldgo
go mod tidy
```
To use the local version of the library (maybe it has been modified or you have problem with go module) just add a similar line to the go.mod file:
```
replace github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go v0.0.0-20220204154707-a45291043d50 => /home/netgroup/polycube-grpc-service-libraries/polycube-grpc-go
```
Where the version is the latest commit and the part after the "=>" represents where is the library



