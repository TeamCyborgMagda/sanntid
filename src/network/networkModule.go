//package main
package main

import (
    "fmt"
    "driver" 
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
	l_adress, err := net.ResolveUDPAddr("udp", ":20008")
	listen, err := net.ListenUDP("udp",l_adress)

	w_adress, err := net.ResolveUDPAddr("udp", "129.241.187.255"+":20008")
	conn, err := net.DialUDP("udp",nil,w_adress)
	return ip, port, listen, conn,  err

// 	returns 3 variables, brodcast_connection, ip, port (?)
}

//			initializes/reset state, resets/initialize number of slaves, and initialize master-ip, initialize connections array		
func StateInit(conn *net.UDPConn)(string, string, int){
	buffer := make([]byte,128)
	for{ 
		conn.SetDeadline(time.Now().Add(6*time.Second))
		_,err := conn.Read(buffer)
		if err != nil{
//	step) hvis man ikke hører en master, returner "master" + "nil"
			fmt.Println("finner ingen master", err )			
			return "master", "", 1
		}
//	step) hvis man hørte en master, returner "slave" + ip adresse  
		read_ip:=string(buffer)
      master_ip:= strings.Split(read_ip, "\x00")[0]
		return "slave", master_ip, -1
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

/*
func CheckConnection(state string, master_ip string, broadcast_listener *net.UDPConn, ip string)(int){
	buffer := make([]byte, 128)
	var read_ip string	
	if state == "slave"{
		broadcast_listener.SetDeadline(time.Now().Add(400*time.Millisecond))		
		_,err := broadcast_listener.Read(buffer)
		if err != nil{
			return 0
		}		
		read_ip =string(buffer)
        master_ip_r:= strings.Split(read_ip, "\x00")[0]
		if master_ip_r != master_ip{
			return 0
		}
	}
	if state == "master"{
		broadcast_listener.SetDeadline(time.Now().Add(400*time.Millisecond))		
		_,err := broadcast_listener.Read(buffer)
		read_ip := string(buffer)				
		if err != nil{
			fmt.Println("skjekk: ", err)
		}		
		read_ip =string(buffer)
        master_ip_r:= strings.Split(read_ip, "\x00")[0]
		if master_ip_r != ip{
			fmt.Println("nummer to trigger også!: ", master_ip_r)
		}
	}
	return 1
 
}
*/

func main(){
	var order_array [10]driver.Data
	

//udefinert state loop,	Lurt å ha heis funksjon sammen med ip?, bare ha ting i init som man er SIKKER på at kun skal kjøres EN gang?
	ip, master_port, broadcast_listener, broadcast_writer, err := Init()
	fmt.Println("Init gikk bra")
	if err != nil{return}
	
	/*
	go func(){
		for{
			select{
			case data := <- remove_order:
				remove_order_copy = data
			case data := <- order_queue:
				order_queue_copy = data
			case data := <- cost:
				cost_copy = data
			}
			time.Sleep(1*time.Millisecond)
		}
	}()
	*/
		
	for{
		
		state, master_adress, elevator_nr := StateInit(broadcast_listener)	//bestemme funksjon
		fmt.Println("State init gikk bra")
		
		connections := make([]net.Conn,10)
		nr_of_slaves := 0
 		slave_listener, err := SlaveListener(master_port)		//			master loop, n = number of slaves 
 							
 		for (state=="master"){
//			if CheckConnection(state, master_adress,broadcast_listener,ip) != 1{
//				break
//			}
			fmt.Println("Heisen er en master")
			broadcast_writer.Write([]byte(ip+"\x00"))
			fmt.Println(ip)
			elevator_number <- elevator_nr   // Må opprette channel
			
			slave_listener.SetDeadline(time.Now().Add(2*time.Second))
			connections[nr_of_slaves],err = slave_listener.Accept()
			if err != nil{
				fmt.Println("Finner ingen slaver: ", err)
			}else{
				nr_of_slaves = nr_of_slaves + 1
				connections[nr_of_slaves].Write([]byte(string(nr_of_slaves+1) + "\x00"))  //need int to string conversion
				fmt.Println("inkrementerer variabel")
   		}
   		
			i := 0
			buffer := make([]byte, 128)
			for i < nr_of_slaves {
				connections[i].Write([]byte("send\x00"))
				connections[i].Read(buffer)
				cost_array[i+1] = buffer[0]
				order_array[i+1] = buffer[1]				//Need conversions from byte to driver.Data type. 
				remove_array[i+1] = buffer[2]
				i += 1
			}
			j := 0
			for i <= nr_of_slaves{
				j = 0
				for j<8{
					if(order_array[i].Array[j] == 1){
						order_list[j] = 1
					}
					j += 1
				}
				j = 0
				for j <8{
					if(remove_array[i].Array[j] == 1){
						order_list[j] = 0
					}
					j += 1
				}
				i +=1
			}
			j = 0
			for j<8{
				i = 0
				lowest_costs[j] = 0
				for i<nr_of_slaves{
					if(cost_array[i+1].Array[j] < cost_array[lowest_cost[j]].Array[j]){
						lowest_cost[j] = i
					}
				
				}
			}
			j = 0
			for j<8{ 	
				order_list[j] = order_list[j]*(lowest_cost[j] + 1) // 0 - (nr_of_slaves +1) based on the elevator with lowest cost. 
			}
			broadcast_writer.Write(string(order_list) + "\x00")
			order_list <- order_list_yo							//needs to be made channel name/or changed variable name. 
			time.Sleep(1*time.Millisecond)	
		}
// 	"Videre: prossesering av data og sending av informasjon til slaver




 		if state == "slave" {
			connections[0],err = ConnectMaster(master_adress + master_port)
			if err != nil{
				fmt.Println("Cannot connect to master: ", err);
				continue
			}
			buffer := make([]byte,128)
			connections[0].Read(buffer)
			elevator_nr = buffer		//need string to int conversions
			elevator_number <- elevator_nr  // må opprette kanal 		
		}			
//									Slave loop
		for state == "slave"{
//			if !CheckConnection(state){
//				break						
//			}

			fmt.Println("Heisen er en slave")
			buffer = make([]byte, 128)
			
			connections[0].Read(buffer)
			if string(buffer) == "send"{
				connections[0].Write([]byte(string(order_queue.Array)+ string(cost.Array) + string(remove_orders.Array) + "\x00" )
				broadcast_listener.Read(buffer)
				order_list_yo = buffer //converted to driver.Data yay! ++ standariser navn yo
				order_list <- order_list_yo // navn og channel ikke opprettet. 
			}
			
			
			time.Sleep(100*time.Millisecond)
//			HandleMastersOrders()
		}
	}
}




/*
master løække:
 for N_slaves{
	gi klarsignal
	les melding
 	hvis(error timeout) - inkrementer feil
 		hvis feil = maks tillat. Fjern fra lista. 
 	oppdater ordre
 	broadkast ordre
 
 

 
 vente på send
 Hvis man venter for lenge i strekk, avbryt og se etter ny master. 
 sende oppdatterrtte lister
 lese broadcast for ny ordre_liste
 
 	
 	



*/ 
 


