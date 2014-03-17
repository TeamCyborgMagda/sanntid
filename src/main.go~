package main

import(
   "driver"
   "network"
   "iomodule"
   "fmt"
   "heis"
   "time"
)



func main(){
   
   order_queue := make(chan driver.Data)
   command_list := make(chan driver.Data)
   order_list := make(chan driver.Data)
   cost := make(chan driver.Data)
   remove_order := make(chan driver.Data)
   remove_command :=  make(chan driver.Data)
   elevator_number := make(chan int)
   order_list_lights := make(chan driver.Data)
   driver.Init()
   fmt.Printf("Iomodule i choose you\n")
   go iomodule.IoManager(order_queue, command_list, remove_command, order_list_lights)
   fmt.Printf("Heis i choose you\n")
   go heis.Heis(order_list, command_list, cost, remove_order, remove_command, elevator_number)
   fmt.Printf("go network\n")
   go network.Network(order_queue, remove_order, cost, elevator_number, order_list, order_list_lights)
  
   for {
      
      time.Sleep(time.Second)
   }
   
}
