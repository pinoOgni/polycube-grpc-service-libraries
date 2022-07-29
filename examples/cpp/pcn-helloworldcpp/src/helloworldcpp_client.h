
#include "polycube-grpc-cpp/polycube_client.h"
#include <chrono>
#include <thread>
#include <stdlib.h> 
#include "Helloworldcpp_dp_egress.h"
#include "Helloworldcpp_dp_ingress.h"
#include "polycube/services/json.hpp"

enum class HelloworldActionEnum {
  DROP, SLOWPATH, FORWARD
};

using nlohmann::json;

using namespace polycubeclient;

std::vector<std::string> ports;


ToPolycubed getAction(ToServiced request, PolycubeClient* polycubeclientgrpc);
ToPolycubed setAction(ToServiced request, PolycubeClient* polycubeclientgrpc);

ToPolycubed update_ports_map(ToServiced request, PolycubeClient* polycubeclientgrpc);

/*
  First you make a request at the Polycubed level for adding the Port and then 
  move on to the part of the dataplane that depends on the service
*/
ToPolycubed addPorts(ToServiced request, PolycubeClient* polycubeclientgrpc);


/*
  First you make a request at the Polycubed level for adding the Port and then 
  move on to the part of the dataplane that depends on the service
*/

// pino: TODO control if the throw is really ok by del a port that does not exists
ToPolycubed delPort(ToServiced request, PolycubeClient* polycubeclientgrpc); 
ToPolycubed setPeer(ToServiced request, PolycubeClient* polycubeclientgrpc);
ToPolycubed showPorts(ToServiced request, PolycubeClient* polycubeclientgrpc);

ToPolycubed destroyCube(ToServiced request, PolycubeClient* polycubeclientgrpc);
ToPolycubed createCube(ToServiced request, PolycubeClient* polycubeclientgrpc);

ToPolycubed helpCommands(ToServiced request, PolycubeClient* polycubeclientgrpc);

