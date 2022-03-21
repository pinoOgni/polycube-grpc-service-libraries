# Polycube gRPC Cpp library

## Note

**The polycube-grpc-cpp library is in a beta version, after an initial phase, efforts have focused on Go and Python We are trying to keep it updated as much as we can**

This is the Polycube gRPC library in C ++.


## TLDR

polycube-grpc-cpp is the library used to allow you to write the controlplane of a Polycube service in C++, not subject to the conventional and current rules for writing a Polycube service. The library is divided into 3 parts:
* The `.proto` file where there is the definition of messages and rpc calls used. It resides in the protos directory, since the `.proto` file is common to all libraries
* Two pairs of self-generated `.h/.cc` files starting from the .proto file, in particular they are divided into `.grpc.pb` which implements the logic of gRPC and `.pb` which implements the logic of ProtoBuf
* A pair of `.h/.cc` files representing the layer that stands between the self-generated code and the actual Controlplane. Thus, the developer will only have to call the methods defined in this part of the library, ignoring all the gRPC part that is carried out within it.


## Use it
To use it just run the following commands:

    ```
    cd polcyube-grpc-cpp
    mkdir build
    cd build
    cmake ..
    make -j4
    ```
If everything went well, the polycube-grpc-cpp library should be installed on the machine (in `/usr/include`)


## Create a Polycube gRPC service

ATTENTION: still a work in progress!

To write a Polycube gRPC service in C++, simply include this library
```
#include "polycube-grpc-cpp/polycube_client.h"
``` 

In the `main` you have to follow some simple rules:
* Create a `PolycubeClient polycubeclientgrpc;` object in this way you can use the library methods
* Sign up for the service to Polycube using the `Subscribe` method, where ServiceName is the service name and Uuid is a number generated as you wish
```
polycubeclientgrpc.Subscribe(ServiceName,Uuid);
```
* Register the handlers (i.e. the controlplane methods) for each URL, taking into account that the base URL is polycube/v1/serviceName (see examples)
```
polycubeclientgrpc.RegisterHandlerPost(":cubeName/",&createCube);
```
* Then, there are two possible choices
    * the developer decides that the Controlplane should only handle requests from the user then calls the ReadTheStream method normally
    * the developer decides that the Controlplane in addition to handling user calls, must also do something else so call the ReadTheStream method using a thread and then join

Turning now to the methods that implement the Controlplane, they are simple C ++ methods that must mirror the logic of the service, inside which a call is made to the gRPC method of the library. For example:
```
ToPolycubed setAction(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  
  ....
  // setting all variables

  auto tableSetReply = polycubeclientgrpc->TableSet(request.serviced_info().cube_name(), "action_map", program_type,  key_pass, sizeof(key), value_pass, sizeof(action));
  
  ToPolycubed ret_reply;
  if(tableSetReply.status())
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
  else 
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
  return ret_reply;
}

```