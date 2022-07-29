include(${PROJECT_SOURCE_DIR}/cmake/LoadFileAsVariable.cmake)
include(../cmake/common.cmake)

# Needed to load files as variables
include_directories(${CMAKE_CURRENT_BINARY_DIR})

# Target client
foreach(_target
  helloworldcpp_client)
  add_executable(${_target} "${_target}.cc")
  target_link_libraries(${_target}
    polycube-grpc-cpp
    ${_REFLECTION}
    ${_GRPC_GRPCPP}
    ${_PROTOBUF_LIBPROTOBUF})
endforeach()
