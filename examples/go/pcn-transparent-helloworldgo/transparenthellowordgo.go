package main

import (
	"fmt"
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"io/ioutil"
	"sync"
	"time"
	"unsafe"
	"os"
)

var ingress string
var egress string

// developer's choice regarding port management and also cube management
var cubes = make(map[string][]string)
// the slice need to be initialized after the creation of the cube in this way
// cubes[cube_name] = make([]string, 0, 2)

func main() {
	client := pb.CreatePolycubeClient()
	if client == false {
		fmt.Println("CreatePolycubeClient error")
		os.Exit(1)
	}
	time.Sleep(1 * time.Second)

	// SUBSCRIBE
	subscribe := pb.Subscribe("transparenthelloworldgo",pb.ServicedInfo_TRANSPARENT)
	if subscribe == false {
		fmt.Println("Subscribe error")
		os.Exit(1)
	} 
	
	var wg sync.WaitGroup
	wg.Add(1)
	go pb.ReadTheStreamGoRoutine(&wg)
	

	// GET DATAPLANE CODE
	ingress = GetDataplaneCode("TransparentHelloworldgo_dp_ingress.c")
	egress = GetDataplaneCode("TransparentHelloworldgo_dp_egress.c")

	
	// ROUTES MANAGEMENTS
	pb.RegisterHandlerDelete(":cubeName/",DestroyCube)
	pb.RegisterHandlerPost(":cubeName/",CreateCube)

	pb.RegisterHandlerGet(":cubeName/ingress-action/", GetIngressAction)
	pb.RegisterHandlerGet(":cubeName/egress-action/", GetEgressAction)
	pb.RegisterHandlerPatch(":cubeName/ingress-action/", SetIngressAction)
	pb.RegisterHandlerPatch(":cubeName/egress-action/", SetEgressAction)



	pb.RegisterHandlerOptions("helloworldgo",HelpCommands)

	pb.RegisterHandlerPacketIn(PacketIn)
	
	fmt.Println("the end routes management")

	wg.Wait()
	fmt.Println("the end")
	// time.Sleep(100*time.Second)
}

func first(a string, b string) string {

	return a + " --first-- " + b
}

/*
	This method is used to get the action saved in the action_map table of the INGRESS CODE
*/
func GetIngressAction(request pb.ToServiced) pb.ToPolycubed {

	// (pinoOgni) TODO // the key is setted by the controlplane? In this case yes, but in others?
	var key uint32 = 0
	var value uint8
	fmt.Println("request.GetServicedInfo().GetCubeName() ", request.GetServicedInfo().GetCubeName())
	tableGetReply := pb.TableGet(request.GetServicedInfo().GetCubeName(), "action_map", pb.CubeInfo_INGRESS, key, unsafe.Sizeof(key), unsafe.Sizeof(value))
		if (tableGetReply == nil) {
			return pb.SetToPolycubedReply(false,`{"message":"error while getting the action"}`)		
		} else {
			fmt.Println("tableGetReply: ", tableGetReply)
			var action string
			switch tableGetReply.ByteValues[0] {
			case 0:
				action = `drop`
			case 1:
				action = `pass`
			case 2:
				action = `slowpath`
			default:
				action = `drop`
			}
			return pb.SetToPolycubedReply(true,`{"action":"`+action+`"}`)
		}		
}

/*
	This method is used to get the action saved in the action_map table of the EGRESS CODE
*/
func GetEgressAction(request pb.ToServiced) pb.ToPolycubed {

	// (pinoOgni) TODO // the key is setted by the controlplane? In this case yes, but in others?
	var key uint32 = 0
	var value uint8
	fmt.Println("request.GetServicedInfo().GetCubeName() ", request.GetServicedInfo().GetCubeName())
	tableGetReply := pb.TableGet(request.GetServicedInfo().GetCubeName(), "action_map", pb.CubeInfo_EGRESS, key, unsafe.Sizeof(key), unsafe.Sizeof(value))
		if (tableGetReply == nil) {
			return pb.SetToPolycubedReply(false,`{"message":"error while getting the action"}`)		
		} else {
			fmt.Println("tableGetReply: ", tableGetReply)
			var action string
			switch tableGetReply.ByteValues[0] {
			case 0:
				action = `drop`
			case 1:
				action = `pass`
			case 2:
				action = `slowpath`
			default:
				action = `drop`
			}
			return pb.SetToPolycubedReply(true,`{"action":"`+action+`"}`)
		}		
}

