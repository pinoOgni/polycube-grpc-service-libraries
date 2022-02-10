package main

import (
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"time"
	"fmt"
	"io/ioutil"
	"unsafe"
)


type rt_k struct {
	netmask_len uint32
	network uint32
}


type rt_v struct {
	port uint32;
	nexthop uint32
	typep uint8;
}

func main() {
	pb.CreatePolycubeClient()
	time.Sleep(1 * time.Second)
	
	// TABLE GET AND TABLE SET
	pb.Subscribe("helloworldgo")
	time.Sleep(4*time.Second)
	ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
	egress := GetDataplaneCode("Helloworld2_dp_egress.c")
	// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
    conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
	pb.CreateCube("h1",ingress,egress,conf)
	time.Sleep(4*time.Second)

	var key uint32 = 0
	var value uint8
	r_get := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
	fmt.Println(r_get)
	time.Sleep(4*time.Second)
	value = 2 // slowpath
	r_set := pb.TableSet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),value,unsafe.Sizeof(value))
	fmt.Println(r_set)  // should be true
	time.Sleep(4*time.Second)
	r_get2 := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
	fmt.Println(r_get2)
	time.Sleep(4*time.Second)
	fmt.Println("the end")
	time.Sleep(100*time.Second)
}



func GetDataplaneCode(filename string) string {
	b, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Print(err)
    }
	return string(b)
}



// TESTING EXAMPLES
/*
	// SUBSCRIBE AND UNSUBSCRIBE METHODS
	
	pb.Subscribe("helloworldgo")
	fmt.Println("ciao")
	time.Sleep(4*time.Second)
	pb.Unsubscribe()


	// CREATE CUBE AND DESTROY CUBE

	ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
	egress := GetDataplaneCode("Helloworld2_dp_egress.c")
	// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
    conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
	pb.CreateCube("h1",ingress,egress,conf)
	time.Sleep(4*time.Second)
	pb.DestroyCube("h1")
	
	// TABLE GET AND TABLE SET
	pb.Subscribe("helloworldgo")
	ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
	egress := GetDataplaneCode("Helloworld2_dp_egress.c")
	// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
    conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
	pb.CreateCube("h1",ingress,egress,conf)
	time.Sleep(4*time.Second)

	var key uint32 = 0
	var value uint8
	r := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
	fmt.Println(r)

	value = 1 // forward
	r := pb.TablesET("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),value,unsafe.Sizeof(value))
	fmt.Println(r)  // should be true

	r := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
	fmt.Println(r)

*/