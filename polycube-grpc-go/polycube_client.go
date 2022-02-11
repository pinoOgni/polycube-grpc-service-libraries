/*
	IMPORTANT LINKS
	* occhio a Kind() https://stackoverflow.com/questions/6996704/how-to-check-variable-type-at-runtime-in-go-language
	* grazie a questo funziona https://pkg.go.dev/encoding/binary#example-Write
	* occhio, forse riesco a mettere il risultato di typeof per il type assertion https://stackoverflow.com/questions/46840691/golang-error-i-is-not-a-type/46840716
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
*/



package polycube_grpc_go 

import (
	context "context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"unsafe"
	"net/http"
	"strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"github.com/go-chi/chi/v5"
	"bytes"
	"encoding/binary"
	"sync"
)

const (
	address     = "127.0.0.1:9001"
	defaultName = "world"
	uuid_seed = 9999
)

var DebugServiceName string = "helloworldgo"
var DebugUuid int32 = 3153
var DebugCubeName string = "h1"
var DebugMapName string = "action_map" // routing_table

//     auto tableGetReply = polycube.TableGet("router", 1111, "r1", "routing_table", program_type,key_pass,sizeof(key),sizeof(value));

/*
	This representes the Long lived polycube client with all the field that we need
*/
type LonglivedPolycubeClient struct {
	client PolycubeClient   // it represents the Polycube Client implementation
	conn *grpc.ClientConn  // simple grpc connection (used to create the client)
	stream Polycube_SubscribeClient // the stream, i.e. the long lived stream between the client and the server
	ctx context.Context  // context used fo in the Subscribe method
	router *chi.Mux
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
	uuid int32
}
/*
	client parameters used to set some timeouts in the client
	almost the same parameters are setted in the server
*/

var TheServiced Serviced

var kacp = keepalive.ClientParameters{
	Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             5 * time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
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
func CreatePolycubeClient(){

	fmt.Println("1")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	fmt.Println("2")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		fmt.Println("connected")
	}

	// pino: TODO OCCHIO, non so se metterla o no
	// defer conn.Close()

	TheClient = &LonglivedPolycubeClient{
		client: NewPolycubeClient(conn), // method of auto-generated library
		conn:   conn,
		stream: nil,
		ctx: context.Background(),
		router: chi.NewRouter(),
	} 
	
}



/*
	This method is used to make a subscription of the client (Serviced) to the 
	server (Polycubed). It creates the long-lived stream which will live for the entire
	life of the Serviced and where ToPolycubed and ToServiced messages are exchanged

*/
func Subscribe(service_name string) {

	// Background returns a non-nil, empty Context. It is never canceled, has no values,
	// and has no deadline. It is typically used by the main function, initialization, 
	// and tests, and as the top-level Context for incoming requests. 
	
	stream, err := TheClient.client.Subscribe(TheClient.ctx)
	TheClient.stream = stream
	if err != nil {
		log.Fatalf("Subscribe, error creating stream: %v", err)
	}
	
	// the service_name is saved internally so the developer does not need to write every time
	TheServiced.serviceName = service_name 
	
	// creation of the uuid
	rand.Seed(time.Now().UnixNano())
	TheServiced.uuid = rand.Int31n(uuid_seed)
	
	// cube_name := "cube_name"
	request := &ToPolycubed {ServicedInfo: &ServicedInfo{ServiceName: TheServiced.serviceName , Uuid: TheServiced.uuid}, ToPolycubed: &ToPolycubed_Subscribe {Subscribe: true}}
	
	if err := stream.Send(request); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	fmt.Println("CIAO ", request.ServicedInfo.ServiceName)
	fmt.Println("CIAO ", request.ServicedInfo.Uuid)
	fmt.Println("CIAO ", request.ToPolycubed)
	fmt.Println("unsafe.Sizeof(request) ", unsafe.Sizeof(request))
}

