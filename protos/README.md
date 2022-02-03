# PROTO

In questo file si cerca di fare una sintesi dei metodi e dei messaggi che sono implementati nel file proto

## Methods

*   rpc Unsubscribe(ServicedInfo) returns (Bool) {}
    * Viene usato per disiscrivere un servizio da Polycube
*   rpc Subscribe(stream ToPolycubed) returns (stream ToServiced) {} 
    * Viene usato per iscrivere un servizio a Polycube
    * Crea anche il long-lived-stream tra Polycube e il servizio
*   rpc DestroyCube(CubeManagement) returns (Bool) {}
*   rpc CreateCube(CubeManagement) returns (Bool) {}
*   rpc TableSet(ServicedToDataplaneRequest) returns (Bool) {}
*   rpc TableGetAll(ServicedToDataplaneRequest) returns (MapValue) {}
*   rpc TableGet(ServicedToDataplaneRequest) returns (MapValue) {}
*   rpc TableRemove(ServicedToDataplaneRequest) returns (Bool) {}
*   rpc SetPort(ServicedToDataplaneRequest) returns (Bool)  {}


## OLD METHODS

*   rpc AddProgram(Dataplane) returns (Index) {}
*   rpc DelProgram(Dataplane) returns (Empty) {}

## TO IMPLEMENT METHODS

*   rpc Reload(Dataplane) returns (Empty) {}
*   rpc TableRemoveAll(ServicedToDataplaneRequest) returns (Bool) {}

## Messages

* Index: è un indice che doveva essere ritornato dalla add_program ma che ora non serve a niente
* Empty: rappresenta un messaggio vuoto, per i metodi void
* Dataplane: doveva essere usato dalla add_program, del_program e reload e contiene varie info sul dataplane. NON viene usato
* Bool: rappresenta un valore boolean
* ServicedInfo: è un insieme di varie info sul servizio come nome, uuid, nome del cubo e nome della porta
* ToServiced: è un messaggio wrapper con ServicedInfo e poi un oneof tra DataplaneToservicedRequest e RestUserToServicedRequest. Come suggerito dal nome prende tutte le info dirette verso il Serviced
* ToPolycubed:  è un messaggio wrapper con ServicedInfo e poi un oneof tra ServicedToDataplaneReply e ServicedToRestUserReply. Come suggerito dal nome prende tutte le info dirette verso Polycubed
* RestUserToServicedRequest: rappresenta una richiesta rest dall'utente verso il servizio. Contiene il verbo http, l'url e il corpo della richiesta http
* DataplaneToServicedRequest
* ServicedToDataplaneReply
* ServicedToRestUserReply
* ServicedToDataplaneRequest
* CubeInfo
* CubeManagement
* Port
* MapValue: rappresenta un valore di una mappa e viene usato nei casi di risposta ad una richiesta user to dataplane oppure service to dataplane; in entrambi i casi va da Polycube al Servizio
* SetRequest: serve per settare un valore in una mappa
* GetRequest: serve per ottenere un valore da una mappa
* RemoveRequest: serve nella rimozione di un valore da una mappa



