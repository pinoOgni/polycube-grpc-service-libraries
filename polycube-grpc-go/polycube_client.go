/*
	IMPORTANT LINKS
	* Kind() https://stackoverflow.com/questions/6996704/how-to-check-variable-type-at-runtime-in-go-language
	* grazie a questo funziona https://pkg.go.dev/encoding/binary#example-Write
	* forse riesco a mettere il risultato di typeof per il type assertion https://stackoverflow.com/questions/46840691/golang-error-i-is-not-a-type/46840716
		cioè a fare una cosa del tipo
			var r = reflect.TypeOf(i).Kind()
			fmt.Printf("Other:%v\n", r)

			f, ok := i.(r)
			fmt.Println(f, ok)
		solo che ora come ora mi dice ./prog.go:13:12: r is not a type
	* spiegazione di type assertion https://stackoverflow.com/questions/38816843/explain-type-assertions-in-go
	* type switch https://go.dev/doc/effective_go#type_switch
	* kind https://pkg.go.dev/reflect#Kind
	* qui capisci come lavora una goroutine https://go.dev/play/p/8VE3KBa1Pkn


	* go routine wait https://go.dev/tour/concurrency/5
	* https://go.dev/play/p/S98GjeaGBX0
	* https://gobyexample.com/waitgroups

	MAP
	* https://www.codegrepper.com/code-examples/go/golang+map+functions


	ROUTER
	* https://benhoyt.com/writings/go-routing/
	* https://groups.google.com/g/golang-nuts/c/hbNCHMIA05g


	URL
	* https://pkg.go.dev/net/url
	* https://golangbyexample.com/parse-url-golang/


	* https://github.com/liqotech/liqo/blob/119018e0ec725c1e8f7f4146494d19d2ca948088/cmd/liqo-webhook/main.go#L31
*/

package polycube_grpc_go

import (
	"bytes"
	context "context"
	"encoding/binary"
	"fmt"
	"github.com/pinoOgni/urlpath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
	"unsafe"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const (
	address     = "127.0.0.1:9001"
	defaultName = "world"
	uuid_seed   = 9999
	base_url    = "/polycube/v1/"
)

// used for debugging, remove them later (pinoOgni)
var DebugServiceName string = "helloworldgo" // router
var DebugUuid int32 = 3153
var DebugCubeName string = "h1"        // r1
var DebugMapName string = "action_map" // routing_table


/*
	This structure represents the abse confguration for the creation of a cube of a generic service
	In particular:
		* Loglevel: trace, debug, info, warning, error, critical
		* Shadow: true, false (only for standard cubes)
		* Span: true, false (only for standard cubes)
		* Type: TC, XDP_SKB, XDP_DRV

	For a more detailed description, refer to the Polycube documentation.
*/
type BaseConf struct {
	Loglevel string
	Shadow   bool
	Span     bool
	Type     string
}

/*
	Struct used within the method map
*/
type key struct {
	templateUrl string
	httpVerb    string
}

/*
 	This map is used to keep track of all the methods of the controlplane to which a key is associated
	 given by the set of the url template (i.e. the generic url) and by the HTTP verb
*/
var methods_map = make(map[key]func(ToServiced) ToPolycubed)

/*
	Represents a handler to a method that will be implemented by the developer in the controlplane
	(therefore specific for service) and that is used at the time of a request from the dataplane to
	the controlplane (example: arrival of a packet and print)
*/
var packet_in func(ToServiced)

/*
	This representes the Long lived polycube client with all the field that we need
*/
type LonglivedPolycubeClient struct {
	client PolycubeClient           // it represents the Polycube Client implementation
	conn   *grpc.ClientConn         // simple grpc connection (used to create the client)
	stream Polycube_SubscribeClient // the stream, i.e. the long lived stream between the client and the server
	ctx    context.Context          // context used fo in the Subscribe method
}

/*
	this represent the client, i.e. this particulare instance of the service
	the client live inside the library, in this way the developer which will create
	the logic of the service, i.e. the controlplane, need only to understand which methods to call
	and does not worry about the implementation of the grpc client
*/
var TheClient *LonglivedPolycubeClient

/*
	Thi represents the Serviced with the service name and the uuid, in this way, the developer
	needs to choose only the service name when the Subscribe method is called and so
	the service name and the uuid (created in the Subscribe method) will remain inside the library
*/
type Serviced struct {
	serviceName string
	uuid        int32
	serviceType ServicedInfo_ServiceType
}

/*
	client parameters used to set some timeouts in the client
	almost the same parameters are setted in the server
*/

var TheServiced Serviced

var kacp = keepalive.ClientParameters{
	Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             5 * time.Second, // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,            // send pings even without active streams
}

/*
  ====================================================================================================
                                        SERVICE MANAGEMENT
	* CreatePolycubeClient()
	* func Subscribe(service_name string)
	* func Unsubscribe()
  ====================================================================================================
*/

/*
	This method is used to create a PolycubeClient, setting all the things required by
	gRPC to make the connection bewteen the client and the server

	Note: a single method could be created where the client was created and the Polycubed subscription was made.
	Two methods were preferred to keep things separate.
*/
func CreatePolycubeClient() bool {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("Did not connect ", err)
		return false
	} else {
		fmt.Println("Connected")
		TheClient = &LonglivedPolycubeClient{
			client: NewPolycubeClient(conn), // method of auto-generated library
			conn:   conn,
			stream: nil,
			ctx:    context.Background(),
		}
		if TheClient == nil {
			return false
		} else {
			return true
		}
	}

	// (pinoOgni): TODO investigate if you should put this instruction or a similar one
	// defer conn.Close()

}

