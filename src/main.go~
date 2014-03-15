package main

import(
   "driver"
   //"network"
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
   driver.Init()
   fmt.Printf("Iomodule i choose you\n")
   go iomodule.IoManager(order_queue, command_list, order_list, cost, remove_order, remove_command)
   fmt.Printf("Heis i choose you\n")
   go heis.Heis(order_list, command_list, cost, remove_order, remove_command)
  
   for {
      
      time.Sleep(time.Second)
   }
   
}
