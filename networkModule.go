package main

import (
    "fmt" 
	"net"
	"time"
	"strings"
)

//		initializes the pc's ip adress, standard master port and the UDP broadcast connection 
func Init()(string, string, *net.UDPConn,*net.UDPConn, error) {
	AllAddr, err := net.InterfaceAddrs()
	if err != nil{
		fmt.Println("couldn't find ip")
	}
	ip := strings.Split(AllAddr[1].String(),"/")[0]
	port := ":33546"
	l_adress, err := net.ResolveUDPAddr("udp", ":20002")
	listen, err := net.ListenUDP("udp",l_adress)

	w_adress, err := net.ResolveUDPAddr("udp", "192.168.0.255"+":20002")
	conn, err := net.DialUDP("udp",nil,w_adress)
	return ip, port, listen, conn,  err

// 	returns 3 variables, brodcast_connection, ip, port (?)
}

//			initializes/reset state, resets/initialize number of slaves, and initialize master-ip, initialize connections array		
func StateInit(conn *net.UDPConn)(string, string){
	buffer := make([]byte,128)
	for{ 
		conn.SetDeadline(time.Now().Add(6*time.Second))
		_,err := conn.Read(buffer)
		if err != nil{
//	step) hvis man ikke hører en master, returner "master" + "nil"
			fmt.Println(err)			
			return "master", ""
		}
//	step) hvis man hørte en master, returner "slave" + ip adresse  
		read_ip:=string(buffer)
        master_ip:= strings.Split(read_ip, "\x00")[0]
		return "slave", master_ip
	}
}

//		Listens for slaves and returns the slave connection
func SlaveListener(port string)(*net.TCPListener, error){
	listen_adress, err := net.ResolveTCPAddr("tcp", port)
	slave_listener, err := net.ListenTCP("tcp", listen_adress)
	return slave_listener, err
}

// 		attempts to connect to the given adress
func ConnectMaster(adress string)(*net.TCPConn, error){
	master_adr, err :=  net.ResolveTCPAddr("tcp", adress)
	conn, err := net.DialTCP("tcp", nil, master_adr)
	return conn, err 
}





func NetworkModule(){
//udefinert state loop,	Lurt å ha heis funksjon sammen med ip?, bare ha ting i init som man er SIKKER på at kun skal kjøres EN gang?
	ip, master_port, broadcast_listener, broadcast_writer, err := Init()
	fmt.Println("Init gikk bra")
	if err != nil{return}
	for{
		state, master_adress := StateInit(broadcast_listener)	//bestemme funksjon
		fmt.Println("State init gikk bra")
		connections := make([]net.Conn,10)
		nr_of_slaves := 0
 		slave_listener, err := SlaveListener(master_port)		//			master loop, n = number of slaves 					
 		for (state=="master"){
//			if !CheckConnection(state, connections){
//				break
//			}
			fmt.Println("Heisen er en master")
			broadcast_writer.Write([]byte(ip+"\x00"))
			
			slave_listener.SetDeadline(time.Now().Add(2*time.Second))
			connections[nr_of_slaves],err = slave_listener.Accept()
			if err != nil{
				fmt.Println("Finner ingen slaver: ", err)
			}			
			if err == nil{
				nr_of_slaves = nr_of_slaves + 1
				fmt.Println("inkrementerer variabel")
   			}
			i := 0
			buffer := make([]byte, 128)
			for i < nr_of_slaves {
				connections[i].Write([]byte("HAI there"))
				connections[i].Read(buffer)
				fmt.Println(string(buffer))
				i += 1
			}
			time.Sleep(100*time.Millisecond)	
		}
// 		"Videre: prossesering av data og sending av informasjon til slaver


 		if state == "slave" {
			connections[0],err = ConnectMaster(master_adress + master_port)
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
			fmt.Println("Heisen er en slave")
			buffer := make([]byte, 128)
			connections[nr_of_slaves].Write([]byte("BAI there" + "\x00"))
			connections[nr_of_slaves].Read(buffer)
			fmt.Println(string(buffer))
			time.Sleep(100*time.Millisecond)	
//			HandleMastersOrders()
		}
	}
}

func main(){
	NetworkModule()
}