/*
	This method is used to make a subscription of the client (Serviced) to the
	server (Polycubed). It creates the long-lived stream which will live for the entire
	life of the Serviced and where ToPolycubed and ToServiced messages are exchanged

*/
func Subscribe(service_name string, service_type ServicedInfo_ServiceType) bool {

	// Background returns a non-nil, empty Context. It is never canceled, has no values,
	// and has no deadline. It is typically used by the main function, initialization,
	// and tests, and as the top-level Context for incoming requests.

	stream, err := TheClient.client.Subscribe(TheClient.ctx)
	TheClient.stream = stream
	if err != nil {
		log.Printf("Subscribe, error creating stream: %v", err)
		return false
	}

	// the service_name is saved internally so the developer does not need to write every time
	TheServiced.serviceName = service_name
	TheServiced.serviceType = service_type

	// only standard cube have the concept of set peer
	// using the SetPeerUserRequest we can setting a peer and the developer of the controlpane
	// does not care about it
	if service_type == ServicedInfo_STANDARD {
		RegisterHandlerPatch(":cubeName/ports/:portName/peer/", SetPeerUserRequest)
	}

	// creation of the uuid
	rand.Seed(time.Now().UnixNano())
	TheServiced.uuid = rand.Int31n(uuid_seed)

	// cube_name := "cube_name"
	request := &ToPolycubed{ServicedInfo: &ServicedInfo{ServiceName: TheServiced.serviceName, ServiceType: &TheServiced.serviceType, Uuid: TheServiced.uuid}, ToPolycubed: &ToPolycubed_Subscribe{Subscribe: true}}

	if err := stream.Send(request); err != nil {
		log.Printf("error sending on stream: %v", err)
		return false
	}

	fmt.Println("ServiceName", request.ServicedInfo.ServiceName)
	fmt.Println("Uuid ", request.ServicedInfo.Uuid)
	fmt.Println("ServiceType ", request.ServicedInfo.ServiceType)

	// (pinoOgni): maybe it will be removed
	// we wait 4 seconds to give Polycube time to set up
	// everything necessary for the new service
	time.Sleep(4 * time.Second)

	// if we are here, everything is ok
	return true
}

/*
	This method is used to unsubscribe from the service
*/
func Unsubscribe() bool {
	ctx := context.Background()
	request := &ServicedInfo{ServiceName: TheServiced.serviceName, Uuid: TheServiced.uuid}
	reply, err := TheClient.client.Unsubscribe(ctx, request)
	if err != nil {
		log.Printf("Unsubscribe failed: %v", err)
		return false
	}
	log.Printf("Reply from Unsubscribe: %s", reply.GetStatus())
	return reply.GetStatus()
}

/*
  ====================================================================================================
                                        STREAM MANAGEMENT
	func ReadTheStreamGoRoutine(wg *sync.WaitGroup)
	func ReadTheStream()

  ====================================================================================================
*/

