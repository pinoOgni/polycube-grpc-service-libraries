package main

import (
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"time"
	"fmt"
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
	//pb.Subscribe("helloworld")
	fmt.Println("ciao")
	time.Sleep(4*time.Second)
	// go pb.ReadTheStream()
	// pb.DestroyCube("h1")
	// var key uint32 = 0
	// var value uint8 
	// r := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
	// fmt.Println(r)
	
	key_rt_k := rt_k{
		netmask_len: 24,
		network: 196618,
	}
	var value_rt_v rt_v
	r1 := pb.TableGet("r1","routing_table",pb.CubeInfo_INGRESS,key_rt_k,unsafe.Sizeof(key_rt_k),unsafe.Sizeof(value_rt_v))
	fmt.Println(r1)

	/*

    rt_k key{
        .netmask_len = 24, 
        .network = 196618, // x0a\x00\x03\x00 is in big endian, 00, 03, 00, 0a is in little endian
                          // 0a ==> 10, 00 ==> 0, 03 ==> 0*16^5 + 3*16^4 ==> 196608 , 00 ==> 0 ====> 10+196608 = 196618
    };
	*/


}
