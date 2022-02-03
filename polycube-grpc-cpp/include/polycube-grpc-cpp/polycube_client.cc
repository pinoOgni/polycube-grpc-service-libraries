#include "polycube_client.h"

namespace polycubeclient {

/*
  ====================================================================================================
                                        SERVICE MANAGEMENT
  * PolycubeClient();
  * void Subscribe(const std::string service_name, const int32_t service_uuid);
  * Bool Unsubscribe();
  ====================================================================================================
*/

// TODO pino: understand the use of CreateChannel, CreateCustomChannel,
// keepalive and channel in general
PolycubeClient::PolycubeClient() {
  std::string target_str = target_ip + ":" + target_port;

  grpc::ChannelArguments channelArguments;
  channelArguments.SetInt(GRPC_ARG_KEEPALIVE_TIME_MS, 5000);
  channelArguments.SetInt(GRPC_ARG_KEEPALIVE_TIMEOUT_MS, 5000);
  channelArguments.SetInt(GRPC_ARG_KEEPALIVE_PERMIT_WITHOUT_CALLS, 1);
  channelArguments.SetInt(GRPC_ARG_HTTP2_BDP_PROBE, 1);
  channelArguments.SetInt(GRPC_ARG_HTTP2_MIN_RECV_PING_INTERVAL_WITHOUT_DATA_MS,
                          5000);
  channelArguments.SetInt(GRPC_ARG_HTTP2_MIN_SENT_PING_INTERVAL_WITHOUT_DATA_MS,
                          5000);
  channelArguments.SetInt(GRPC_ARG_HTTP2_MAX_PINGS_WITHOUT_DATA, 0);

  std::shared_ptr<grpc::Channel> channel = grpc::CreateCustomChannel(
      target_str, grpc::InsecureChannelCredentials(), channelArguments);

  stub_ = Polycube::NewStub(channel);
}

/*
  This method is used to register a service with the name service_name to
  Polycube. Save the context and stream for future use. Basically it serves to
  set up a long-lived grpc
*/
void PolycubeClient::Subscribe(const std::string service_name,
                               const int32_t service_uuid) {
  this->service_name = service_name;
  this->service_uuid = service_uuid;
  this->service_base_url = base_url + "/" + service_name;

  std::cout << "Subscribe " << service_name << " " << service_uuid << std::endl;
  theStream = stub_->Subscribe(&theContext);
  std::cout << "Subscribe theStream->use_count()  " << theStream.use_count()
            << std::endl;
  std::cout << "INSIDE Subscribe() theStream " << theStream << std::endl;
  ToPolycubed request;
  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);
  request.set_subscribe(true);

  auto ret = theStream->Write(request);
  std::cout << "Subscribe returned write " << ret << std::endl;
  // stream->WritesDone(); // yes or no ????

  // pino TODO: maybe here we can read the stream launching a thread for a
  // method ReadTheStream then in the subscribe we can use the join for the
  // thread or std::terminate() or nothing (because the program will close after
  // the unsubscribe)

  // pino TODO: remove these lines, the thread should be moved to the
  // service/main and it will be a MUST for the developer. If the thread is
  // called here, the service/main will stop immediatily if there is no infinite
  // while loop

  /*
  std::cout << "In Subscribe before thread" << std::endl;
  std::thread readthestream (&PolycubeClient::ReadTheStream,this);
  std::cout << "In Subscribe between thread and join" << std::endl;


  // Separates the thread of execution from the thread object, allowing
  execution to continue independently.
  // Any allocated resources will be freed once the thread exits.

  readthestream.detach();
  std::cout << "In Subscribe before after join" << std::endl;
  */
}

/*
  This server method to unsubscribe the service from Polycube
*/
Bool PolycubeClient::Unsubscribe() {
  ServicedInfo request;
  request.set_service_name(service_name);
  request.set_uuid(service_uuid);

  Bool reply;

  ClientContext context;

  // The actual RPC.
  Status status = stub_->Unsubscribe(&context, request, &reply);
  // Act upon its status.
  if (status.ok()) {
    return reply;
  } else {
    std::cout << status.error_code() << ": " << status.error_message()
              << std::endl;
    return reply;
  }
}

/*
  ====================================================================================================
                                        STREAM MANAGEMENT
  * void ReadTheStream();
  ====================================================================================================
*/