/*
	This method is used to read the long lived stream waiting for a request
	(a ToServiced message) from Polycubed to arrive.
	After reading the request and understanding what type it is,
	it will call the correct method based on the registration made at the beginning

	This method can be called in two ways from the controlplane:
	1. as a call to a normal method, therefore of the type pb.ReadTheStream () and after
	this call the controlplane will do nothing but wait for some request from the user
	(precisely by reading from the stream)
	2. as goroutine and in that case the controlpalne will be able to do anything else
	while the thread is in charge of reading the stream. Beware, however, that if the method
	that called the goroutine or the main terminates, the launched goroutine also terminates,
	so the developer will have to handle this case.
*/
func ReadTheStreamGoRoutine(wg *sync.WaitGroup) {
	reply := ToPolycubed{}
	fmt.Println("Waiting for something to read from the stream ....")
	for {
		request, err := TheClient.stream.Recv()
		if err != nil {
			fmt.Println("Failed to receive message: %v", err)
			break
		}
		// router logic with method call of the controlplane
		if request.GetRestUserToServiceRequest() != nil {
			requestUrl := request.GetRestUserToServiceRequest().GetUrl()
			requestHttpVerb := request.GetRestUserToServiceRequest().GetHttpVerb()
			for key, method := range methods_map {
				match, ok := urlpath.MatchPath(key.templateUrl, requestUrl)
				if ok == true && requestHttpVerb == key.httpVerb {
					fmt.Println("the url is present in the method_map ")
					cubeName := match.Params["cubeName"]
					request.ServicedInfo.CubeName = &cubeName

					portName := match.Params["portName"]
					request.ServicedInfo.PortName = &portName

					// call the correct controlplane method
					reply = method(*request)
					break
				}
			}
			// writing reply to Polycubed
			TheClient.stream.Send(&reply)

		} else if request.GetDataplaneToServicedPacket() != nil {
			packet_in(*request)
		} else if request.GetDataplaneToServicedRequest() != nil {
			/*
				Still to be implemented, the message has been created in the proto to represent the case of the dataplane 
				requesting something from the service, i.e. the call of a control plane method by the data plane
			*/
			fmt.Println("GetDataplaneToServicedRequest")
		} else {
			fmt.Println("ERROR: message case not implemented!")
			// (pinoOgni): what to do in this case?
		}
		fmt.Println("Waiting for something to read from the stream ....")
	}
	wg.Done()
}

// (pinoOgni): TODO 
func ReadTheStream() {
	for {
		request, err := TheClient.stream.Recv()
		if err != nil {
			fmt.Println("Failed to receive message: %v", err)
			continue
		}
		fmt.Println("Message from Polycubed ", request.GetServicedInfo().GetServiceName())
	}
}

/*
	reply := &ToPolycubed {
		ServicedInfo: &ServicedInfo{
					ServiceName: TheServiced.serviceName ,
					Uuid: TheServiced.uuid},
					ToPolycubed: &ToPolycubed_ServicedToRestUserReply {
						ServicedToRestUserReply: &ServicedToRestUserReply {
							Success: true,
						},
					},
				}
	TheClient.stream.Send(reply)

*/

/*
  ====================================================================================================
                                        CUBE MANAGEMENT
  ====================================================================================================
*/

/*
	This method is used to make a request to create a cube to Polycubed.

	// (pinoOgni): TODO
	To implement
		* extract other information from the conf and not just the name
		(which is not taken from the conf) such as the type of cube, if it is shadow or not ...

	Question: Why is the userRequest needed and this method must be called from the Control Plane?
	Answer: because the Control Plane has to understand and manage the creation and destruction of a cube 
	(Example: vector with names of the cubes)	
*/
func CreateCube(userRequest ToServiced, ingress_code string, egress_code string) *Bool {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// request := &CubeManagement{CubeName: cube_name, ServiceName: &TheServiced.serviceName}
	request := CubeManagement{}
	request.CubeName = userRequest.GetServicedInfo().GetCubeName()
	request.ServiceName = &TheServiced.serviceName
	request.Uuid = &TheServiced.uuid

	// it is a little bit strange but we need to convert from servicedInfo serviceType
	// to cubemanagement servicetype in order to give to polycube the correct service type
	// yes, we can mantain this info in the server but for now this is ok
	if TheServiced.serviceType == 0 {
		serviceType := CubeManagement_STANDARD
		request.ServiceType = &serviceType
	} else {
		serviceType := CubeManagement_TRANSPARENT
		request.ServiceType = &serviceType
	}
	// represents the basis of the configuration of a service with fields common to all services
	conf := GetBaseConf(userRequest)
	request.Conf = &conf
	request.IngressCode = &ingress_code
	request.EgressCode = &egress_code
	reply, err := TheClient.client.CreateCube(ctx, &request)
	if err != nil {
		log.Printf("Could not create cube: %v", err)
	}
	return reply
}

