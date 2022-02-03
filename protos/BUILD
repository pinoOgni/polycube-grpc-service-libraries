package(default_visibility = ["//visibility:public"])

load("@rules_proto//proto:defs.bzl", "proto_library")
load("//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("//bazel:grpc_build_system.bzl", "grpc_proto_library")


# The following three rules demonstrate the usage of the cc_grpc_library rule in
# in a mode compatible with the native proto_library and cc_proto_library rules.
proto_library(
    name = "commons_proto",
    srcs = ["commons.proto"],
)

cc_proto_library(
    name = "commons_cc_proto",
    deps = [":commons_proto"],
)

cc_grpc_library(
    name = "commons_cc_grpc",
    srcs = [":commonsd_proto"],
    grpc_only = True,
    deps = [":commons_cc_proto"],
)


