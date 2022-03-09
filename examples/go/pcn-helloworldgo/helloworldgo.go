package main

import (
	"fmt"
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"sync"
	"time"
	"unsafe"
	"os"
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
	subscribe := pb.Subscribe("helloworldgo",pb.ServicedInfo_STANDARD)
	if subscribe == false {
		fmt.Println("Subscribe error")
		os.Exit(1)
	} 
	
	var wg sync.WaitGroup
	wg.Add(1)
	go pb.ReadTheStreamGoRoutine(&wg)
	

	// GET DATAPLANE CODE
	ingress = pb.GetDataplaneCode("Helloworldgo_dp_ingress.c")
	egress = pb.GetDataplaneCode("Helloworldgo_dp_egress.c")

	
	// ROUTES MANAGEMENTS
	pb.RegisterHandlerDelete(":cubeName/",DestroyCube);
	pb.RegisterHandlerPost(":cubeName/",CreateCube);

	pb.RegisterHandlerPost(":cubeName/ports/:portName/",AddPorts)
	pb.RegisterHandlerDelete(":cubeName/ports/:portName/",DelPort)
	
	pb.RegisterHandlerGet(":cubeName/ports/",ShowPorts)

	pb.RegisterHandlerGet(":cubeName/action/", GetAction)
	pb.RegisterHandlerPatch(":cubeName/action/", SetAction)


	pb.RegisterHandlerOptions("helloworldgo", HelpCommands)

	pb.RegisterHandlerPacketIn(PacketIn)
	
	fmt.Println("the end routes management")

	wg.Wait()
	fmt.Println("the end")
	// time.Sleep(100*time.Second)
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
		if (tableGetReply == nil) {
			return pb.SetToPolycubedReply(false,`{"message":"error while getting the action"}`)		
		} else {
			fmt.Println("tableGetReply: ", tableGetReply)
			var action string
			switch tableGetReply.ByteValues[0] {
			case 0:
				action = `drop`
			case 1:
				action = `slowpath`
			case 2:
				action = `forward`
			default:
				action = `drop`
			}
			return pb.SetToPolycubedReply(true,`{"action":"`+action+`"}`)
		}		
}