// (pinoOgni): TODO add comments
func GetBaseConf(userRequest ToServiced) string {
	cube_name := userRequest.GetServicedInfo().GetCubeName()
	request_body := userRequest.GetRestUserToServiceRequest().GetRequestBody()

	var baseConf BaseConf
	var stringBaseConf string
	json.Unmarshal([]byte(request_body), &baseConf)
	if baseConf.Loglevel == "" {
		baseConf.Loglevel = "INFO"
	}
	if baseConf.Type == "" {
		baseConf.Type = "TC"
	}
	if TheServiced.serviceType == 0 { 
		stringBaseConf = `{"name":"` + cube_name + `","type":"` + baseConf.Type +  `","loglevel":"` + baseConf.Loglevel + `","shadow":` + strconv.FormatBool(baseConf.Shadow) + `,"span":` + strconv.FormatBool(baseConf.Span) +`}`
	} else if TheServiced.serviceType == 1 {
		stringBaseConf = `{"name":"` + cube_name + `","type":"` + baseConf.Type +  `","loglevel":"` + baseConf.Loglevel + `"}`
	}
	return stringBaseConf
}

// (pinoOgni): TODO
/*
	in the grpc-go examples, in unary requests, a context is done with withtimeout,
	with 1 or 10 seconds. My question is, is it wrong to create a context that always
	lives with Background? Once you exit the Unary method, isn't it automatically deleted?
*/
/*
	This method is used to request the destruction of a cube from Polycubed

*/
func DestroyCube(cube_name string) *Bool {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// request := &CubeManagement{CubeName: cube_name, ServiceName: &TheServiced.serviceName}
	request := CubeManagement{}
	request.CubeName = cube_name
	request.ServiceName = &TheServiced.serviceName
	request.Uuid = &TheServiced.uuid
	reply, err := TheClient.client.DestroyCube(ctx, &request)
	if err != nil {
		log.Printf("Could not destroy cube: %v", err)
	}
	return reply
}

/*
  ====================================================================================================
                                        PORT MANAGEMENT
  * SetPort(cube_name string, port_name string) *Bool
  * DelPort(cube_name string, port_name string) *Bool
  * SetPeer(cube_name string, port_name string, peer_name string) *Bool
  ====================================================================================================
*/

/*
  This method is used to add a port to a cube. The concept of a door and
  therefore of adding a door is different from adding a value to a map of a
  cube. In fact you have to go through Polycubed which has a view of the doors
  that exist. This is why this method was created.
*/
func SetPort(cube_name string, port_name string) *Bool {
	ctx := context.Background()
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo{
		ServiceName: TheServiced.serviceName,
		Uuid:        TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo{
		CubeName: cube_name,
	}

	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
		},
	}
	reply, err := TheClient.client.SetPort(ctx, &request)
	if err != nil {
		log.Printf("Could not use SetPort: %v", err)
	}

	return reply
}

/*
  This method is used to remove a port to a cube. The concept of a door and
  therefore of adding a door is different from adding a value to a map of a
  cube. In fact you have to go through Polycubed which has a view of the doors
  that exist. This is why this method was created.
*/
func DelPort(cube_name string, port_name string) *Bool {
	ctx := context.Background()
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo{
		ServiceName: TheServiced.serviceName,
		Uuid:        TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo{
		CubeName: cube_name,
	}

	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
		},
	}
	reply, err := TheClient.client.DelPort(ctx, &request)
	if err != nil {
		log.Printf("Could not use DelPort: %v", err)
	}
	return reply
}

/*
  This method is used to set a peer for a given port
  This method is called directly by the SetPeerUserRequest method which is present in the library and not by the control plane,
  but it was done for future versions of the library (example, you change the method to read the stream and you want to use the
	SetPeer method in the controlplane for one specific service)
*/

func SetPeer(cube_name string, port_name string, peer_name string) *Bool {
	ctx := context.Background()
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo{
		ServiceName: TheServiced.serviceName,
		Uuid:        TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo{
		CubeName: cube_name,
	}

	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
			Peer: &peer_name,
		},
	}
	reply, err := TheClient.client.SetPeer(ctx, &request)
	if err != nil {
		log.Printf("could not use SetPeer: %v", err)
	}

	return reply
}