/*
 This method is used to continuously read the stream. It is launched in a thread
*/
void PolycubeClient::ReadTheStream() {
  std::cout << "ReadTheStream  theStream " << this->theStream << std::endl;
  std::cout << "ReadTheStream theStream->use_count()  "
            << this->theStream.use_count() << std::endl;
  ToServiced request;
  ToPolycubed reply;
  // theStream->WritesDone();
  std::cout << "\n Waiting for something to read from the stream .... \n "
            << std::endl;
  while (this->theStream->Read(&request)) {
    std::cout << "true " << request.has_rest_user_to_service_request()
              << std::endl;
    std::cout << "false " << request.has_dataplane_to_serviced_request()
              << std::endl;
    std::cout << "method " << request.rest_user_to_service_request().http_verb()
              << std::endl;
    std::cout << "body "
              << request.rest_user_to_service_request().request_body()
              << std::endl;
    std::cout << "url " << request.rest_user_to_service_request().url()
              << std::endl;

    /*
      For now two cases: rest user to serviced request and dataplane to serviced
      request
    */
    try {
      if (request.has_rest_user_to_service_request()) {
        PathMatch match = this->router.matchPath(
            request.rest_user_to_service_request().url(),
            request.rest_user_to_service_request().http_verb());
        std::cout << "pathTemplate " << match.pathTemplate()
                  << std::endl;  // returns the generic url template
        std::cout
            << "path " << match.path()
            << std::endl;  // returns the requested url created by matchPath

        auto method_call = this->methods_map.at(
            std::make_pair(match.pathTemplate(), match.http_verb()));
        if (method_call == nullptr)
          std::cout << "method_call è nullptr";
        else
          std::cout << "tutto ok method_call";

        std::string cubeName = match["cubeName"];
        std::cout << "cube name from match " << cubeName << std::endl;
        request.mutable_serviced_info()->set_cube_name(cubeName);

        std::string portName = match["portName"];
        std::cout << "port name from match " << portName << std::endl;
        request.mutable_serviced_info()->set_port_name(portName);

        // pino: this is a trick, maybe there is a bettere way, the peer name is
        // "peerName" and not peerName, so we need to remove the two double
        // quotas
        std::string peer =
            (request.rest_user_to_service_request().request_body().empty() ==
             true)
                ? " "
                : request.rest_user_to_service_request().request_body();
        if (!peer.empty()) {
          peer = peer.substr(1, peer.size() - 2);
          // (pinoOgni) I dunno if there is a bettere way
          request.mutable_serviced_info()->set_peer_name(peer);
          std::cout << "peer name from match " << peer << std::endl;
        }
        // maybe i could launch a thread to fulfill the request (pinoOgni)

        /*
          here the correct method is called using the map where each key is
          associated with a handler that represents a method implemented by the
          final developer that is the Service Controlplane
        */
        reply = method_call(request, this);
        std::cout << "MESSAGE PINO "
                  << reply.serviced_to_rest_user_reply().message() << std::endl;
      } else if (request.has_dataplane_to_serviced_request()) {
      } else if (request.has_dataplane_to_serviced_packet()) {
        std::cout << "request.dataplane_to_serviced_packet().cube_id() "
                  << request.dataplane_to_serviced_packet().cube_id()
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().port_id() "
                  << request.dataplane_to_serviced_packet().port_id()
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().packet_len() "
                  << request.dataplane_to_serviced_packet().packet_len()
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().reason() "
                  << request.dataplane_to_serviced_packet().reason()
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().traffic_class() "
                  << request.dataplane_to_serviced_packet().traffic_class()
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().metadata[0] "
                  << request.dataplane_to_serviced_packet().metadata(0)
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().metadata[1] "
                  << request.dataplane_to_serviced_packet().metadata(1)
                  << std::endl;
        std::cout << "request.dataplane_to_serviced_packet().metadata[2] "
                  << request.dataplane_to_serviced_packet().metadata(2)
                  << std::endl;
        for (auto i = 0;
             i < request.dataplane_to_serviced_packet().packet_size(); i++)
          std::cout << "i = " << i << " : "
                    << request.dataplane_to_serviced_packet().packet(i)
                    << std::endl;

        /*
          (pinoOgni)
          here I should call a method as in the case of before using the methods_map map where 
          the real logic written by the controlplane developer will be implemented
        */
      } else {
        std::cout << "ERROR: message case not implemented!" << std::endl;
      }

      std::cout << "WRITING REPLY TO POLYCUBED" << std::endl;
      this->theStream->Write(reply);
      std::cout << "REPLY TO POLYCUBED WRITTEN" << std::endl;
      std::cout << "\n Waiting for something to read from the stream .... \n "
                << std::endl;

    } catch (const std::exception &exc) {
      std::cout << "EXCEPTION " << exc.what() << std::endl;
      reply.mutable_serviced_to_rest_user_reply()->set_success(false);
      reply.mutable_serviced_to_rest_user_reply()->set_message(exc.what());
      this->theStream->Write(reply);
    }
  }
}

