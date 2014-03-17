package iomodule

import(
   "driver"
   "fmt"
   "time"
)

func IoManager(order_queue chan driver.Data, command_list chan driver.Data, remove_command chan driver.Data, order_list_lights chan driver.Data){
   
   //Initialize data structures. 
   var order_queue_copy driver.Data
   order_queue_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var order_list_copy driver.Data
   order_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var command_list_copy driver.Data
   command_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var remove_commands driver.Data
   remove_commands.Array = [8]int{0,0,0,0,0,0,0,0}
   
   //Gorutine for reading relevant channels 
   go func(){
      for{
         select{
         case data := <- order_list_lights:
            order_list_copy = data
         case data := <- remove_command:
            remove_commands = data         
         }
         time.Sleep(1*time.Millisecond)
      }
   }()
   
   for {
   
      //Stores button signals in arrays. 
      i := 0
      for i<4{
         if driver.GetButtonSignal("command", i) == 1{
            command_list_copy.Array[i] = 1
         } 
         if driver.GetButtonSignal("down", i) == 1{
            order_queue_copy.Array[2*i] = 1
         }
         if driver.GetButtonSignal("up", i) == 1{
            fmt.Println("knapp nr: ", i, " er blitt satt til 1")
            order_queue_copy.Array[2*i + 1] = 1
         }
         i += 1
         
      }
      
      
      // Sets lights on according to input. 
      i= 0 
      for (i < 4){
         driver.SetButtonLamp("command", i, command_list_copy.Array[i])
         if (i == driver.GetFloor()){
            driver.SetFloorIndicator(i)
         }
         if (i > 0) && (order_list_copy.Array[i*2] != 0) {   
            driver.SetButtonLamp("down", i , 1)
         }else if (i >0){
         	driver.SetButtonLamp("down", i , 0)
         }
         if (i < 3)&&  (order_list_copy.Array[i*2+1] != 0){
            driver.SetButtonLamp("up", i, 1)
         }else if (i<3){
         	driver.SetButtonLamp("up", i, 0)
         }
         i += 1
      }
      
      //Removes orders and sets bits low in the remove orders arrays because they no longer need to be removed. 
      
      i = 0
      for i<4{
         if remove_commands.Array[i] == 1{
            command_list_copy.Array[i] = 0
            remove_commands.Array[i] = 0
         } 
         /*
         if remove_orders[2*i] == 1{
            order_list_copy.Array[2*i] = 0
            remove_orders[2*i] = 0
         }
         if remove_orders[2*i+1] == 1{
            order_list_copy.Array[2*i+1] =  0
            remove_orders[2*i+1] = 0
         }
         */
         i += 1
      }
      
      //Writing the lists to other files. (is outdated and needs som change)
      
     
     command_list <- command_list_copy
     
     order_queue <- order_queue_copy
     
     order_queue_copy = driver.DataInit()
     
       
   time.Sleep(1*time.Millisecond)
   }

}