/*
	This method is used to unsubscribe from the service
*/
func Unsubscribe() {
	ctx := context.Background()
	request := &ServicedInfo{ServiceName: TheServiced.serviceName , Uuid: TheServiced.uuid}
	reply, err := TheClient.client.Unsubscribe(ctx, request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Reply from Unsubscribe: %s", reply.GetStatus())
}



/*
  ====================================================================================================
                                        STREAM MANAGEMENT
  ====================================================================================================
*/




/*
	This method is used to read the long lived stream waiting for a request 
	(a ToServiced message) from Polycubed to arrive.
	After reading the request and understanding what type it is, 
	it will call the correct method based on the registration made at the beginning

	IL METODO CORRETTO DA CHIAMARE NON È QUI NELLA LIBRERIA MA NEL CONTROLPLANE
	OSSIA NEL MAIN.GO

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
	for {
		request, err := TheClient.stream.Recv()
		if err != nil {
			fmt.Println("Failed to receive message: %v", err)
			continue
		}
		fmt.Println("Message from Polycubed ", request.GetServicedInfo().GetServiceName())
	
	// router logic with method call of the controlplane
		if request.GetRestUserToServiceRequest() != nil {
			fmt.Println("ToServiced_RestUserToServiceRequest ")



			// pino: TODO testing a simple call ==> it works
			var key uint32 = 0
			var value uint8
			r1 := TableGet("h1","action_map",CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
			fmt.Println(r1)
			
			value = 2 // forward
			r2 := TableSet("h1","action_map",CubeInfo_INGRESS,key,unsafe.Sizeof(key),value,unsafe.Sizeof(value))
			fmt.Println(r2)  // should be true
			
			r3 := TableGet("h1","action_map",CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
			fmt.Println(r3)
		}
	
	}
	wg.Done()
}



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
                                        EBPF MAPS MANAGEMENT
  ====================================================================================================
*/



func TableGet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value_size uintptr) interface{} {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
		MapName: map_name,
		ProgramType: program_type,
	}
	fmt.Println("key_size ", key_size)
	fmt.Println("value_size ", value_size)

	fmt.Println("TABLEGET ", key)
	key_bytes := new(bytes.Buffer)
	if ret, ok := key.(interface{}); ok {
		
		err := binary.Write(key_bytes, binary.LittleEndian, ret)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("nooooooooo")
	}
	

/*
	switch v := key.(type) {
		case nil:
			fmt.Println("x is nil", v)           // here v has type interface{}
		case uint32: 
			ret, ok := key.(uint32)
			buf.Write([]byte{ret})
			
			fmt.Println("x is", v)
			fmt.Println("x is", ret)
			fmt.Println("x is", ok)             // here v has type int
		case bool:
			fmt.Println("x is bool", v) 			// here v has type interface{}
		case string:
			fmt.Println("x is string", v)			// here v has type interface{}
		default:
			fmt.Println("type unknown", v)        // here v has type interface{}
	}
	*/
	
	/*
	var real_key []byte
	if ret, ok := key.([]byte); ok {
		fmt.Println("yeeeeee", ret)
	} else {
		fmt.Println("nooooooooo")
	}
	
	*/
	b, ok := key.(interface{})
	fmt.Println("\n -------------------------")
	fmt.Println(ok)
	fmt.Println(b)
	fmt.Println("\n -------------------------")
	
	
	// key_bytes := make([]byte, key_size)
	// pino: TODO, per ora così
	request.RequestType = &ServicedToDataplaneRequest_GetRequest{
			GetRequest: &GetRequest{
				All: false,
				Key: key_bytes.Bytes(),
				KeySize: uint64(key_size),
				ValueSize: uint64(value_size),
			},
	}
	fmt.Println("TABLEGET ", []byte(fmt.Sprint(key)))
	reply, err := TheClient.client.TableGet(ctx, &request)
	if err != nil {
		log.Fatalf("could not use TableGet: %v", err)
	}
	fmt.Println("Reply from TableGet: ", reply.ByteValues[0])
	fmt.Println("unsafe.Sizeof(reply.ByteValues)) ", unsafe.Sizeof(reply.ByteValues))
	fmt.Println("len(reply.ByteValues)) ", len(reply.ByteValues))


	return "ciao"
}




