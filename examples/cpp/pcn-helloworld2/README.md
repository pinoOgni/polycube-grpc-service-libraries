Questo è il primo esempio di grpc C++ service. All'interno di questa directory ci sono vari tipi di file e directory.

Directory:
* cmake: contiene  i cmake utili per la compilazione del servizio
* datamodel: contiene lo yang del servizio (PER ORA NON VIENE USATO)
* src: contiene tutto il codice del servizio helloworld 1.0 e altri file che compongono il servizio grpc. Di tutti i file contenuti qui dentro vengono usati solo:
    * helloworld2_client.cc  rappresenta il CONTROLPLANE. Qui c'è il main del servizio e devono essere fatte alcune cose per far funzionare tutto:
        * PolycubeClient polycubeclientgrpc; in questo modo possiamo usare i metodi della libreria polycube-grpc (c++)
        * usare i metodi RegisterHandlerHTTPVERB per associare i metodi del controlplane ad un url
        * lanciare in un thread il metodo ReadTheStream della libreria polycube-grpc
        * fare un join del thread appena lanciato
    * helloworld2_client.h rappresenta il CONTROLPLANE. 
        * Qui non c'è solo il prototipo del metodo ma anche l'implementazione. Ogni metodo che viene implementato qui e che necessita di contattare Polycube per il proprio lavoro, dovrà effettuare una chiamata ad un metodo opportuno della libreria polycube-grpc. 
        * Per esempio, il metodo setAction dell'Helloworld serve a settare una action nel DATAPLANE ossia a settare un valore in una mappa eBPF, ciò si riduce a chiamare il metodo TableSet con opportuni parametri.
        * Da notare che c'è anche il metodo create_cube che serve a creare un cubo
        * IMPORTANTE: qui, lo sviluppatore potrà decidere la logica con cui gestire il rapporto tra servizio e cubo, quindi per esempio gestire i cubi con degli ID e avere un vettore di cubi o simile
        * IMPORTANTE: qui, lo sviluppatore deve anche gestire la logica delle porte, ovviamente se il servizio ha in sè il concetto di porta
        * LESS IMPORTANT: ora come ora in quesot file, lo sviluppatore dovrà mettere anche il metodo helpocommands che verrà registrato poi su l'url opportuno e che ritorna una descrizione su come usare il servizio da polycubectl
    * Helloworld2_dp_egress.h, Helloworld2_dp_egress.c Helloworld2_dp_ingress.h e Helloworld2_dp_ingress.c che corrispondono al dataplane
    * Helloworld2-lib.cpp, Helloworld2.cpp, Helloworld2.h, Ports.cpp e Ports.h ORA come ora non servono, non so se verrano usati in futuro
* first_client.cc, first_polycube_client.cc, router_client.cc e second_client.cc NON fanno pare del servizio ma sono file di test


OSSERVAZIONE: per far funzionare tutto, si dovrà includere nei file del controlplane la libreria polycube-grpc






# Examples gRPC remote services in C++

## 






## Osservations

Within each example, for logic reasons, or because it is still being implemented, there are some directories or files that are not necessary to develop a simple service. In particular, the directories
* cmake: contains some files for loading the files as variables and another one used for the compiler of gRPC/ProtoBuf
* datamodel: contains the datamodel of the Polycube helloworld service