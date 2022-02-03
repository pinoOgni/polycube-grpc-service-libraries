# Polycube-grpc-cpp

Polycube-grpc-cpp is the library used to allow you to write the controlplane of a Polycube service in C ++, not subject to the conventional and current rules for writing a Polycube service. The library is divided into 3 parts:
* The .proto file where there is the definition of messages and rpc calls used
* Two pairs of self-generated .h / .cc files starting from the .proto file, in particular they are divided into ".grpc.pb" which implements the logic of gRPC and ".pb" which implements the logic of ProtoBuf
* A pair of .h / .cc files representing the layer that stands between the self-generated code and the actual Controlplane. Thus, the developer will only have to call the methods defined in this part of the library, ignoring all the gRPC part that is carried out within it.
Pro: 
* il servizio può essere remoto
* il servizio  non deve sottostare alle regole che ci sono ora
* lo sviluppatore è più libero


Contro:
* lo sviluppatore ha alcune cose in più da conoscere e da attenzionare
* lo sviluppatore è più libero


Questa libreria è solo una delle 3 che verrano create, cioè polycube-grpc-cpp, polycube-grpc-go e polycube-grpc-python. In pratica nelle ultime due, aumentano i PRO perché sarà possibile scrivere il controlplane in GO o in Python. 

Le librerie si metteranno tra il controlplane scritto in un linguaggio e Polycube, una sorta di API. Quindi sarà tipo:
* Polycube -- grpc server -- LIBRERIA -- grpc client -- Servizio (controlplane)