This is the first example of a grpc C++ service. Within this directory there are various types of files and directories.

Directory:
* cmake: contains the cmake useful for building the service
* datamodel: contains the yang of the service (FOR NOW NOT USED)
* src: Contains all the helloworld 1.0 service code and other files that make up the grpc service. Of all the files contained here, only the following are used:
    * helloworldcpp_client.cc represents the CONTROLPLANE. Here is the main of the service and a few things need to be done to get everything to work:
        * PolycubeClient polycubeclientgrpc; in this way we can use the methods of the polycube-grpc library (C++)
        * use the RegisterHandlerHTTPVERB methods to associate the controlplane methods to a url
        * launch the ReadTheStream method of the polycube-grpc library in a thread
        * join the newly launched thread
    * helloworldcpp_client.h represents the CONTROLPLANE.
        * Here there is not only the prototype of the method but also the implementation. Each method that is implemented here and that needs to contact Polycube for its work, will have to make a call to a suitable method of the polycube-grpc library.
        * For example, the Helloworld setAction method is used to set an action in the DATAPLANE that is to set a value in an eBPF map, this is reduced to calling the TableSet method with appropriate parameters.
        * Note that there is also the create_cube method which is used to create a cube
        * IMPORTANT: here, the developer will be able to decide the logic with which to manage the relationship between service and cube, so for example manage the cubes with IDs and have a vector of cubes or similar
        * IMPORTANT: here, the developer must also manage the logic of the ports, obviously if the service has the concept of port in itself
        * LESS IMPORTANT: right now in this file, the developer will also have to put the helpocommands method which will then be registered on the appropriate url and which returns a description on how to use the service from polycubectl
    * Helloworldcpp_dp_egress.h, Helloworldcpp_dp_egress.c Helloworldcpp_dp_ingress.h and Helloworldcpp_dp_ingress.c which correspond to the dataplane
    * Helloworldcpp-lib.cpp, Helloworldcpp.cpp, Helloworldcpp.h, Ports.cpp and Ports.h NOW as of now they are not needed, I don't know if they will be used in the future
* first_client.cc, first_polycube_client.cc, router_client.cc and second_client.cc DO NOT do the service but are test files


NOTE: To make everything work, you will need to include the polycube-grpc library in the controlplane files






# Examples gRPC remote services in C++

##






## Observations

Within each example, for logic reasons, or because it is still being implemented, there are some directories or files that are not necessary to develop a simple service. In particular, the directories
* cmake: contains some files for loading the files as variables and another one used for the compiler of gRPC/ProtoBuf
* datamodel: contains the datamodel of the Polycube helloworld service