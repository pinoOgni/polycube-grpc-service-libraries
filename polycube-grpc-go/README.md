# Polycube gRPC GO library (WORK IN PROGRESS)

Questa è la libreria polycube-grpc-go che permette di scrivere una servizio Polycube in GO.


## Use it/Work on it

Creare una cartella con il nome del progetto/servizio, per esempio `pcn-helloworld2`
```
mkdir pcn-helloworld2
cd pcn-helloworld2
go mod init pcn-helloworld2
```
Creare il file dove scrivere il controlplane per esempio `helloworld2.go` e dentro importare la libreria 
```
pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
```
In seguito eseguire i seguenti comandi (TODO pinoOGni)

FORSE

```
git config --global url.git@github.com:.insteadOf https://github.com/
```

```
export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
```

Infine
```
go build helloworld2.go
```

Se si vuole apportare qualche modifica al file proto, per aggiungere una nuova feature, bisogna rigenerare i file. Per qualsiasi dubbio qui c'è la documentazione [ufficiale](https://developers.google.com/protocol-buffers/docs/reference/go-generated
)
1. scaricare protobuf e grpc-go
```
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
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