func TableSet(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr, value interface{}, value_size uintptr) bool {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
		MapName: map_name,
		ProgramType: program_type,
	}
	fmt.Println("key_size ", key_size)
	fmt.Println("value_size ", value_size)

	fmt.Println("TABLESET ", key)
	fmt.Println("TABLESET ", value)
	key_bytes := new(bytes.Buffer)
	if ret_key, ok_key := key.(interface{}); ok_key {
		
		err_key := binary.Write(key_bytes, binary.LittleEndian, ret_key)
		if err_key != nil {
			fmt.Println("binary.Write failed:", err_key)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("nooooooooo")
	}
	
	b_key, ok_key := key.(interface{})
	fmt.Println("\n -------------------------")
	fmt.Println(ok_key)
	fmt.Println(b_key)
	fmt.Println("\n -------------------------")
	
	value_bytes := new(bytes.Buffer)
	if ret_value, ok_value := value.(interface{}); ok_value {
		
		err_value := binary.Write(value_bytes, binary.LittleEndian, ret_value)
		if err_value != nil {
			fmt.Println("binary.Write failed:", err_value)
		}
		fmt.Printf("% x", value_bytes.Bytes())
	} else {
		fmt.Println("nooooooooo")
	}

	
	
	// key_bytes := make([]byte, key_size)
	// pino: TODO, per ora così
	request.RequestType = &ServicedToDataplaneRequest_SetRequest{
			SetRequest: &SetRequest{
				Key: key_bytes.Bytes(),
				KeySize: uint64(key_size),
				Value: value_bytes.Bytes(),
				ValueSize: uint64(value_size),
			},
	}
	fmt.Println("TABLESET ", []byte(fmt.Sprint(key)))
	fmt.Println("TABLESET ", []byte(fmt.Sprint(value)))
	reply, err := TheClient.client.TableSet(ctx, &request)
	if err != nil {
		log.Fatalf("could not use TableGet: %v", err)
	}
	fmt.Println("Reply from TableSet: ", reply.GetStatus())

	return reply.GetStatus()
}




func TableRemove(cube_name string, map_name string, program_type CubeInfo_ProgramType, key interface{}, key_size uintptr) bool {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
		MapName: map_name,
		ProgramType: program_type,
	}
	fmt.Println("key_size ", key_size)

	fmt.Println("TABLEREMOVE ", key)
	key_bytes := new(bytes.Buffer)
	if ret_key, ok_key := key.(interface{}); ok_key {
		
		err_key := binary.Write(key_bytes, binary.LittleEndian, ret_key)
		if err_key != nil {
			fmt.Println("binary.Write failed:", err_key)
		}
		fmt.Printf("% x", key_bytes.Bytes())
	} else {
		fmt.Println("nooooooooo")
	}
	
	b_key, ok_key := key.(interface{})
	fmt.Println("\n -------------------------")
	fmt.Println(ok_key)
	fmt.Println(b_key)
	fmt.Println("\n -------------------------")
	
	
	// key_bytes := make([]byte, key_size)
	// pino: TODO, per ora così
	request.RequestType = &ServicedToDataplaneRequest_RemoveRequest{
			RemoveRequest: &RemoveRequest{
				All: false,
				Key: key_bytes.Bytes(),
				KeySize: uint64(key_size),
			},
	}
	fmt.Println("TABLEREMOVE", []byte(fmt.Sprint(key)))
	reply, err := TheClient.client.TableRemove(ctx, &request)
	if err != nil {
		log.Fatalf("could not use TableRemove: %v", err)
	}
	fmt.Println("Reply from TableRemove: ", reply.GetStatus())

	return reply.GetStatus()
}





/*
  ====================================================================================================
                                        CUBE MANAGEMENT
  ====================================================================================================
*/



/*
	questo metodo verrà chiamato dal main.go oppure da readthestream che chiama
	il metodo CreateCube del controlplane che a sua volta chiama questo
*/
func CreateCube(cube_name string, ingress_code string, egress_code, conf string) {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// request := &CubeManagement{CubeName: cube_name, ServiceName: &TheServiced.serviceName}
	request := CubeManagement{}
	request.CubeName = cube_name
	request.ServiceName = &TheServiced.serviceName
	request.Uuid = &TheServiced.uuid
	request.Conf = &conf
	request.IngressCode = &ingress_code
	request.EgressCode = &egress_code
	fmt.Println("CubeManagement ", request.CubeName, " , ", request.ServiceName, "uuid ", request.Uuid)
	reply, err := TheClient.client.CreateCube(ctx, &request)
	if err != nil {
		log.Fatalf("could not create cube: %v", err)
	}
	log.Printf("Reply from CreateCube: %s", reply.GetStatus())
}




