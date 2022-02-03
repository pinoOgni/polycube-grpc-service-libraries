# Polycube gRPC GO library (WORK IN PROGRESS)

Questa è la libreria polycube-grpc-go che permette di scrivere una servizio Polycube in GO.


## Use it/Work on it

Per usare la libreria basta usare il tool Go per buildare il controlplane del servizio e per importare la libreria da github



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



