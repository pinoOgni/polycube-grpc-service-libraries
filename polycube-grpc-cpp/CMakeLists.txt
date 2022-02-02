cmake_minimum_required(VERSION 3.5.1)

project(polycube-grpc-cpp)

add_subdirectory(include)

# (pinoOgni)
install(DIRECTORY include/ DESTINATION /usr/include)