/*
	This method is used to SET the action saved in the action_map table 
*/
func SetAction(request pb.ToServiced) pb.ToPolycubed {
	var action uint8
	fmt.Println("ciao ciao ", request.GetRestUserToServiceRequest().GetRequestBody())
	actionString := request.GetRestUserToServiceRequest().GetRequestBody()
	actionString = actionString[1 : len(actionString)-1]
	switch {
	case actionString == `drop`:
		action = 0
	case actionString == `slowpath`:
		action = 1
	case actionString == `forward`:
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

func UpdatePortsMap(request pb.ToServiced) pb.ToPolycubed {
	fmt.Println("UpdatePortsMap")
	cube_name := request.GetServicedInfo().GetCubeName()
	var key uint32 = 0
	var key_pass interface{}
	var value uint16
	var value_pass interface{}
	var ok bool = true
	// we need the capacity of the ports array not the len
	// because in the case of del ports this will be an error in the case of only one port present in the array
	for key = uint32(0); key < uint32(cap(cubes[cube_name])); key++ {
		key_pass = &key
		fmt.Println("UpdatePortsMap first for", key)
		// (pinoOgni) code taken from classic helloworld, could be improved
		value =  uint16(key) + uint16(100) // like port->index()  
		value_pass = &value
		tableSetReply := pb.TableSet(cube_name,"ports_map",pb.CubeInfo_INGRESS,key_pass,unsafe.Sizeof(key),value_pass,unsafe.Sizeof(value))
		if tableSetReply.GetStatus() == false {
			ok = false
			break
		}
	}

	value = ^uint16(0)
	for key < 2 && ok == true {
		fmt.Println("UpdatePortsMap second for", key)
		tableSetReply := pb.TableSet(cube_name,"ports_map",pb.CubeInfo_INGRESS,key_pass,unsafe.Sizeof(key),value_pass,unsafe.Sizeof(value))
		key++
		if tableSetReply.GetStatus() == false {
			ok = false
			break
		}
	}

	if ok == true {
		return pb.SetToPolycubedReply(true,"OK")
	} else {
		return pb.SetToPolycubedReply(false,"KO")
	}
}

// pino: TODO find a better solution to fill the ToPolycubed message
func AddPorts(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	fmt.Println("AddPorts method")
	if len(cubes[cube_name]) >= 2 {
		fmt.Println("Maximum number of ports reached")
		// pino: TODO the message is not displayed by the termianl of polycubectl
		return pb.SetToPolycubedReply(false,`{"message":"Maximum number of ports reached"}`)
	}
	port_name := request.GetServicedInfo().GetPortName()
	
	// ports = append(ports, port_name)
	cubes[cube_name] = append(cubes[cube_name], port_name)
	fmt.Println("AddPorts method ", cubes[cube_name][0])
	setPortReply := pb.SetPort(cube_name,port_name)
	fmt.Println("AddPorts method setPortReply ", setPortReply.GetStatus())
	if setPortReply.GetStatus() == true {
		return UpdatePortsMap(request)
	} else {
		return pb.SetToPolycubedReply(false,`{"message":"Problem while adding port"}`)
	}
}


func DelPort(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	fmt.Println("DelPort method")
	if len(cubes[cube_name]) == 0 {
		fmt.Println("There are no port to be deleted")
		// pino: TODO the message is not displayed by the terminal of polycubectl
		return pb.SetToPolycubedReply(false,`{"message":"There are no port to be deleted"}`)
	}

	var pos = cap(cubes[cube_name])
	port_name := request.GetServicedInfo().GetPortName()
	for i := range cubes[cube_name] {
		if cubes[cube_name][i] == port_name {
			pos = i
		}
	}
	if pos == cap(cubes[cube_name]) {
		fmt.Println("The port does not exists!")
		// pino: TODO the message is not displayed by the terminal of polycubectl
		return pb.SetToPolycubedReply(false,`{"message":"The port does not exists!"}`)
	}

	// controlplane level 
	// Remove the element at index pos from a.
	cubes[cube_name][pos] = cubes[cube_name][len(cubes[cube_name])-1] // Copy last element to index pos.
	cubes[cube_name][len(cubes[cube_name])-1] = ""   // Erase last element (write zero value).
	cubes[cube_name] = cubes[cube_name][:len(cubes[cube_name])-1]   // Truncate slice.

	delPortReply := pb.DelPort(cube_name,port_name)

	if delPortReply.GetStatus() == true {
		return UpdatePortsMap(request)
	} else {
		return pb.SetToPolycubedReply(false,`{"message":"Problem while adding port"}`)
	}
}


/*
	This method has nothing to do with the server, it only concerns the controlplane 
	and above all reflects the logic of the service and how the final developer wants 
	to implement the controlplane. In this case we have to do with the concept of Port

	The only thing it has to respect is the method signature in order 
	to receive a request from the user and respond
*/
func ShowPorts(request pb.ToServiced) pb.ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	if len(cubes[cube_name]) != 0 {
		message := ""
		for _, p := range cubes[cube_name] {
			message += p + " "
		} 
		return pb.SetToPolycubedReply(true,`{"ports":"`+message+`"}`)
	} else {
		message := "There are no ports on the cube " + cube_name 
		return pb.SetToPolycubedReply(false,`{"ports":"`+message+`"}`)
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






func PacketIn(request pb.ToServiced) {
	fmt.Println("Packet in HelloworldGo")
	cubeId := request.GetDataplaneToServicedPacket().GetCubeId()
	portId := request.GetDataplaneToServicedPacket().GetPortId()
	fmt.Println("======= PACKET FROM " , cubeId, " and port: ", portId)
}