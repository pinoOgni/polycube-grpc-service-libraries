/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"flag"
	"fmt"
	
	"log"
	
	"net"
	"unsafe"
	"io"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"


	pb "polycube-grpc-go/commons"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 9001, "The server port")
)

type polycubeServer struct {
	pb.UnimplementedPolycubeServer
}


func (s *polycubeServer) Subscribe(stream pb.Polycube_SubscribeServer) error {
	fmt.Println("SUBSCRIBE")
	for {
		in, err := stream.Recv()
		fmt.Println("unsafe.Sizeof() ", unsafe.Sizeof(in))
		if err == io.EOF {
			fmt.Println("err == io.EOF")
			return nil
		}
		if err != nil {
			fmt.Println("err != nil")
			return err
		}
		fmt.Println("ciao ", in.GetServicedInfo().GetServiceName())
		fmt.Println("ciao ", in.GetServicedInfo().GetCubeName())

		reply := &pb.ToServiced {ServicedInfo: &pb.ServicedInfo{ServiceName: "NOME", Uuid: 1111111}}

		if err := stream.Send(reply); err != nil {
				return err
		}
	}
}


func newServer() *polycubeServer {
	s := &polycubeServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPolycubeServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

