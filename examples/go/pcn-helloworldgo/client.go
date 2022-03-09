package main

import (
	pb "polycube-grpc-go/commons"
	context "context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
	"unsafe"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
)


var kacp = keepalive.ClientParameters{
	Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             5 * time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}



const (
	address     = "localhost:9001"
	defaultName = "world"
)


func main() {
	Cazzo()

}


func ProvaStream() {
	fmt.Println("1")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	fmt.Println("2")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		fmt.Println("connected")
	}
	defer conn.Close()
	client := pb.NewPolycubeClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

 
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("%v.Subscribe(_) = _, %v", client, err)
	}
	suca := "suca"
	request := &pb.ToPolycubed {ServicedInfo: &pb.ServicedInfo{ServiceName: "NOME", Uuid: 1111111, CubeName: &suca }, ToPolycubed: &pb.ToPolycubed_Subscribe {Subscribe: true}}
	fmt.Println("unsafe.Sizeof() ", unsafe.Sizeof(request))

	if err := stream.Send(request); err != nil {
		log.Fatalf("Failed to send a toPolycubed message: %v", err)
	}

	cancel()
	
}


func Cazzo() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPolycubeClient(conn)

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
	request := &pb.ToPolycubed {ServicedInfo: &pb.ServicedInfo{ServiceName: "NOME", Uuid: 1111111, CubeName: &suca }, ToPolycubed: &pb.ToPolycubed_Subscribe {Subscribe: true}}
	// Send some test messages.
	if err := sendMessage(stream, request); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	// Ensure the RPC is working.
	recvMessage(stream, codes.OK)

	// fmt.Println("cancelling context")
	

	// ora come ora funziona però da come problema de = DeadlineExceeded desc = context deadline exceeded; 
	// in pratica devo risettare i parametri nel server e settarli anche nel client
	// cancel()
	// il contesto non deve essere cancellato, devo seguire il programma di long-lived

	recvMessage(stream, codes.OK)
}

func sendMessage(stream pb.Polycube_SubscribeClient, msg *pb.ToPolycubed) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(msg)
}

func recvMessage(stream pb.Polycube_SubscribeClient, wantErrCode codes.Code) {
		res, err := stream.Recv()
		if status.Code(err) != wantErrCode {
			log.Fatalf("stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, wantErrCode)
		}
		if err != nil {
			fmt.Printf("stream.Recv() returned expected error %v\n", err)
			return
		}
		fmt.Printf("received message %q\n", res.GetServicedInfo().GetServiceName())
}


func ProvaUnar() {
	fmt.Println("1")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	fmt.Println("2")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		fmt.Println("connected")
	}
	defer conn.Close()
	client := pb.NewPolycubeClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	request := &pb.ToPolycubed {ServicedInfo: &pb.ServicedInfo{ServiceName: "NOME", Uuid: 1111111}, ToPolycubed: &pb.ToPolycubed_Subscribe {Subscribe: true}}
	fmt.Println("unsafe.Sizeof() ", unsafe.Sizeof(request))
	r, err := client.SimpleUnaryMethod(ctx,request)
	if err != nil {
		log.Fatalf("%v.Subscribe(_) = _, %v", client, err)
	}
	fmt.Println("Risposta ", r.GetServicedInfo().GetServiceName())

	
}