/*
	This method is simply to make a call to the SetPeer method

	Logical operation of setting up a peer:
	1. The user executes the command "polycubectl nameServiceNameCube portsNamePort set peer=namePeer"
	2. The request goes through the rest server and arrives at the grpc server (to the dispatcher method)
	3. A ToServiced message is created and written to the stream
	4. The request arrives at the service and is intercepted by the library (by the ReadTheStream or ReadTheStreamGoRoutine method)
	5. The library itself makes a request to Polycube at the SetPeer endpoint of the grpc method

	This "double" request has been chosen (ie user -> Polycubed -> Service -> Polycubed) so that the Service is updated with the peer settings
*/
func SetPeerUserRequest(request ToServiced) ToPolycubed {
	cube_name := request.GetServicedInfo().GetCubeName()
	port_name := request.GetServicedInfo().GetPortName()
	var peer_name string
	// here, I'm sure that in the request body there is the peer name
	if len(request.GetRestUserToServiceRequest().GetRequestBody()) > 0 {
		peer_name = request.GetRestUserToServiceRequest().GetRequestBody()
		peer_name = peer_name[1 : len(peer_name)-1]
	} else {
		peer_name = " "
	}

	
	setPeerReply := SetPeer(cube_name, port_name, peer_name)
	if setPeerReply.GetStatus() == true {
		return SetToPolycubedReply(true, `{"message":"Peer setted"}`)
	} else {
		return SetToPolycubedReply(false, `{"message":"Problem while setting peer"}`)
	}
}


/*
  ====================================================================================================
										MESSAGE SETTING MANAGEMENT
	* func SetToPolycubedReply(success bool, message string) ToPolycubed 
	* func SetTableServicedToDataplaneRequest(cube_name string, map_name string, program_type CubeInfo_ProgramType) ServicedToDataplaneRequest 
  ====================================================================================================
*/


/*
	Simple method to populate a ToPolycubed message where the only important
	fields are Subscribe and Message
*/
func SetToPolycubedReply(success bool, message string) ToPolycubed {
	ret_reply := ToPolycubed{
		ToPolycubed: &ToPolycubed_ServicedToRestUserReply{
			ServicedToRestUserReply: &ServicedToRestUserReply{
				Success: success,
				Message: &message,
			},
		},
	}
	return ret_reply
}


/*
	Simple method to populate a SetTableServicedToDataplaneRequest message 
*/
func SetTableServicedToDataplaneRequest(cube_name string, map_name string, program_type CubeInfo_ProgramType) ServicedToDataplaneRequest {
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo{
		ServiceName: TheServiced.serviceName,
		Uuid:        TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo{
		CubeName:    cube_name,
		MapName:     map_name,
		ProgramType: program_type,
	}
	return request
}




/*
  ====================================================================================================
                                        EBPF MAPS MANAGEMENT
	* func TableGet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value_size uintptr) interface{}
	* func TableSet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value interface{}, value_size uintptr) *Bool
	* func TableRemove(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr) *Bool

	// (pinoOgni): TODO
	TableRemoveAll
	TableGetAll
  ====================================================================================================
*/

/*
  Given a key, this method is used to ask Polycube for data from an eBPF map of
  a cube
*/
func TableGet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value_size uintptr) *MapValue {
	ctx := context.Background()
	request := SetTableServicedToDataplaneRequest(cube_name,map_name,program_type)

	key_bytes := new(bytes.Buffer)
	if ret, ok := key.(interface{}); ok {

		err := binary.Write(key_bytes, binary.LittleEndian, ret)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("Error")
	}


	b, ok := key.(interface{})

	// key_bytes := make([]byte, key_size)
	// (pinoOgni): TODO, first version 
	request.RequestType = &ServicedToDataplaneRequest_GetRequest{
		GetRequest: &GetRequest{
			All:       false,
			Key:       key_bytes.Bytes(),
			KeySize:   uint64(key_size),
			ValueSize: uint64(value_size),
		},
	}
	reply, err := TheClient.client.TableGet(ctx, &request)
	if err != nil {
		log.Printf("could not use TableGet: %v", err)
		return nil
	}

	return reply
}

