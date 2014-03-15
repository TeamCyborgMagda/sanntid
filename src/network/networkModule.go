package main

import (
    "fmt"
    "driver" 
	"net"
	"time"
	"strings"
	"encoding/json"
	"strconv"
)

type pablo struct{
	Cost [8]int
	OrderQueue [8]int
	RemoveOrder [8]int
}


func main(){
	// temporary variables
	var(
		order_list chan driver.Data
		order_queue chan driver.Data
		remove_order chan driver.Data
		cost chan driver.Data
		order_queue_copy driver.Data
		remove_order_copy driver.Data
		cost_copy driver.Data
		elevator_number chan int //counts slaves
	)
	
	

	var (
		order_queue_array [10]driver.Data //masters array av alle order_queue fått av slavene
		cost_array [10]driver.Data //masters array av alle cost fått av slavene
		remove_order_array [10]driver.Data //masters array av alle remove_order fått av slavene
		trans pablo // brukes til å sende cost, order_queue og remove_order fra slave til master
		order_list_array [8]int //lista som lages av master og sendes til slavene
		order_list_yo driver.Data // ferdig ordreliste i Dataformat
	)

//udefinert state loop,	Lurt å ha heis funksjon sammen med ip?, bare ha ting i init som man er SIKKER på at kun skal kjøres EN gang?
	ip, master_port, broadcast_listener, broadcast_writer, broadcast_orders, recieve_orders, err := Init()
	fmt.Println("Init gikk bra")
	if err != nil{return}
	
	connection_timeouts := [10]int{0,0,0,0,0,0,0,0,0,0}
	
	go func(){
		for{
			select{
			case data := <- remove_order:
				remove_order_copy = data
			case data := <- order_queue:
				order_queue_copy = data
			case data := <- cost:
				cost_copy = data
			case data := <- elevator_number:
				connection_timeouts[4] = data	
			}
			time.Sleep(1*time.Millisecond)
			
		}
	}()
	
	
		
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
			broadcast_writer.Write([]byte(ip+"\x00")) // her sender master ipen sin over broadcast
			fmt.Println(ip)
			//elevator_number <- elevator_nr   // Må opprette channel
			fmt.Println("heisen har nr: ", elevator_nr, "og har: ", nr_of_slaves, " slaver")
			
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
				connections[i].SetDeadline(time.Now().Add(100*time.Millisecond))
				_, err := connections[i].Write([]byte("send\x00"))
				if err != nil{
					connection_timeouts[i] += 1
					if (connection_timeouts[i] > 3){
						// sortere ordre liste
						fmt.Println("Slave number: ", i, "is inactive, teminating connection") 
						nr_of_slaves-= 1
					}
					continue
				}else{
					connections[i].SetDeadline(time.Now().Add(100*time.Millisecond))
					connections[i].Read(buffer)
				}
				
				err = json.Unmarshal(buffer, &trans)
				if (err != nil){
					fmt.Println("klarte ikke pakke opp cost, order queue og remover orders: ", err)
				}else{
					cost_array[i+1].Array = trans.Cost
					order_queue_array[i+1].Array = trans.OrderQueue
					remove_order_array[i+1].Array = trans.RemoveOrder
				}
				
				i += 1
			}
			// assigning orders
			i = 0
			j := 0
			for i <= nr_of_slaves{
				j = 0
				for j<8{
					if(order_queue_array[i].Array[j] == 1){
						order_list_array[j] = 1
					}
					j += 1
				}
				j = 0
				for j <8{
					if(remove_order_array[i].Array[j] == 1){
						order_list_array[j] = 0
					}
					j += 1
				}
				i += 1
			}
			j = 0
			var lowest_cost [8]int
			for j<8{
				i = 0
				lowest_cost[j] = 1
				for i<nr_of_slaves{
					if(cost_array[i].Array[j] < cost_array[lowest_cost[j]].Array[j]){
						lowest_cost[j] = i
					}
					i += 1
				}
				j += 1
			}
			j = 0
			for j<8{ 	
				order_list_array[j] = order_list_array[j]*(lowest_cost[j] + 1) // 0 - (nr_of_slaves +1) based on the elevator with lowest cost. 
				j += 1
			}
			order_list_yo.Array = order_list_array
			buffer, err = json.Marshal(order_list_yo)
			if (err!=nil){
				fmt.Println("error converting order_list to type driver.Data: ", err)
			}else{
				broadcast_orders.Write(buffer) // broadcast sender order list??
			}
			fmt.Println("Hopper over for siden nr_of_slaves = 0")
			//order_list <- order_list_yo
			time.Sleep(1*time.Millisecond)	
		}
// 		"Videre: prossesering av data og sending av informasjon til slaver - CHECK




 		if state == "slave" {
			connections[0],err = ConnectMaster(master_adress + master_port)
			if err != nil{
				fmt.Println("Cannot connect to master: ", err);
				continue
			}
			buffer := make([]byte,128)
			connections[0].Read(buffer)
			
			elevator_nr_str := string(buffer)		//
			elevator_nr,_ = strconv.Atoi(elevator_nr_str)
			//elevator_number <- elevator_nr  //	
		}			
//		Slave loop
		for state == "slave"{
			
			
			fmt.Println("Heisen er en slave")
			buffer := make([]byte, 128)
			
			connections[0].SetDeadline(time.Now().Add(3*time.Second))
			_, err := connections[0].Read(buffer)
			if string(buffer) == "send"{
				

				trans.Cost = cost_copy.Array
				trans.OrderQueue = order_queue_copy.Array
				trans.RemoveOrder = remove_order_copy.Array
				
				buffer,err = json.Marshal(trans)
				if err != nil{
					fmt.Println("klarte ikke pakke ned cost, order queue og remove orders: ", err)
				}else{
					connections[0].Write(buffer)
				}				
				recieve_orders.Read(buffer) // her skal man lese order list fra broadcast. Hvor/når sendes ip?
				err = json.Unmarshal(buffer, &order_list_yo) // er order_list_yo opprettet før detta?
				if (err != nil){
					fmt.Println("klarte ikke pakke opp order list: ", err)
				}else{
					//order_list <- order_list_yo
				}
				
				//order_list_yo = buffer //converted to driver.Data yay! ++ standariser navn yo
				//order_list <- order_list_yo // navn og channel ikke opprettet. 
			}else if err != nil{
				connection_timeouts[0] += 1
				if connection_timeouts[0] > 1{
					connection_timeouts[0] = 0
					fmt.Println("connection timeout, terminating connection to master")
					break
				}
			}
			
			
			time.Sleep(1*time.Millisecond)
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
 
//		initializes the pc's ip adress, standard master port and the UDP broadcast connection 
func Init()(string, string, *net.UDPConn,*net.UDPConn,*net.UDPConn,*net.UDPConn, error) {
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
	
	l_adress, err = net.ResolveUDPAddr("udp", ":20009")
	recieve_orders, err := net.ListenUDP("udp",l_adress)

	w_adress, err = net.ResolveUDPAddr("udp", "129.241.187.255"+":20009")
	broadcast_orders, err := net.DialUDP("udp",nil,w_adress)
	
	return ip, port, listen, conn, broadcast_orders, recieve_orders, err

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

