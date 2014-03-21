package network

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


func Network(order_queue chan driver.Data ,remove_order chan driver.Data ,cost chan driver.Data, elevator_number chan int, order_list chan driver.Data, order_list_lights chan driver.Data){
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
				remove_order_array[0] = data
			case data := <- order_queue:
				i := 0
				for i < 8{
					if data.Array[i] == 1{
						order_queue_array[0].Array[i] = 1
					}
					i += 1
				}
			case data := <- cost:
				cost_temp := data
				cost_array[0] = cost_temp
			}
			time.Sleep(1*time.Millisecond)
		}
	}()
	
	
	for{
		state, master_adress, elevator_nr := StateInit(broadcast_listener)	//bestemme funksjon
		
		connections := make([]net.Conn,10)
		nr_of_slaves := 0
			//			master loop, n = number of slaves 
 							
 		for (state=="master"){
			
			broadcast_writer.Write([]byte(ip+"\x00")) // her sender master ipen sin over broadcast
			fmt.Println(ip)
			elevator_number <- elevator_nr   // Må opprette channel
			fmt.Println("heisen har nr: ", elevator_nr, "og har: ", nr_of_slaves, " slaver")
			
			if(nr_of_slaves == 0){
				if CheckConnection(broadcast_listener, ip) == -1{
					order_queue_array[0].Array = order_list_array
					fmt.Println("Unexpected error: Another master on the net. Reassigning elevator state")
					break 
				} 
			}
			
			slave_listener, err := SlaveListener(master_port)
			slave_listener.SetDeadline(time.Now().Add(1500*time.Millisecond))
			connections[nr_of_slaves],err = slave_listener.Accept()
			if err != nil{
				//fmt.Println("Finner ingen slaver: ", err)
			}else{
				nr := strconv.Itoa(nr_of_slaves + 2)
				connections[nr_of_slaves].Write([]byte(nr+"\x00"))  //need int to string conversion
				nr_of_slaves = nr_of_slaves + 1
				fmt.Println(nr_of_slaves)
   			}
   			
			i := 0
			buffer := make([]byte, 128)
			var n int
			for i < nr_of_slaves {
				connections[i].SetDeadline(time.Now().Add(100*time.Millisecond))
				_, err := connections[i].Write([]byte("send\x00"))
				if err != nil{
					connection_timeouts[i] += 1
					if (connection_timeouts[i] > 3){
						fmt.Println("connection timeout for slave: ", i)
						connections[i].Close()
						j := i
						k := nr_of_slaves
						for j < k{
							connections[j] = connections[j+1]
							if(j < nr_of_slaves -1){
								connections[j].SetDeadline(time.Now().Add(1*time.Second))
								_, err1 := connections[j].Write([]byte("decr\x00"))
								if(err1 != nil){
									fmt.Println("error updating connection list")
									continue 
								}
							}
							
							j += 1 
						}
						fmt.Println("Slave number: ", i+1, "is inactive, teminating connection")
						connection_timeouts[i] = 0
						nr_of_slaves-= 1
					}
					continue
				}else{
					connections[i].SetDeadline(time.Now().Add(100*time.Millisecond))
					n, err = connections[i].Read(buffer)
					fmt.Println("slave has sendt package", n)
				}
				var information pablo
				err = json.Unmarshal(buffer[0:n], &information)
				if (err != nil){
					fmt.Println("klarte ikke pakke opp cost, order queue og remover orders: ", err)
					connection_timeouts[i] += 1
					if connection_timeouts[i] > 3{
						fmt.Println("too many wrong messages from slave, terminating connection")
						connections[i].Close()
						j := i
						k := nr_of_slaves
						for j < k{
							connections[j] = connections[j+1]
							if(j < nr_of_slaves -1){
								connections[j].SetDeadline(time.Now().Add(1*time.Second))
								_, err1 := connections[j].Write([]byte("decr\x00"))
								if(err1 != nil){
									fmt.Println("error updating connection list")
									continue 
								}
							}
							
							j += 1 
						}
						connection_timeouts[i] = 0
						nr_of_slaves -= 1
					}		
				}else{
					cost_array[i+1].Array = information.Cost
					order_queue_array[i+1].Array = information.OrderQueue
					remove_order_array[i+1].Array = information.RemoveOrder
				}
				
				i += 1
			}
			//assigning orders
			i = 0
			j := 0
			for i <= nr_of_slaves{
				j = 0
				for j<8{
					if(order_queue_array[i].Array[j] == 1){
						order_list_array[j] = 1
						order_queue_array[i].Array[j] = 0
					}
					j += 1
				}
				j = 0
				for j <8{
					if(remove_order_array[i].Array[j] == 1){
						order_list_array[j] = 0
						remove_order_array[i].Array[j] = 0
					}
					j += 1
				}
				i += 1
			}
			j = 0
			var lowest_cost [8]int
			for j<8{
				i = 0
				lowest_cost[j] = 0
				for i <= nr_of_slaves{
					if(cost_array[i].Array[j] < cost_array[lowest_cost[j]].Array[j]){
						lowest_cost[j] = i
					}
					i += 1
				}
				j += 1
			}
			
			j = 0
			for j<8{ 
				if order_list_array[j] != 0{	
					order_list_array[j] = ((lowest_cost[j] + 1))  // 0 - (nr_of_slaves +1) based on the elevator with lowest cost. 
				}else{
					order_list_array[j] = 0
				}
				j += 1
			}
			
			order_list_yo.Array = order_list_array
			buffer, err = json.Marshal(order_list_yo)
			if (err!=nil){
				fmt.Println("error converting order_list to type driver.Data: ", err)
			}else{
				broadcast_orders.Write(buffer) // broadcast sender order list??
			}
			
			
			order_list <- order_list_yo
			order_list_lights <- order_list_yo
			
			time.Sleep(10*time.Millisecond)	
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
			read_msg := string(buffer)
			elevator_nr_str := strings.Split(read_msg, "\x00")[0]		//
			elevator_nr,_ = strconv.Atoi(elevator_nr_str)
			fmt.Println(elevator_nr, " er det heisen får tilsent :O ")
			elevator_number <- elevator_nr  //	
		}			
//		Slave loop
		for state == "slave"{
			
			
			fmt.Println("Heisen er en slave", elevator_nr)
			buffer := make([]byte, 128)
			
			connections[0].SetDeadline(time.Now().Add(1500*time.Millisecond))
			_, err := connections[0].Read(buffer)
		
			read_msg:=string(buffer)
			if strings.Split(read_msg, "\x00")[0]	 == "send"{
				connection_timeouts[0] = 0	//kompensasjon for å gjøre slaven raskere.
				trans.Cost = cost_array[0].Array
				trans.OrderQueue = order_queue_array[0].Array
				trans.RemoveOrder = remove_order_array[0].Array
				buffer = make([]byte, 128)
				buffer,err = json.Marshal(trans)
				if err != nil{
					fmt.Println("klarte ikke pakke ned cost, order queue og remove orders: ", err)
				}else{
					fmt.Println("Slave sending package: ", cost_array[0].Array) 
					connections[0].Write(buffer)
					order_queue_array[0] = driver.DataInit()
					remove_order_array[0] = driver.DataInit()
				}				
				n, err := recieve_orders.Read(buffer)
				var new_orders driver.Data // her skal man lese order list fra broadcast. Hvor/når sendes ip?
				err = json.Unmarshal(buffer[0:n], &new_orders) // er order_list_yo opprettet før detta?
				if (err != nil){
					fmt.Println("klarte ikke pakke opp order list: ", err)
				}else{
					fmt.Println(new_orders.Array)
					order_list_array = new_orders.Array
					order_list_lights <- new_orders 
					order_list <- new_orders
					
					
				}
				
				//order_list_yo = buffer //converted to driver.Data yay! ++ standariser navn yo
			}else if strings.Split(read_msg, "\x00")[0]	 == "decr"{
				elevator_nr -= 1
				elevator_number <- elevator_nr
			}else if err != nil{
				connection_timeouts[0] += 1
				fmt.Println("Connections timeout: ", err)
				if connection_timeouts[0] > 3{
					connection_timeouts[0] = 0
					fmt.Println("connection timeout, terminating connection to master")
					connections[0].Close()
					break
				}
			}
			
			
			time.Sleep(1*time.Millisecond)
//			HandleMastersOrders()
		}
		time.Sleep(100*time.Millisecond)	
	}
}
 
 
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
	fmt.Println("skal lete etter master nå")
	

	
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
    fmt.Println("den fant noe?: ", master_ip)
	return "slave", master_ip, -1
	
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


func CheckConnection(broadcast_listener *net.UDPConn, ip string)(int){
	buffer := make([]byte, 128)
	var read_ip string	
	broadcast_listener.SetDeadline(time.Now().Add(200*time.Millisecond))		
	_,err := broadcast_listener.Read(buffer)
	if err != nil{
		return 0
	}
	read_ip =string(buffer)
   master_ip_r:= strings.Split(read_ip, "\x00")[0]
	if master_ip_r != ip{
			return -1
	}
	return 1
 
}


