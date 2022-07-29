
#include "helloworldcpp_client.h"






int main(int argc, char **argv) {


  // Instantiate the client. It requires a channel, out of which the actual RPCs
  // are created. This channel models a connection to an endpoint specified by
  // the argument "--target=" which is the only expected argument.
  // We indicate that the channel isn't authenticated (use of
  // InsecureChannelCredentials()).


  PolycubeClient polycubeclientgrpc;

  
  srand (time(NULL));
  auto uuid = rand() % 10000 + 1;

  // salvare helloworldcpp in una var dentro polycubeclient
  
  polycubeclientgrpc.Subscribe("helloworldcpp",uuid);
  std::cout << "After Subscribe " << std::endl;

  // register paths
  polycubeclientgrpc.RegisterHandlerDelete(":cubeName/",&destroyCube);
  polycubeclientgrpc.RegisterHandlerPost(":cubeName/",&createCube);

  // pino: TODO a little bit strange this first solution
  polycubeclientgrpc.RegisterHandlerOptions("helloworldcpp",&helpCommands); 

  polycubeclientgrpc.RegisterHandlerGet(":cubeName/action/",&getAction);
  polycubeclientgrpc.RegisterHandlerPatch(":cubeName/action/",&setAction);
 
  polycubeclientgrpc.RegisterHandlerPost(":cubeName/ports/:portName/",&addPorts);
  polycubeclientgrpc.RegisterHandlerDelete(":cubeName/ports/:portName/",&delPort);
  polycubeclientgrpc.RegisterHandlerPatch(":cubeName/ports/:portName/peer/",&setPeer);
  polycubeclientgrpc.RegisterHandlerGet(":cubeName/ports/",&showPorts);
  // polycubeclientgrpc.RegisterHandlerGet(":cubeName/",&showCube);
  


  std::cout << "In Subscribe before thread" << std::endl;
  std::thread readthestream (&PolycubeClient::ReadTheStream,&polycubeclientgrpc);
  std::cout << "In Subscribe between thread and join" << std::endl;
  
  std::this_thread::sleep_for(std::chrono::milliseconds(5000));


  readthestream.join();
  return 0;
}

/*
	This method is used to get the action saved in the action_map table
*/
ToPolycubed getAction(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  
  // pino: TODO, grpc
  // uint8_t value = get_array_table<uint8_t>("action_map").get(0x0);
  
  commons::CubeInfo_ProgramType program_type = commons::CubeInfo_ProgramType_INGRESS;
  uint32_t key = 0;
  void *key_pass = &key;

  auto tableGetReply =  polycubeclientgrpc->TableGet(request.serviced_info().cube_name(), "action_map", program_type,  key_pass, sizeof(key), sizeof(uint8_t));
  
  int i = 0; 
  for (auto val : tableGetReply.bytevalues()) 
      printf("\n %d - \\x%.2x \n", i++, val);

  //std::cout << "GETACTION ACTION " << static_cast<std::underlying_type<HelloworldcppActionEnum>::type>(tableGetReply.bytevalues().operator[0]) << std::endl;
  std::cout << "GETACTION ACTION " << tableGetReply.bytevalues()[0] << std::endl;


  ToPolycubed ret_reply;
  ret_reply.mutable_serviced_to_rest_user_reply()->set_message("{\"action\":\"drop\"}");
  return ret_reply;
}


/*
  This method is used to set the action in the action_map table
*/
ToPolycubed setAction(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  
  // auto request_body = nlohmann::json::parse(std::string { request.rest_user_to_service_request().request_body() });
  // HelloworldcppActionEnum value = HelloworldcppJsonObject::string_to_HelloworldcppActionEnum(request_body);
  // uint8_t action = static_cast<uint8_t>(HelloworldcppActionEnum::DROP);
  uint8_t action = 1;


  commons::CubeInfo_ProgramType program_type = commons::CubeInfo_ProgramType_INGRESS;
  uint32_t key = 0;
  void *key_pass = &key;

  void *value_pass = &action;
  auto tableSetReply = polycubeclientgrpc->TableSet(request.serviced_info().cube_name(), "action_map", program_type,  key_pass, sizeof(key), value_pass, sizeof(action));
  ToPolycubed ret_reply;
  if(tableSetReply.status())
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
  else 
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
  return ret_reply;
}

