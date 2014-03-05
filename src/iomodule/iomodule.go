package iomodule

import(
   "driver"
   "fmt"
//   "time"
)

func IoManager(order_queue chan driver.Data, command_list chan driver.Data, order_list chan driver.Data, cost chan driver.Data){
   for {
      fmt.Printf("while loop io \n")
      // READ INPUT
     // order_queue_copy := <- order_queue
    //  order_queue <- order_queue_copy
      
      command_list_copy := <- command_list
      
      
      order_list_copy := <- order_list
      
    //  cost_copy := <- cost
      
      fmt.Printf("begynner for loop io \n")
      i:= 0
      for i<4{
         if driver.GetButtonSignal("command", i) == 1{
            command_list_copy.Array[i] = 1
         } 
         if driver.GetButtonSignal("down", i) == 1{
            order_list_copy.Array[2*i] = 1
         }
         if driver.GetButtonSignal("up", i) == 1{
            order_list_copy.Array[2*i + 1] =driver.GetButtonSignal("up", i)
         }
         i += 1
         
      }
      
     // temp = <- order_queue
      //order_queue <- order_queue_copy
     // temp = <- command_list
    //  command_list <- command_list_copy
     // fmt.Println(temp.Array)
   
      // Panel thingy
      
      
      i= 0 
      for (i < 4){
         driver.SetButtonLamp("command", i, command_list_copy.Array[i])
         if (i > 0){   
            driver.SetButtonLamp("down", i , order_list_copy.Array[i*2])
         }
         if (i < 3){
            driver.SetButtonLamp("up", i, order_list_copy.Array[i*2+1])
         }
         i += 1
      }
   order_list <- order_list_copy
   command_list <- command_list_copy
   //cost <- cost_copy
   
   }

}
