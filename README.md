# polycube-grpc-service-libraries
GRPC libraries to allow the writing of a stand-alone Polycube service also in other languages


This repository contains the client side grpc libraries to interface with Polycube and is composed as follows:

* protos: this directory contains the proto file on which all the libraries are based (both client and server side) and from which the respective files are auto-generated where the grpc calls and the various data structures are implemented. We have chosen to have only one for all N libraries (for now only the one in C++ and Go), so as to have a single modification point. There is another identical proto file in the Polycube repository that serves for the server side
* examples: contains some examples of writing a service (i.e. the controlplane), of how to use the grpc client libraries provided in this repository and also provides a guideline on some constraints to be respected (for example, a call to the Subscribe method must be made or record the url of the various methods)
* polycube-grpc-cpp: contains the grpc client library to be able to write a grpc service in C++. Unlike the classic Polycube service which is always written in C++, this one may not respect some constraints and can (in theory) be run on a different machine from the one where Polycube is located.
* polycube.-grpc-go: contains the grpc client library to be able to write a grpc service in Go