ToPolycubed update_ports_map(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  // pino: TODO for now I'll put it INGRESS but then it will have to be taken somewhere
  commons::CubeInfo_ProgramType program_type = commons::CubeInfo_ProgramType_INGRESS;
  uint32_t key = 0;
  void *key_pass;
  uint16_t value;
  void *value_pass;

  Bool tableSetReply;
  bool ok = true;

  for (key=0;key<ports.size();key++ && ok == true) {
    key_pass = &key;
    value = key + 100; // like port->index() 
    value_pass = &value;
    tableSetReply = polycubeclientgrpc->TableSet(request.serviced_info().cube_name(), "ports_map", program_type,  key_pass, sizeof(key), value_pass, sizeof(value));
    if(tableSetReply.status() == false) {
      ok = false;
    }
  }


  value = UINT16_MAX;
  while (key < 2 && ok == true) {
    tableSetReply = polycubeclientgrpc->TableSet(request.serviced_info().cube_name(), "ports_map", program_type,  key_pass, sizeof(key), value_pass, sizeof(value));
    key++;
    if(tableSetReply.status() == false) {
      ok = false;
    }
  }


  ToPolycubed ret_reply;
  if(ok)
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
  else 
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
  return ret_reply;
}

/*
  First you make a request at the Polycubed level for adding the Port and then 
  move on to the part of the dataplane that depends on the service
*/
ToPolycubed addPorts(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  std::cout << "ToPolycubed addPorts(ToServiced request, PolycubeClient* polycubeclientgrpc) " << std::endl;
  if(ports.size() == 2)
        throw std::runtime_error("maximum number of ports reached");

  //it is like add_port<PortsJsonObject>(name, conf);
  ports.push_back(request.serviced_info().port_name());  //controlplane level 

  
  auto setPortReply = polycubeclientgrpc->SetPort(request.serviced_info().cube_name(), request.serviced_info().port_name());
 
  std::cout << "ToPolycubed addPorts(, setPortreply " << setPortReply.status() << std::endl;
  if(setPortReply.status()) 
    return update_ports_map(request,polycubeclientgrpc);
  else 
      throw std::runtime_error("Problem while adding port");
}

/*
  First you make a request at the Polycubed level for adding the Port and then 
  move on to the part of the dataplane that depends on the service
*/

// pino: TODO control if the throw is really ok by del a port that does not exists
ToPolycubed delPort(ToServiced request, PolycubeClient* polycubeclientgrpc) {
std::cout << "ToPolycubed delPort(ToServiced request, PolycubeClient* polycubeclientgrpc) " << std::endl;
if(ports.size() == 0)
        throw std::runtime_error("There are no port to be deleted");
std::cout << "aaaaaaaaa " << std::endl;
auto port_name = request.serviced_info().port_name();
std::cout << "bbbbbbbbbbbbbb " << std::endl;
auto pos = std::find_if(ports.begin(), ports.end(),
                          [port_name](std::string it) -> bool {
                            return it == port_name;
                          });
std::cout << "ccccccccccccccccccccc " << std::endl;
if (pos == std::end(ports)) {
    throw std::runtime_error("The port does not exists!");
}
std::cout << "dddddddddddddd " << std::endl;
  //it is like del_port<PortsJsonObject>(name, conf); or similar (I do not rember the name)
ports.erase(pos);  //controlplane level 
std::cout << "eeeeeeeeeeeee " << std::endl;
  
auto delPortReply = polycubeclientgrpc->DelPort(request.serviced_info().cube_name(), request.serviced_info().port_name());
std::cout << "ffffffffffffffff " << std::endl;
std::cout << "ToPolycubed addPorts(, delPortReply " << delPortReply.status() << std::endl;

if(delPortReply.status()) 
  return update_ports_map(request,polycubeclientgrpc);
else 
  throw std::runtime_error("Problem while deleting port");
}


