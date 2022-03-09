# Skeleteon

Here the skeleton of a Polycube gRPC service written in GO is presented. Starting from this, the developer only needs to worry about:
1. write the Control Plane with the logic of the service you want to implement
2. know the library proto file, with related messages and methods (ie RPC calls)
3. respect some constraints to allow the library to work correctly


Constraints:
* manage the creation of N cubes: it is in fact up to the programmer to decide how to manage the cubes, for example a vector of names or a map whose key is the name of the cubes and as a value a struct with specific information for this service
* the methods of the Control Plane must all have the same signature in order to be registered at a specific url. The signature is as follows:
```
func FuncName(request pb.ToServiced) pb.ToPolycubed { ... }
```
* 

