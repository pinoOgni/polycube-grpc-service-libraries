#include <arpa/inet.h>
#include <grpcpp/grpcpp.h>

#include <cstdio>
#include <ctime>
#include <iostream>
#include <sstream>
#include <memory>
#include <string>
#include <unistd.h>
#include "polycube_grpc_commons.grpc.pb.h"
#include <thread>
#include "polycube/services/json.hpp"
#include "../Routing/Router.h"



using commons::Polycube;
using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using grpc::ClientReader;
using grpc::ClientReaderWriter;
using grpc::ClientWriter;
using commons::Polycube;
using commons::ToPolycubed;
using commons::ToServiced;
using commons::ServicedInfo;
using commons::Bool;
using commons::Index;
using commons::Empty;
using commons::Dataplane;
using commons::RestUserToServicedRequest;
using commons::DataplaneToServicedRequest;
using commons::ServicedToDataplaneReply;
using commons::ServicedToRestUserReply;
using commons::ServicedToDataplaneRequest;
using commons::CubeInfo;
using commons::MapValue;
using commons::SetRequest;
using commons::GetRequest;
using commons::RemoveRequest;
using commons::CubeManagement;
using Routing::Router;
using Routing::PathMatch;

namespace polycubeclient {

class PolycubeClient {
 public:

/*
  ====================================================================================================
                                        SERVICE MANAGEMENT
  ====================================================================================================
*/

  PolycubeClient();
  // PolycubeClient(std::shared_ptr<Channel> channel);

  // why the destructor does not work? Problem with compilation (pinoOgni)
  // ~PolycubeClient();

  void Subscribe(const std::string service_name, const int32_t service_uuid);
  Bool Unsubscribe();
  
/*
  ====================================================================================================
                                        STREAM MANAGEMENT
  ====================================================================================================
*/

  void ReadTheStream();
  // void Dispatcher(commons::ToServiced request, commons::ToPolycubed *reply);

/*
  ====================================================================================================
                                        CUBE MANAGEMENT
  ====================================================================================================
*/

  // maybe add also cubetype (pinoOgni)
  Bool CreateCube(const std::string cube_name, const std::string ingress_code, const std::string egress_code, const std::string conf);
  Bool DestroyCube(const std::string cube_name);
  // (pinoOgni)
  // void Reload(const std::string cube_name, const int32_t index, const commons::Dataplane_ProgramType program_type, const std::string code);

/*
  ====================================================================================================
                                        PORT MANAGEMENT
  ====================================================================================================
*/ 

  Bool SetPort(const std::string cube_name, const std::string port_name);
  Bool DelPort(const std::string cube_name, const std::string port_name);
  Bool SetPeer(const std::string cube_name, const std::string port_name, const std::string peer_name);

/*
  ====================================================================================================
                                        EBPF MAP MANAGEMENT
  ====================================================================================================
*/
  MapValue TableGetAll(const std::string cube_name, const std::string map_name, const commons::CubeInfo_ProgramType program_type, size_t key_size, size_t value_size);
  Bool TableSet(const std::string cube_name, const std::string map_name, const commons::CubeInfo_ProgramType program_type, void *key, size_t key_size, void *value, size_t value_size);
  MapValue TableGet(const std::string cube_name, const std::string map_name, const commons::CubeInfo_ProgramType program_type, void *key, size_t key_size, size_t value_size);
  Bool TableRemove(const std::string cube_name, const std::string map_name, const commons::CubeInfo_ProgramType program_type, void *key, size_t key_size);

/*
  ====================================================================================================
                                        ROUTES MANAGEMENT
  ====================================================================================================
*/

  void RegisterHandler(const std::string path,const std::string http_verb, ToPolycubed(*method)(ToServiced,PolycubeClient*));
  void RegisterHandlerGet(const std::string path, ToPolycubed(*method)(ToServiced,PolycubeClient*));
  void RegisterHandlerPost(const std::string path, ToPolycubed(*method)(ToServiced,PolycubeClient*));
  void RegisterHandlerPatch(const std::string path, ToPolycubed(*method)(ToServiced,PolycubeClient*));
  void RegisterHandlerDel(const std::string path, ToPolycubed(*method)(ToServiced,PolycubeClient*));
  void RegisterHandlerOptions(const std::string path, ToPolycubed(*method)(ToServiced,PolycubeClient*));

/*
  ====================================================================================================
                                        TODO OR MAYBE TODO
  ====================================================================================================
*/

  // (pinoOgni)
  void create_cube_by_id_handler(const char*name, const char*value);

 private:
  std::unique_ptr<Polycube::Stub> stub_;
  std::string target_ip = "127.0.0.1"; /* first version: grpc used on the same host, so the address is localhost */
  std::string target_port = "9001"; /* the grpc server on Polycubed is listening on thi port */
  Router router; /* router instance of: Simple C++ URI routing library */
  std::string base_url = "/polycube/v1";
  std::string service_name;  /* this is the service name of Serviced, gien by the Controlplane written by the developer */
  int32_t service_uuid; /* this is the uuid of Serviced, gien by the Controlplane written by the developer */
  std::string service_base_url; /* base_url + "/" + service_name */

public:
  /* 
    It represents a context that lasts for the life of 
    the Serviced, without it, the stream is deleted
  */
  ClientContext theContext;
  /* 
    This represents the long-lived stream, that is the connection thread between 
    the Serviced and Polycubed that will be used for the exchange of messages 
  */
  std::shared_ptr<ClientReaderWriter<ToPolycubed, ToServiced>> theStream; 
  
  // (pinoOgni)
  enum StatusPort {  UP, DOWN };
  struct Port {
    std::string name;
    std::string uuid;
    Status status;
    std::string peer;
    std::vector<std::string> tcubes;

  };

  
  /*
    This map is used to register for each key given by the set of template (URL) and http verb, 
    the corresponding handler. The handler corresponds to a method implemented by the Controlplane 
    and therefore by the developer, following the logic of the specific service
  */
  std::map<std::pair<std::string,std::string>, ToPolycubed(*)(ToServiced,PolycubeClient*)> methods_map;  

};



}