ToPolycubed setPeer(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  ToPolycubed ret_reply;
  std::cout << "ToPolycubed setPeer(ToServiced request, PolycubeClient* polycubeclientgrpc) " << std::endl;
    if(ports.size() == 0)
         throw std::runtime_error("There are no ports to connect the peer to");
    auto port_name = request.serviced_info().port_name();
    auto pos = std::find_if(ports.begin(), ports.end(),
                            [port_name](std::string it) -> bool {
                              return it == port_name;
                            });

    if (pos == std::end(ports)) {
      throw std::runtime_error("The port does not exists!");
    }

    auto peer_name = request.serviced_info().peer_name();
    // auto peer_name = request.rest_user_to_service_request().request_body();
    auto setPeerReply = polycubeclientgrpc->SetPeer(request.serviced_info().cube_name(), port_name, peer_name);
    std::cout << "ToPolycubed setPeer(..) setPeerReply " << setPeerReply.status() << std::endl;
    if(setPeerReply.status())
      ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
    else 
      ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
    return ret_reply;
}


ToPolycubed showPorts(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  std::cout << "ToPolycubed showPorts(ToServiced request, PolycubeClient* polycubeclientgrpc) " << std::endl;
  ToPolycubed ret_reply;

  std::string ports_name = "PORTS: ";
  int i;
  for(i=0;i<ports.size();i++)
    ports_name += ports[i] + " ";
    if(i+1 != ports.size())
      ports_name += "; " ;
  std::cout << "showPorts " << ports_name << std::endl;

  ret_reply.mutable_serviced_to_rest_user_reply()->set_message(ports_name);

  return ret_reply;
}



ToPolycubed destroyCube(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  ToPolycubed ret_reply;
  auto destroy_cube_reply = polycubeclientgrpc->DestroyCube(request.serviced_info().cube_name());
  if(destroy_cube_reply.status())
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
  else 
    ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
  return ret_reply;
}


ToPolycubed createCube(ToServiced request, PolycubeClient* polycubeclientgrpc) {
    nlohmann::json jconf = json::parse(request.rest_user_to_service_request().request_body());
    // pino: TODO se conf.at non trova la chiave rimane in pending, come risolvere? se un non mette tutti i campi
    // del tipo shadow=true non arriva la chiave nella request body

    std::string conf = "{\"name\":\"" + request.serviced_info().cube_name() + "\",\"type\":\"TC\",\"loglevel\":\"INFO\",\"shadow\":false,\"span\":false,\"action\":\"drop\"}";
    
    // std::cout << "conf " << conf;
    // pino: TODO, extracto from request the name of cube and substitute with h1 and also with other fiels(shadow,....)
    // std::string conf = R"CONF({"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"})CONF";
    ToPolycubed ret_reply;
    auto create_cube_reply = polycubeclientgrpc->CreateCube(request.serviced_info().cube_name(),helloworldcpp_code_ingress,helloworldcpp_code_egress,conf);
    if(create_cube_reply.status() == true) {
        std::cout << "CREATED " << std::endl;
        ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ok");
        } else {
        std::cout << "NOT CREATED" << std::endl;
        ret_reply.mutable_serviced_to_rest_user_reply()->set_message("ko");
    }

    return ret_reply;
}




ToPolycubed helpCommands(ToServiced request, PolycubeClient* polycubeclientgrpc) {
  std::string help = R"HELP({
  "commands": [
    "add",
    "del",
    "show"
  ],
  "params": {
    "name": {
      "name": "name",
      "type": "key",
      "simpletype": "string",
      "description": "Name of the helloworld service",
      "example": "helloworld1"
    }
  },
  "elements": []
})HELP";
  std::cout << help << std::endl;
  ToPolycubed ret_reply;
  ret_reply.mutable_serviced_to_rest_user_reply()->set_message(help);
  return ret_reply;
}