/*
	This method is used to SET the action saved in the action_map table of the INGRESS CODE
*/
func SetIngressAction(request pb.ToServiced) pb.ToPolycubed {
	var action uint8
	actionString := request.GetRestUserToServiceRequest().GetRequestBody()
	switch {
	case actionString == `"drop"`:
		action = 0
	case actionString == `"pass"`:
		action = 1
	case actionString == `"slowpath"`:
		action = 2
	default: {
		action = 0
	}
		
	}

	var key uint32 = 0
	// var value uint8 = 1
	tableSetReply := pb.TableSet(request.GetServicedInfo().GetCubeName(),"action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),action,unsafe.Sizeof(action))
	fmt.Println("request.GetServicedInfo().GetCubeName() ", request.GetServicedInfo().GetCubeName())
	fmt.Println("tableSetReply: ", tableSetReply)
	if tableSetReply.GetStatus() == true {
		return pb.SetToPolycubedReply(true,"ACTION SETTED")
	} else {
		return pb.SetToPolycubedReply(false,"ACTION NOT SETTED")
	}
}



/*
	This method is used to SET the action saved in the action_map table of the EGRESS CODE
*/
func SetEgressAction(request pb.ToServiced) pb.ToPolycubed {
	var action uint8
	actionString := request.GetRestUserToServiceRequest().GetRequestBody()
	switch {
	case actionString == `"drop"`:
		action = 0
	case actionString == `"pass"`:
		action = 1
	case actionString == `"slowpath"`:
		action = 2
	default: {
		action = 0
	}
		
	}

	var key uint32 = 0
	// var value uint8 = 1
	tableSetReply := pb.TableSet(request.GetServicedInfo().GetCubeName(),"action_map",pb.CubeInfo_EGRESS,key,unsafe.Sizeof(key),action,unsafe.Sizeof(action))
	fmt.Println("request.GetServicedInfo().GetCubeName() ", request.GetServicedInfo().GetCubeName())
	fmt.Println("tableSetReply: ", tableSetReply)
	if tableSetReply.GetStatus() == true {
		return pb.SetToPolycubedReply(true,"ACTION SETTED")
	} else {
		return pb.SetToPolycubedReply(false,"ACTION NOT SETTED")
	}
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


/*

	This method is used to create a helloworldgo type cube, it internally 
	calls the CreateCube of the polycube-grpc-go library

*/
func CreateCube(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()

	// simple control with cubes map
	if _, ok := cubes[cube_name]; ok {
		return pb.SetToPolycubedReply(false,cube_name + " already created")
	}
	create_cube_reply := pb.CreateCube(request, ingress, egress)
	if create_cube_reply.GetStatus() == true {
		// initialize the ports slice for this cube
		cubes[cube_name] = make([]string, 0, 2)
		return pb.SetToPolycubedReply(true,"CREATED")
	} else {
		return pb.SetToPolycubedReply(false,"NOT CREATED")
	}
}



/*

	This method is used to destroy a helloworldgo type cube, it internally 
	calls the DestroyCube of the polycube-grpc-go library
	
*/
func DestroyCube(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	destroy_cube_reply := pb.DestroyCube(cube_name)
	if destroy_cube_reply.GetStatus() == true {
		// destroy element from the array of cubes
		delete(cubes, cube_name)
		return pb.SetToPolycubedReply(true,"DESTROYED")
	} else {
		return pb.SetToPolycubedReply(false,"NOT DESTROYED")
	}
}



// (pinoOgni)
/*
In theory it works, but:
1. the message is not displayed by polycubectl
2. if I try to put all the helloworld help message, polycubectl gives an error while parsing
*/
func HelpCommands(request pb.ToServiced) pb.ToPolycubed {
	message := `{"commands":"add, del, show"}`
	return pb.SetToPolycubedReply(true,message)
}




func PacketIn(request pb.ToServiced){
	fmt.Println("Packet in TransparentHelloworldGo")
	cubeId := request.GetDataplaneToServicedPacket().GetCubeId()
	portId := request.GetDataplaneToServicedPacket().GetPortId()
	fmt.Println("======= PACKET FROM " , cubeId, " and port: ", portId)

}