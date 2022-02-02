cmake_minimum_required(VERSION 3.5.1)
include(../cmake/common.cmake)

# Proto file
get_filename_component(pc_proto "../../protos/polycube_grpc_commons.proto" ABSOLUTE)
get_filename_component(pc_proto_path "${pc_proto}" PATH)

# Generated sources
set(pc_proto_srcs "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp/polycube_grpc_commons.pb.cc")
set(pc_proto_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp/polycube_grpc_commons.pb.h")
set(pc_grpc_srcs "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp/polycube_grpc_commons.grpc.pb.cc")
set(pc_grpc_hdrs "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp/polycube_grpc_commons.grpc.pb.h")

message(STATUS "CMAKE_SOURCE_DIR = ${CMAKE_SOURCE_DIR}")
message(STATUS "{CMAKE_CURRENT_SOURCE_DIR} = ${CMAKE_CURRENT_SOURCE_DIR}")
message(STATUS "{CMAKE_CURRENT_BINARY_DIR} = ${CMAKE_CURRENT_BINARY_DIR}")

add_custom_command(
      OUTPUT "${pc_proto_srcs}" "${pc_proto_hdrs}" "${pc_grpc_srcs}" "${pc_grpc_hdrs}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${pc_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${pc_proto}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp"
        --cpp_out "${CMAKE_CURRENT_SOURCE_DIR}/polycube-grpc-cpp"
        -I "${pc_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${pc_proto}"
      DEPENDS "${pc_proto}")

message(STATUS "pc_proto_srcs = ${pc_proto_srcs}")
message(STATUS "pc_proto_hdrs = ${pc_proto_hdrs}")
message(STATUS "pc_grpc_srcs = ${pc_grpc_srcs}")
message(STATUS "pc_grpc_hdrs = ${pc_grpc_hdrs}")
message(STATUS "CMAKE_INSTALL_LIBDIR = ${CMAKE_INSTALL_LIBDIR}")


# Include generated *.pb.h files

include_directories(${CMAKE_CURRENT_BINARY_DIR})

# client library sources and header

set(polycube_client_srcs "${CMAKE_SOURCE_DIR}/include/polycube-grpc-cpp/polycube_client.cc")
set(polycube_client_hdrs "${CMAKE_SOURCE_DIR}/include/polycube-grpc-cpp/polycube_client.h")

set(simple_cpp_router_hdrs 
    "${CMAKE_SOURCE_DIR}/include/Routing/Router.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/PathMatch.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/Exceptions.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/PathTemplate.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/RegexConverter.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/PathTemplateParser.h"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/StringTokenizer.h"
)

set(simple_cpp_router_srcs
    "${CMAKE_SOURCE_DIR}/include/Routing/Router.cpp"
    "${CMAKE_SOURCE_DIR}/include/Routing/PathMatch.cpp"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/PathTemplate.cpp"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/RegexConverter.cpp"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/PathTemplateParser.cpp"
    "${CMAKE_SOURCE_DIR}/include/Routing/Details/StringTokenizer.cpp"
)

# polycube-grpc-cpp
add_library(polycube-grpc-cpp STATIC
    ${simple_cpp_router_srcs}
    ${simple_cpp_router_hdrs}
    ${pc_grpc_srcs} 
    ${pc_grpc_hdrs}  
    ${pc_proto_srcs} 
    ${pc_proto_hdrs} 
    ${polycube_client_srcs} 
    ${polycube_client_hdrs})

set_property(TARGET polycube-grpc-cpp PROPERTY POSITION_INDEPENDENT_CODE ON)

target_link_libraries(polycube-grpc-cpp
    ${_REFLECTION} 
    ${_GRPC_GRPCPP} 
    ${_PROTOBUF_LIBPROTOBUF})


# Install libpolycube-grpc-cpp, it is used by grpc services that are compiled standalone
set(CMAKE_INSTALL_LIBDIR /usr/lib)

install(
    TARGETS
    polycube-grpc-cpp
    DESTINATION
    "${CMAKE_INSTALL_LIBDIR}"
)