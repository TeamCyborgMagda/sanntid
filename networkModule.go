package main

import (
    "fmt" 
	"net"
	"time"
)

//		initializes the pc's ip adress, standard master port and the UDP broadcast connection 
func Init()(string, string, *UDPConn, error) {
	AllAddr, err := net.InterfaceAddrs()
	if err != nil{
		fmt.Println("couldn't find ip")
	}
	ip := strings.Split(AllAddr[1].String(),"/")[0]
	port := ":33546"
	adress, err := net.ResolveUDPAddr("udp", ":20002")
	conn, err := net.ListenUDP("udp4",adress)
	return ip, port, conn, err

// 	returns 3 variables, brodcast_connection, ip, port (?)
}

//			initializes/reset state, resets/initialize number of slaves, and initialize master-ip, initialize connections array		
func StateInit(conn *UDPConn)(string, string){
	buffer := make([]byte,128)
	for{ 
		conn.SetDeadline(time.Now().Add(4*time.Second))
		_,err := conn.Read(buffer)
		if err != nil{
//	step) hvis man ikke hører en master, returner "master" + "nil"
			return "master", ""
		}
//	step) hvis man hørte en master, returner "slave" + ip adresse  
		read_ip:=string(buffer)
        master_ip,_:= strings.Split(read_ip, "\x00")[0]
		return "slave", master_ip
	}
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





func NetworkModule(){
//udefinert state loop,	Lurt å ha heis funksjon sammen med ip?, bare ha ting i init som man er SIKKER på at kun skal kjøres EN gang?
	ip, master_port, broadcast_conn, err := Init()
	for{
		state, master_adress := StateInit(broadcast_conn)	//bestemme funksjon
		connections := make([]*net.TCPConn,10)
		nr_of_slaves := 0
// 							master loop, n = number of slaves 					
 		for (state=="master"){
//			if !CheckConnection(state, connections){
//				break
//			}
			broadcast_conn.Write(ip)
			connections[nr_of_slaves],err = MakeSlave(master_port)
			if err != nil{
				nr_of_slaves = nr_of_slaves + 1
   			}
			i := 0
			buffer := make([]byte, 128)
			for i < nr_of_slaves {
				connections[i].Write("HAI there")
				connection[i].Read(buffer)
				fmt.printf(string(buffer))
				i += 1
			}	
		}
// 		"Videre: prossesering av data og sending av informasjon til slaver


 		if state == "slave" {
			connections[nr_of_slaves],err := ConnectMaster(master_adress + master_port)
			if err != nil{
				fmt.Println("Cannot connect to master: ", err);
 			
			}		
		}			
//									Slave loop
		for state == "slave"{
//			if !CheckConnection(state){
//				break						
//			}
//			ReadElevatorStatus()
			buffer := make([]byte, 128)
			connections[nr_of_slaves].Write("bai there")
			connections[nr_of_slaves].Read(buffer)
			fmt.printf(string(buffer))
//			HandleMastersOrders()
		}
	}
}

func main(){
	NetworkModule()
}
