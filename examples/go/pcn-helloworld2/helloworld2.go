package main

import (
	"fmt"
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"io/ioutil"
	"sync"
	"time"
	"unsafe"
)

type rt_k struct {
	netmask_len uint32
	network     uint32
}

type rt_v struct {
	port    uint32
	nexthop uint32
	typep   uint8
}

var ingress string
var egress string

func main() {
	pb.CreatePolycubeClient()
	time.Sleep(1 * time.Second)

	// SUBSCRIBE
	pb.Subscribe("helloworldgo")
	time.Sleep(4 * time.Second)

	// GET DATAPLANE CODE
	ingress = GetDataplaneCode("Helloworld2_dp_ingress.c")
	egress = GetDataplaneCode("Helloworld2_dp_egress.c")

	// testing urlpath router in the ReadTheStream method
	pb.RegisterHandlerGet(":cubeName/action/", GetAction)

	pb.RegisterHandlerPost(":cubeName/", CreateCube)

	// pino: TODO se lo metto prima e dopo faccio qualcosa mi da errore e crasha
	var wg sync.WaitGroup
	wg.Add(1)
	go pb.ReadTheStreamGoRoutine(&wg)
	wg.Wait()

	fmt.Println("the end")
	// time.Sleep(100*time.Second)
}

func first(a string, b string) string {

	return a + " --first-- " + b
}

/*
	This method is used to get the action saved in the action_map table
*/
func GetAction(request pb.ToServiced) pb.ToPolycubed {

	// pino: TODO // the key is setted by the controlplane? In this case yes, but in others?
	var key uint32 = 0
	var value uint8
	fmt.Println("request.GetServicedInfo().GetCubeName() ", request.GetServicedInfo().GetCubeName())
	tableGetReply := pb.TableGet(request.GetServicedInfo().GetCubeName(), "action_map", pb.CubeInfo_INGRESS, key, unsafe.Sizeof(key), unsafe.Sizeof(value))
	fmt.Println("tableGetReply: ", tableGetReply)

	message := `{"action":"drop"}`
	//   ret_reply.mutable_serviced_to_rest_user_reply()->set_message("{\"action\":\"drop\"}");
	ret_reply := pb.ToPolycubed{
		ToPolycubed: &pb.ToPolycubed_ServicedToRestUserReply{
			ServicedToRestUserReply: &pb.ServicedToRestUserReply{
				Success: true,
				Message: &message,
			},
		},
	}
	return ret_reply
}

/*
	Controlplane method: is an example of how Dataplane code can be obtained

	So, it is used to obtain a string given the name of a file that contains the dataplane
	code (eBPF) of the new service. Once this is achieved, the string can be passed to the
	CreateCube method of the polycube-grpc-go library
*/

func GetDataplaneCode(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func CreateCube(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
	conf := `{"name":"` + cube_name + `","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
	create_cube_reply := pb.CreateCube(cube_name, ingress, egress, conf)
	if create_cube_reply == true {
		message := "CREATED"
		ret_reply := pb.ToPolycubed{
			ToPolycubed: &pb.ToPolycubed_ServicedToRestUserReply{
				ServicedToRestUserReply: &pb.ServicedToRestUserReply{
					Success: true,
					Message: &message,
				},
			},
		}
		return ret_reply
	} else {
		message := "NOT CREATED"
		ret_reply := pb.ToPolycubed{
			ToPolycubed: &pb.ToPolycubed_ServicedToRestUserReply{
				ServicedToRestUserReply: &pb.ServicedToRestUserReply{
					Success: false,
					Message: &message,
				},
			},
		}
		return ret_reply
	}

}
