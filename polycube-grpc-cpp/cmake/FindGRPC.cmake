# cmake build file for C++ gRPC server.
# Assumes protobuf and gRPC have been installed using cmake.

cmake_minimum_required(VERSION 3.16.3)

find_package(Threads REQUIRED)


# This branch assumes that gRPC and all its dependencies are already installed
# on this system, so they can be located by find_package().

# Find Protobuf installation
# Looks for protobuf-config.cmake file installed by Protobuf's cmake installation.
set(protobuf_MODULE_COMPATIBLE TRUE)

find_package(Protobuf REQUIRED) # To solve problem with CMAKE_PREFIX_PATH or Protobuf_DIR
message(STATUS "Using protobuf ${Protobuf_VERSION}")
message(STATUS "Protobuf_FOUND ${Protobuf_FOUND}")
message(STATUS "Protobuf_INCLUDE_DIR ${Protobuf_INCLUDE_DIR}")

set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_REFLECTION gRPC::grpc++_reflection)
if(CMAKE_CROSSCOMPILING)
  find_program(_PROTOBUF_PROTOC protoc)
else()
  set(_PROTOBUF_PROTOC $<TARGET_FILE:protobuf::protoc>)
endif()


# Find gRPC installation
# Looks for gRPCConfig.cmake file installed by gRPC's cmake installation.
find_package(gRPC CONFIG REQUIRED) # To solve with CMAKE_PREFIX_PATH or grpc_DIR
message(STATUS "Using gRPC ${gRPC_VERSION}")

set(_GRPC_GRPCPP gRPC::grpc++)
if(CMAKE_CROSSCOMPILING)
  find_program(_GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin)
else()
  set(_GRPC_CPP_PLUGIN_EXECUTABLE $<TARGET_FILE:gRPC::grpc_cpp_plugin>)
endif()



