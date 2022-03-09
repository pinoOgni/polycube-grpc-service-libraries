# Polycube gRPC GO library (WORK IN PROGRESS)

Questa è la libreria polycube-grpc-go che permette di scrivere una servizio Polycube in GO.


## Use it/Work on it

Creare una cartella con il nome del progetto/servizio, per esempio `pcn-helloworldgo`
```
mkdir pcn-helloworldgo
cd pcn-helloworldgo
go mod init pcn-helloworldgo
```
Creare il file dove scrivere il controlplane per esempio `helloworldgo.go` e dentro importare la libreria 
```
pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
```


NOTA: il go.mod e il go.sum potrebbero fare riferimento ad una versione della libreria vecchia rispetto alla attuale, in quel caso il go build dell'esempio darebbe vari errori (in particolare invalid version: unknown revision), quindi, cancellare il go.mod e il go.sum ed eseguire i segunti comandi:


```

go mod init helloworldgo
go mod tidy
go build helloworldgo.go

```

NOTA: solo per il caso in cui questa libreria faccia parte ancora di una repository privata

```
git config --global url.git@github.com:.insteadOf https://github.com/
```

```
export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
```

Infine di nuovo
```
go build helloworldgo.go
```

Se si vuole apportare qualche modifica al file proto, per aggiungere una nuova feature, bisogna rigenerare i file. Per qualsiasi dubbio qui c'è la documentazione [ufficiale](https://developers.google.com/protocol-buffers/docs/reference/go-generated
)
1. scaricare protobuf e grpc-go
```
go get -d -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -d -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
2. Potresti aver bisogno di esportare il PATH
```
export PATH="$PATH:$PATH:$(go env GOPATH)/bin"
```
3. Dentro la directory `polycube-grpc-go`
```
protoc --proto_path=../protos --go_out=commons --go_opt=paths=source_relative --go-grpc_out=commons --go-grpc_opt=paths=source_relative ../protos/polycube_grpc_commons.proto
```
Grazie a questo comando, verranno rigenerati i file e posizionati dentro la cartella commons, utilizzando il file `.proto` posizionato nella directory `protos`


Infine per poter usare la libreria locale, bisogna modificare il file go.mod del servizio a cui si sta lavorando aggiungendo questo

```
replace "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" v0.0.0-20220204154707-a45291043d50 => "local_path_to_polycube_grpc_go_library"
```





NUOVO 24 02 2022

1. occhio al gopath (da me è netgroup/go)
2. segui qui https://grpc.io/docs/languages/go/quickstart/ e fai due go install per protoc-grn-go e protoc-gen-grpc 
3. vai dentro polycube-grpc-go e dai il seguente comando
```
protoc --proto_path=../protos --go_out=commons --go_opt=paths=source_relative --go-grpc_out=commons --go-grpc_opt=paths=source_relative ../protos/polycube_grpc_commons.proto
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