/*
  Given a key and a value, this method is used to request Polycube to write a
  data in an eBPF map of a cube
*/
func TableSet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value interface{}, value_size uintptr) *Bool {
	ctx := context.Background()
	request := SetTableServicedToDataplaneRequest(cube_name,map_name,program_type)

	key_bytes := new(bytes.Buffer)
	if ret_key, ok_key := key.(interface{}); ok_key {
		err_key := binary.Write(key_bytes, binary.LittleEndian, ret_key)
		if err_key != nil {
			fmt.Println("key binary.Write failed:", err_key)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("Error")
	}

	b_key, ok_key := key.(interface{})

	value_bytes := new(bytes.Buffer)
	if ret_value, ok_value := value.(interface{}); ok_value {
		err_value := binary.Write(value_bytes, binary.LittleEndian, ret_value)
		if err_value != nil {
			fmt.Println("value binary.Write failed:", err_value)
		}
		fmt.Printf("% x", value_bytes.Bytes())
	} else {
		fmt.Println("Error")
	}

	// key_bytes := make([]byte, key_size)
	// (pinoOgni): TODO, first version 
	request.RequestType = &ServicedToDataplaneRequest_SetRequest{
		SetRequest: &SetRequest{
			Key:       key_bytes.Bytes(),
			KeySize:   uint64(key_size),
			Value:     value_bytes.Bytes(),
			ValueSize: uint64(value_size),
		},
	}
	reply, err := TheClient.client.TableSet(ctx, &request)
	if err != nil {
		log.Printf("could not use TableSet: %v", err)
	}
	return reply
}

/*
  Given a key, this method is used to request Polycube to delete a data in an
  eBPF map of a cube
*/
// (pinoOgni): TODO, BUG, method being debugged
func TableRemove(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr) *Bool {
	ctx := context.Background()
	request := SetTableServicedToDataplaneRequest(cube_name,map_name,program_type)

	key_bytes := new(bytes.Buffer)
	if ret_key, ok_key := key.(interface{}); ok_key {

		err_key := binary.Write(key_bytes, binary.LittleEndian, ret_key)
		if err_key != nil {
			fmt.Println("binary.Write failed:", err_key)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("Error")
	}

	b_key, ok_key := key.(interface{})

	// key_bytes := make([]byte, key_size)
	// (pinoOgni): TODO, first version 
	request.RequestType = &ServicedToDataplaneRequest_RemoveRequest{
		RemoveRequest: &RemoveRequest{
			All:     false,
			Key:     key_bytes.Bytes(),
			KeySize: uint64(key_size),
		},
	}
	reply, err := TheClient.client.TableRemove(ctx, &request)
	if err != nil {
		log.Printf("could not use TableRemove: %v", err)
	}

	return reply
}

/*

  ====================================================================================================
                                        ROUTES MANAGEMENT
	* RegisterHandler(path string, httpVerb string, handler func(ToServiced) ToPolycubed)
	* RegisterHandlerGet(path string, handler func(ToServiced) ToPolycubed)
	* RegisterHandlerPost(path string, handler func(ToServiced) ToPolycubed)
	* RegisterHandlerDelete(path string, handler func(ToServiced) ToPolycubed)
	* RegisterHandlerPatch(path string, handler func(ToServiced) ToPolycubed)
	* RegisterHandlerOptions(path string, handler func(ToServiced) ToPolycubed)
  ====================================================================================================

*/

func RegisterHandler(path string, httpVerb string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: strings.ToUpper(httpVerb)}] = handler
}

func RegisterHandlerGet(path string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: "GET"}] = handler
}

func RegisterHandlerPost(path string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: "POST"}] = handler
}

func RegisterHandlerDelete(path string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: "DELETE"}] = handler
}

func RegisterHandlerPatch(path string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: "PATCH"}] = handler
}

func RegisterHandlerOptions(path string, handler func(ToServiced) ToPolycubed) {
	templateUrl := base_url + TheServiced.serviceName + "/" + path
	methods_map[key{templateUrl: templateUrl, httpVerb: "OPTIONS"}] = handler
}

func RegisterHandlerPacketIn(handler func(ToServiced)) {
	packet_in = handler
}

func sendMessage(stream Polycube_SubscribeClient, msg *ToPolycubed) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(msg)
}

func recvMessage(wantErrCode codes.Code) {
	res, err := TheClient.stream.Recv()
	if status.Code(err) != wantErrCode {
		log.Printf("TheClient.stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, wantErrCode)
	}
	if err != nil {
		fmt.Printf("TheClient.stream.Recv() returned expected error %v\n", err)
		return
	}
	fmt.Printf("received message %q\n", res.GetServicedInfo().GetServiceName())
}


/*
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