// pino: TODO
/*
	in the grpc-go examples, in unary requests, a context is done with withtimeout, 
	with 1 or 10 seconds. My question is, is it wrong to create a context that always 
	lives with Background? Once you exit the Unary method, isn't it automatically deleted?
*/
/*
	This method is used to request the destruction of a cube from Polycubed

*/
func DestroyCube(cube_name string) {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// request := &CubeManagement{CubeName: cube_name, ServiceName: &TheServiced.serviceName}
	request := CubeManagement{}
	request.CubeName = cube_name
	request.ServiceName = &TheServiced.serviceName
	request.Uuid = &TheServiced.uuid
	fmt.Println("CubeManagement ", request.CubeName, " , ", request.ServiceName, "uuid ", request.Uuid)
	reply, err := TheClient.client.DestroyCube(ctx, &request)
	if err != nil {
		log.Fatalf("could not destroy cube: %v", err)
	}
	log.Printf("Reply from DestroyCube: %s", reply.GetStatus())
}





/*
  ====================================================================================================
                                        PORT MANAGEMENT
  * Bool SetPort(const std::string cube_name, const std::string port_name);
  * Bool DelPort(const std::string cube_name, const std::string port_name);
  * Bool SetPeer(const std::string cube_name, const std::string port_name, const
  std::string peer_name);
  ====================================================================================================
*/


func SetPort(cube_name string, port_name string) bool {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
	}


	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
		},
	}
	fmt.Println("SETPORT", port_name)
	reply, err := TheClient.client.SetPort(ctx, &request)
	if err != nil {
		log.Fatalf("could not use SetPort: %v", err)
	}
	fmt.Println("Reply from SetPort: ", reply.GetStatus())

	return reply.GetStatus()
}





func DelPort(cube_name string, port_name string) bool {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
	}


	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
		},
	}
	fmt.Println("DELPORT", port_name)
	reply, err := TheClient.client.DelPort(ctx, &request)
	if err != nil {
		log.Fatalf("could not use DelPort: %v", err)
	}
	fmt.Println("Reply from DelPort: ", reply.GetStatus())

	return reply.GetStatus()
}



func SetPeer(cube_name string, port_name string, peer_name string) bool {
	ctx := context.Background()	
	request := ServicedToDataplaneRequest{}
	request.ServicedInfo = &ServicedInfo {
			ServiceName: TheServiced.serviceName,
			Uuid: TheServiced.uuid,
	}
	request.CubeInfo = &CubeInfo {
		CubeName: cube_name,
	}


	request.RequestType = &ServicedToDataplaneRequest_Port{
		Port: &Port{
			Name: port_name,
			Peer: &peer_name,
		},
	}
	fmt.Println("SETPEER", port_name)
	reply, err := TheClient.client.SetPeer(ctx, &request)
	if err != nil {
		log.Fatalf("could not use SetPeer: %v", err)
	}
	fmt.Println("Reply from SetPeer: ", reply.GetStatus())

	return reply.GetStatus()
}

















func Pino() {
	prova := chi.NewRouter()
	prova.Get("polycube/v1/{serviceName}/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hellooooooooooooooooooooo")
	})

	url := "polycube/v1/router/*"
	r, _ := http.NewRequest("POST", url, strings.NewReader(""))
	fmt.Println("Ciao ", r.URL)
	fmt.Println("Suca ", chi.URLParam(r,"serviceName"))
}














































func sendMessage(stream Polycube_SubscribeClient, msg *ToPolycubed) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(msg)
}

func recvMessage(wantErrCode codes.Code) {
		res, err := TheClient.stream.Recv()
		if status.Code(err) != wantErrCode {
			log.Fatalf("TheClient.stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, wantErrCode)
		}
		if err != nil {
			fmt.Printf("TheClient.stream.Recv() returned expected error %v\n", err)
			return
		}
		fmt.Printf("received message %q\n", res.GetServicedInfo().GetServiceName())
}


func UniqueFunction() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewPolycubeClient(conn)

	// non scritto da pino
	// Initiate the stream with a context that supports cancellation.

	// Background returns a non-nil, empty Context. It is never canceled, has no values,
	//  and has no deadline. It is typically used by the main function, initialization, 
	// and tests, and as the top-level Context for incoming requests. 
	ctx := context.Background()
	
	stream, err := c.Subscribe(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}
	suca := "cubo"
	request := &ToPolycubed {ServicedInfo: &ServicedInfo{ServiceName: "NOME", Uuid: 1111111, CubeName: &suca }, ToPolycubed: &ToPolycubed_Subscribe {Subscribe: true}}
	// Send some test messages.
	if err := sendMessage(stream, request); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	for {

	}
	
}
