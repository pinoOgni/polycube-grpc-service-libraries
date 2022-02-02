cmake_minimum_required(VERSION 3.5.1)
project (pcn-helloworld2)
include(cmake/common.cmake)

add_subdirectory(src)


# Target client
# foreach(_target
#  first_client router_client second_client)
#  add_executable(${_target} "${_target}.cc")
#  target_link_libraries(${_target}
#    polycube-grpc
#    ${_REFLECTION}
#    ${_GRPC_GRPCPP}
#    ${_PROTOBUF_LIBPROTOBUF})
# endforeach()