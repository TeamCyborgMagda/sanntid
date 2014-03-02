package main

import(
   "./driver"
   "./network"
)



func main(){
   // channel ordertable
   // lese av ordre fil.
   go network(order_list, cost)   //bør inneholde kost function
   go io(order_list, command_list)
   go heis(order_list, command_list, cost)
   sleep(10000*seconds)
}