/*
  ====================================================================================================
                                        CUBE MANAGEMENT
  * Bool CreateCube(const std::string cube_name, const std::string ingress_code,
  const std::string egress_code, const std::string conf);
  * Bool DestroyCube(const std::string cube_name);
  ====================================================================================================
*/

/*
  This method is used to ask Polycube to create a remote service cube (gRPC)
*/
Bool PolycubeClient::CreateCube(const std::string cube_name,
                                const std::string ingress_code,
                                const std::string egress_code,
                                const std::string conf) {
  CubeManagement request;
  request.set_service_name(service_name);
  request.set_uuid(service_uuid);
  request.set_cube_name(cube_name);
  request.set_ingress_code(ingress_code);
  request.set_egress_code(egress_code);
  request.set_conf(conf);

  Bool reply;

  ClientContext context;

  // The actual RPC.
  Status status = stub_->CreateCube(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK CreateCube " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO CreateCube " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  This method is used to request Polycube to destroy a remote service cube
  (gRPC)
*/
Bool PolycubeClient::DestroyCube(const std::string cube_name) {
  CubeManagement request;
  request.set_service_name(service_name);
  request.set_uuid(service_uuid);
  request.set_cube_name(cube_name);

  Bool reply;

  ClientContext context;

  // The actual RPC.
  Status status = stub_->DestroyCube(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK DESTROYED" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO DESTROYED" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  ====================================================================================================
                                        PORT MANAGEMENT
  * Bool SetPort(const std::string cube_name, const std::string port_name);
  * Bool DelPort(const std::string cube_name, const std::string port_name);
  * Bool SetPeer(const std::string cube_name, const std::string port_name, const
  std::string peer_name);
  ====================================================================================================
*/

/*
  This method is used to add a port to a cube. The concept of a door and
  therefore of adding a door is different from adding a value to a map of a
  cube. In fact you have to go through Polycubed which has a view of the doors
  that exist. This is why this method was created.
*/
Bool PolycubeClient::SetPort(const std::string cube_name,
                             const std::string port_name) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_port()->set_name(port_name);

  // Container for the data we expect from the server.
  Bool reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->SetPort(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK status.ok() " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO status.ok()" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  This method is used to remove a port to a cube. The concept of a door and
  therefore of adding a door is different from adding a value to a map of a
  cube. In fact you have to go through Polycubed which has a view of the doors
  that exist. This is why this method was created.
*/
Bool PolycubeClient::DelPort(const std::string cube_name,
                             const std::string port_name) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_port()->set_name(port_name);

  // Container for the data we expect from the server.
  Bool reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->DelPort(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK status.ok() " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO status.ok()" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  This method is used to set a peer for a given port
*/

Bool PolycubeClient::SetPeer(const std::string cube_name,
                             const std::string port_name,
                             const std::string peer_name) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_port()->set_name(port_name);
  request.mutable_port()->set_peer(peer_name);
  // Container for the data we expect from the server.
  Bool reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->SetPeer(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK status.ok() " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO status.ok()" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  ====================================================================================================
                                        EBPF MAP MANAGEMENT
  * MapValue TableGetAll(const std::string cube_name, const std::string
  map_name, const commons::CubeInfo_ProgramType program_type, size_t key_size,
  size_t value_size);
  * Bool TableSet(const std::string cube_name, const std::string map_name, const
  commons::CubeInfo_ProgramType program_type, void *key, size_t key_size, void
  *value, size_t value_size);
  * MapValue TableGet(const std::string cube_name, const std::string map_name,
  const commons::CubeInfo_ProgramType program_type, void *key, size_t key_size,
  size_t value_size);
  * Bool TableRemove(const std::string cube_name, const std::string map_name,
  const commons::CubeInfo_ProgramType program_type, void *key, size_t key_size);
  ====================================================================================================
*/

/*
  This method is used to get an entire map from the service dataplane
*/
MapValue PolycubeClient::TableGetAll(
    const std::string cube_name, const std::string map_name,
    const commons::CubeInfo_ProgramType program_type, size_t key_size,
    size_t value_size) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_cube_info()->set_map_name(map_name);
  request.mutable_cube_info()->set_program_type(program_type);
  request.mutable_cube_info()->set_index(0);

  request.mutable_get_request()->set_all(true);

  request.mutable_get_request()->set_key_size(key_size);
  request.mutable_get_request()->set_value_size(value_size);
  // Container for the data we expect from the server.
  MapValue reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->TableGetAll(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    return reply;
  } else {
    std::cout << status.error_code() << ": " << status.error_message()
              << std::endl;
    return reply;
  }
}

/*
  Given a key, this method is used to ask Polycube for data from an eBPF map of
  a cube
*/
MapValue PolycubeClient::TableGet(
    const std::string cube_name, const std::string map_name,
    const commons::CubeInfo_ProgramType program_type, void *key,
    size_t key_size, size_t value_size) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_cube_info()->set_map_name(map_name);
  request.mutable_cube_info()->set_program_type(program_type);
  request.mutable_cube_info()->set_index(0);

  request.mutable_get_request()->set_all(false);

  char key_bytes[key_size];
  memcpy(key_bytes, key, key_size);

  request.mutable_get_request()->set_key(key_bytes, key_size);
  request.mutable_get_request()->set_key_size(key_size);
  request.mutable_get_request()->set_value_size(value_size);

  std::cout << " \n key size " << key_size << " ---- value size " << value_size
            << std::endl;
  std::cout << " \n sizeof(key_bytes) " << sizeof(key_bytes) << std::endl;
  std::cout << " \n request.get_request().key().size() "
            << request.get_request().key().size() << std::endl;

  std::cout << "\n\n key_bytes " << std::endl;
  for (auto val : key_bytes)
    printf("\\x%.2x", val);

  std::cout << "\n\n request.get_request().key() " << std::endl;
  for (auto val : request.get_request().key())
    printf("\\x%.2x", val);

  // Container for the data we expect from the server.
  MapValue reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->TableGet(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    return reply;
  } else {
    std::cout << status.error_code() << ": " << status.error_message()
              << std::endl;
    return reply;
  }
}

/*
  Given a key and a value, this method is used to request Polycube to write a
  data in an eBPF map of a cube
*/
Bool PolycubeClient::TableSet(const std::string cube_name,
                              const std::string map_name,
                              const commons::CubeInfo_ProgramType program_type,
                              void *key, size_t key_size, void *value,
                              size_t value_size) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_cube_info()->set_map_name(map_name);
  request.mutable_cube_info()->set_program_type(program_type);
  request.mutable_cube_info()->set_index(0);

  char key_bytes[key_size];
  memcpy(key_bytes, key, key_size);
  char value_bytes[value_size];
  memcpy(value_bytes, value, value_size);

  request.mutable_set_request()->set_key(key_bytes, key_size);
  request.mutable_set_request()->set_key_size(key_size);
  request.mutable_set_request()->set_value(value_bytes, value_size);
  request.mutable_set_request()->set_value_size(value_size);

  std::cout << " \n key size " << key_size << " ---- value size " << value_size
            << std::endl;
  std::cout << " \n strlen(key_bytes) " << strlen(key_bytes)
            << " ---- strlen(value_bytes) " << strlen(value_bytes) << std::endl;
  std::cout << " \n sizeof(key_bytes) " << sizeof(key_bytes)
            << " ---- sizeof(value_bytes) " << sizeof(value_bytes) << std::endl;
  std::cout << " \n sizeof(request.set_request().key()) "
            << sizeof(request.set_request().key())
            << " ---- sizeof(request.set_request().value()) "
            << sizeof(request.set_request().value()) << std::endl;
  std::cout << " \n request.set_request().key().size() "
            << request.set_request().key().size()
            << " ---- request.set_request().value().size() "
            << request.set_request().value().size() << std::endl;

  std::cout << "\n\n key_bytes " << std::endl;
  for (auto val : key_bytes)
    printf("\\x%.2x", val);
  std::cout << "\n\n value_bytes " << std::endl;
  for (auto val : value_bytes)
    printf("\\x%.2x", val);

  std::cout << "\n\n request.set_request().key() " << std::endl;
  for (auto val : request.set_request().key())
    printf("\\x%.2x", val);
  std::cout << "\n\n request.set_request().value() " << std::endl;
  for (auto val : request.set_request().value())
    printf("\\x%.2x", val);

  // Container for the data we expect from the server.
  Bool reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->TableSet(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n OK status.ok() " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n KO status.ok()" << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  Given a key, this method is used to request Polycube to delete a data in an
  eBPF map of a cube
*/
Bool PolycubeClient::TableRemove(
    const std::string cube_name, const std::string map_name,
    const commons::CubeInfo_ProgramType program_type, void *key,
    size_t key_size) {
  // Data we are sending to the server.
  ServicedToDataplaneRequest request;

  request.mutable_serviced_info()->set_service_name(service_name);
  request.mutable_serviced_info()->set_uuid(service_uuid);

  request.mutable_cube_info()->set_cube_name(cube_name);
  request.mutable_cube_info()->set_map_name(map_name);
  request.mutable_cube_info()->set_program_type(program_type);
  request.mutable_cube_info()->set_index(0);

  request.mutable_remove_request()->set_all(false);

  char key_bytes[key_size];
  memcpy(key_bytes, key, key_size);

  request.mutable_remove_request()->set_key(key_bytes, key_size);
  request.mutable_remove_request()->set_key_size(key_size);

  std::cout << " \n key size " << key_size << std::endl;
  std::cout << " \n request.remove_request().key().size() "
            << request.remove_request().key().size() << std::endl;

  std::cout << "\n\n key_bytes " << std::endl;
  for (auto val : key_bytes)
    printf("\\x%.2x", val);

  std::cout << "\n\n request.remove_request().key() " << std::endl;
  for (auto val : request.remove_request().key())
    printf("\\x%.2x", val);

  // Container for the data we expect from the server.
  Bool reply;

  // Context for the client. It could be used to convey extra information to
  // the server and/or tweak certain RPC behaviors.
  ClientContext context;

  // The actual RPC.
  Status status = stub_->TableRemove(&context, request, &reply);

  // Act upon its status.
  if (status.ok()) {
    std::cout << "\n\n " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  } else {
    std::cout << "\n\n " << status.error_code() << ": "
              << status.error_message() << std::endl;
    return reply;
  }
}

/*
  ====================================================================================================
                                        ROUTES MANAGEMENT
  * void RegisterHandler(const std::string path,const std::string http_verb,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  * void RegisterHandlerGet(const std::string path,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  * void RegisterHandlerPost(const std::string path,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  * void RegisterHandlerPatch(const std::string path,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  * void RegisterHandlerDel(const std::string path,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  * void RegisterHandlerOptions(const std::string path,
  ToPolycubed(*method)(ToServiced,PolycubeClient*));
  ====================================================================================================
*/

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandler(const std::string path,
                                     const std::string http_verb,
                                     ToPolycubed (*method)(ToServiced,
                                                           PolycubeClient *)) {
  std::string uppercase_http_verb(http_verb);
  for (auto &c : uppercase_http_verb)
    c = std::toupper(c);
  router.registerPath(service_base_url + "/" + path, uppercase_http_verb);
  // std::cout << "registerpath " << service_base_url+"/"+path << std::endl;
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             uppercase_http_verb)] = method;
}

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandlerGet(
    const std::string path,
    ToPolycubed (*method)(ToServiced, PolycubeClient *)) {
  router.registerPath(service_base_url + "/" + path, "GET");
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             "GET")] = method;
}

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandlerDel(
    const std::string path,
    ToPolycubed (*method)(ToServiced, PolycubeClient *)) {
  router.registerPath(service_base_url + "/" + path, "DELETE");
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             "DELETE")] = method;
}

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandlerPatch(
    const std::string path,
    ToPolycubed (*method)(ToServiced, PolycubeClient *)) {
  router.registerPath(service_base_url + "/" + path, "PATCH");
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             "PATCH")] = method;
}

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandlerPost(
    const std::string path,
    ToPolycubed (*method)(ToServiced, PolycubeClient *)) {
  router.registerPath(service_base_url + "/" + path, "POST");
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             "POST")] = method;
}

/*
  This method is used to register a path and add the handler/method to the map
  (key=template and http verb, value=handler)
*/
void PolycubeClient::RegisterHandlerOptions(
    const std::string path,
    ToPolycubed (*method)(ToServiced, PolycubeClient *)) {
  router.registerPath(service_base_url + "/" + path, "OPTIONS");
  methods_map[std::make_pair(std::string(service_base_url + "/" + path),
                             "OPTIONS")] = method;
}

}  // namespace polycubeclient