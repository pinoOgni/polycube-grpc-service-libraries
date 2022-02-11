package main

import (
	pb "github.com/pinoOgni/polycube-grpc-service-libraries/polycube-grpc-go" // pb stands for polycube
	"time"
	"fmt"
	"io/ioutil"
	"sync"
	// "unsafe"
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
	



	pb.Subscribe("helloworldgo")
	ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
	egress := GetDataplaneCode("Helloworld2_dp_egress.c")
	time.Sleep(4*time.Second)
	// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
	conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
	pb.CreateCube("h1",ingress,egress,conf)
	time.Sleep(4*time.Second)


	// pino: TODO se lo metto prima e dopo faccio qualcosa mi da errore e crasha
	var wg sync.WaitGroup
	wg.Add(1)
	go pb.ReadTheStreamGoRoutine(&wg)
	wg.Wait()

	
	fmt.Println("the end")
	// time.Sleep(100*time.Second)
}



func GetDataplaneCode(filename string) string {
	b, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Print(err)
    }
	return string(b)
}



