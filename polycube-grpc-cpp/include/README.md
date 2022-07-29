# Polycube-grpc-cpp

Polycube-grpc-cpp is the library used to allow you to write the controlplane of a Polycube service in C++, not subject to the conventional and current rules for writing a Polycube service. The library is divided into 3 parts:
* The .proto file where there is the definition of messages and rpc calls used
* Two pairs of self-generated .h/.cc files starting from the .proto file, in particular they are divided into ".grpc.pb" which implements the logic of gRPC and ".pb" which implements the logic of ProtoBuf
* A pair of .h/.cc files representing the layer that stands between the self-generated code and the actual Controlplane. Thus, the developer will only have to call the methods defined in this part of the library, ignoring all the gRPC part that is carried out within it.
Pro:
* the service can be remote
* the service does not have to comply with the rules that exist now
* the developer is freer


Versus:
* the developer has a few more things to know and watch out for
* the developer is freer


This library is only one of 3 that will be created, namely polycube-grpc-cpp, polycube-grpc-go and polycube-grpc-python. In practice, in the last two, the PROs increase because it will be possible to write the controlplane in GO or in Python.

The libraries will be placed between the controlplane written in a language and Polycube, a kind of API. So it's going to be like:
* Polycube - grpc server - LIBRARY - grpc client - Service (controlplane)