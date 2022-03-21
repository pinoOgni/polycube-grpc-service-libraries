package main

import (
	"fmt"
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"sync"
	"os"
)

/*
	Definition of the ingress and egress code of the service
*/
var ingress string
var egress string

/*
	Developer's choice regarding port management and also cube management
*/

/*
	example: map with key given by the name of the cube and with value given by a vector of strings 
	representing the names of the ports of the service (Helloworld service)
*/
var cubes = make(map[string][]string)
// the slice need to be initialized after the creation of the cube in this way
// cubes[cube_name] = make([]string, 0, 2)

func main() {

/*
  ====================================================================================================
                                        CREATION OF POLYCUBE CLIENT
	The first thing to do is to call this method.
	It is used to create a Polycube Client object
  ====================================================================================================
*/
	client := pb.CreatePolycubeClient()

	/*
		It is not mandatory but recommended. 
		Check the return value of the CreatePolycubeClient
	*/
	if client == false {
		fmt.Println("CreatePolycubeClient error")
		os.Exit(1)
	}
/*
  ====================================================================================================
                                        SUBSCRIBE
	In this part we have to register the service with Polycube.
	The only two things to do are to decide a name and type of service (transparent or standard)
	and as before, check the return value of the Subscribe 
  ====================================================================================================
*/
	subscribe := pb.Subscribe("serviceName",pb.ServicedInfo_STANDARD)
	if subscribe == false {
		fmt.Println("Subscribe error")
		os.Exit(1)
	} 
/*
  ====================================================================================================
                                        GET DATAPLANE CODE
	In this part we have to register the service with Polycube.
	Here WE get the Data Plane code, using a library method that only needs the file name.
	It is valid for both the ingress and egress part. 
  ====================================================================================================
*/
	ingress = pb.GetDataplaneCode("Helloworldgo_dp_ingress.c")
	egress = pb.GetDataplaneCode("Helloworldgo _dp_egress.c")
/*
  ====================================================================================================
                                        ROUTE REGISTRATION: WORK IN PROGRESS
	Here WE have to register the controlplane methods at each URL.
	It must be taken into account that the base url to start from is polycube/v1/ServiceName/ 
  ====================================================================================================
*/
	/*
		First, let's record the methods for creating and destroying a cube
		These two methods must be implemented in the Control Plane in order to understand which 
		and how many cubes of this service are running
	*/
	pb.RegisterHandlerDelete(":cubeName/",DestroyCube);
	pb.RegisterHandlerPost(":cubeName/",CreateCube);

	/*
		Here we can register the specific methods of the service. For example, if it were the Helloworld 
		service, we could register the GetAction and the SetAction
	*/
	pb.RegisterHandlerGet(":cubeName/action/", GetAction)
	pb.RegisterHandlerPatch(":cubeName/action/", SetAction)

	/*
		Here we need to register the PacketIn method. 
		It is not linked to any url, so please pass in the name of the method that will be implemented here.
	*/
	pb.RegisterHandlerPacketIn(PacketIn)
/*
  ====================================================================================================
                                        READ THE STREAM GOROUTINE
	Here a goroutine is launched, i.e. a light thread.
	The method used is called ReadTheStreamGoRoutine and is used to read the stream where requests from Polycube are written.
	By running it as a thread you can do other things between launch and wait. 
  ====================================================================================================
*/	
	var wg sync.WaitGroup
	wg.Add(1)
	go pb.ReadTheStreamGoRoutine(&wg)
	

	/*
		Here the wait of the previous goroutine is simply done
	*/
	wg.Wait()
	fmt.Println("End of the skeleton")
}


/*

	This represents a simple implementation of the CreateCube. The only things it needs to do are:
		* Have this signature in order to be registered in the method map
		* Implement the management of multiple cubes as it sees fit
		* Call the CreateCube method of the library, passing the appropriate parameters: the request, the ingress and egress code
		* Based on the response from CreateCube, use the SetToPolycubedReply method to respond to the user
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

	This represents a simple implementation of the DestroyCube. The only things it needs to do are:
		* Have this signature in order to be registered in the method map
		* Implement the management of multiple cubes as it sees fit
		* Call the DestroyCube method of the library, passing the appropriate parameters: cube name
		* Based on the response from DestroyCube, use the SetToPolycubedReply method to respond to the user
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

/*
	This represents a simple implementation of the PacketIn. The only things it needs to do are:
		* Have this signature in order to be registered in the method map
		* Implement the logic you want (for example print the packet)
*/
func PacketIn(request pb.ToServiced) {
	fmt.Println("Packet on the way...")
	cubeId := request.GetDataplaneToServicedPacket().GetCubeId()
	fmt.Println("PACKET FROM " , cubeId)
}


func GetAction(request pb.ToServiced) pb.ToPolycubed {
	return pb.ToPolycubed {}
}

func SetAction(request pb.ToServiced) pb.ToPolycubed {
	return pb.ToPolycubed {}
}
/*
  ====================================================================================================
                                        WORK IN PROGRESS
  ====================================================================================================
*/

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



