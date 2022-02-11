# TESTING EXAMPLES

## SUBSCRIBE AND UNSUBSCRIBE METHODS
	
pb.Subscribe("helloworldgo")
fmt.Println("ciao")
time.Sleep(4*time.Second)
pb.Unsubscribe()


## CREATE CUBE AND DESTROY CUBE

ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
egress := GetDataplaneCode("Helloworld2_dp_egress.c")
// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
pb.CreateCube("h1",ingress,egress,conf)
time.Sleep(4*time.Second)
pb.DestroyCube("h1")
	

## TABLE GET AND TABLE SET

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

value = 2 // forward
r := pb.TableSet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),value,unsafe.Sizeof(value))
fmt.Println(r)  // should be true

r := pb.TableGet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),unsafe.Sizeof(value))
fmt.Println(r)


## TABLE REMOVE

// with the code used before
var key uint32 = 0
r_remove := pb.TableRemove("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key))
fmt.Println(r_remove)



## SET PORT

pb.Subscribe("helloworldgo")
time.Sleep(4*time.Second)
ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
egress := GetDataplaneCode("Helloworld2_dp_egress.c")
// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
pb.CreateCube("h1",ingress,egress,conf)
time.Sleep(4*time.Second)

ret_set_port := pb.SetPort("h1","p1")
fmt.Println(ret_set_port)
// polycubectl helloworldgo show cubes (shows 1 item that is the port added)

## DEL PORT

pb.Subscribe("helloworldgo")
time.Sleep(4*time.Second)
ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
egress := GetDataplaneCode("Helloworld2_dp_egress.c")
// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
pb.CreateCube("h1",ingress,egress,conf)
time.Sleep(4*time.Second)

ret_set_port := pb.SetPort("h1","p1")
fmt.Println(ret_set_port)
// polycubectl helloworldgo show cubes (shows 1 item that is the port added)
ret_del_port := pb.DelPort("h1","p1")
fmt.Println(ret_del_port)
// polycubectl helloworldgo show cubes (shows 0 items)


## SET PEER

pb.Subscribe("helloworldgo")
time.Sleep(4*time.Second)
ingress := GetDataplaneCode("Helloworld2_dp_ingress.c")
egress := GetDataplaneCode("Helloworld2_dp_egress.c")
// pino: TODO, ora come ora è lo sviluppatore che deve configurare il proprio conf perché la libreria è generica
conf := `{"name":"h1","type":"TC","loglevel":"INFO","shadow":false,"span":false,"action":"drop"}`
pb.CreateCube("h1",ingress,egress,conf)
time.Sleep(4*time.Second)

ret_set_port := pb.SetPort("h1","p1")
fmt.Println(ret_set_port)

ret_set_peer := pb.SetPeer("h1","p1","veth1")
fmt.Println(ret_set_peer)
time.Sleep(4*time.Second)

// in this way with a simple ping we can see the debug message from polycubed
var key uint32 = 0
var value uint8  = 1 // slowpath
r := pb.TableSet("h1","action_map",pb.CubeInfo_INGRESS,key,unsafe.Sizeof(key),value,unsafe.Sizeof(value))
fmt.Println(r)  // should be true

fmt.Println("the end")
time.Sleep(100*time.Second)