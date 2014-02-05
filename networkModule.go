package main

import (
    	 "fmt" 
	 "net"
	"time"
)
func Init()(string, string, *UDPConn) {
//	step 1) "finne pc'ens ip'
// 	step 2) "finne master port"				
//  step 3) "opprett UDP(?) broadcast listner/connection"
// 	returns 3 variables, brodcast_listner/connection, ip, port (?)
}

func StateInit(connection *UDPConn)(int, string, string){
//	step 1) "høre på broadcast i en viss tid"
//	step 2) if setning
//	step 3) hvis man hørte en master, returner "slave" + ip adresse  
//	step 3.1) hvis man ikke hører en master, returner "master" + "nil"
}

//		Listens for slaves and returns the slave connection
func MakeSlave(port string)(*TCPConn, error){
	listen_adress, err := net.ResolveTCPAddr("tcp", port)
	slave_listener, err := net.ListenTCP("tcp", listen_adress)
	conn, err := slave_listener.Accept()
	return conn, err
}

// 		attempts to connect to the given adress
func ConnectMaster(adress string)(*TCPConn, error){
	master_adr, err :=  net.ResolveTCPAddr("tcp", adress)
	conn, err := net.DialTCP("tcp", nil, master_adr)
	return conn, err 
}

func network(){
//udefinert state loop,	Lurt å ha heis funksjon sammen med ip?, bare ha ting i init som man er SIKKER på at kun skal kjøres EN gang?
	ip, master_port, broadcast_conn := Init()
// for lopp forever
	nr_of_slaves, state, master_adress := StateInit(broadcast_conn)	//"initialiser verdier, bestemme funksjon" 
// 							master loop, n = number of slaves
// 						(slave nr. 1 is the successor(?) neccesary?)
// 	for state = "master"
//		if (CheckConnection()==False){
//			break
//		}
//		broadcast_conn.Write("i'm the master, find me at ip")
//		broadcast_conn.Read(broadcast_data)
//		if (slave blir identifisert){
//			connections[nr_of_slaves],err = MakeSlave(master_port)
//			if err != nil
//				continue
//			nr_of_slaves++
//      }
//		for connection in connections{
//			connection.Write("HAI there")
//			connection.Read(data)
//			fmt.printf(string(data))
//		}	
// 		"Videre: prossesering av data og sending av informasjon til slaver


// 	if state = "slave" {
//		connection,err := ConnectMaster(master_adress + master_port)
//		if err != nil{
//			fmt.Println("Cannot connect to master: ", err);
// 			continue
//		}		
//	}			
//									Slave loop
//	for state = "slave"
//		if !checkConnections(){
//			break						
//		}
//		ReadElevatorStatus()
//		connection.Write(data)
//		connection.Read(data)
//		HandleMastersOrders()
}

func such_network() {
	much_TCP()
	wow()

}
