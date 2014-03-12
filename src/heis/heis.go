package heis

import(
   "driver"
   "time"  
   "math"
   "fmt"
)

func HeisInit()(int, int, int){
   direction := 0;
   for driver.GetFloor() == -1 {
      driver.SetSpeed(-300)
   }
   driver.SetSpeed(0)
   current_floor := driver.GetFloor()
   destination := -1
   return direction, current_floor, destination
}

func Heis(order_list chan driver.Data, command_list chan driver.Data, cost chan driver.Data, remove_order chan driver.Data, remove_command chan driver.Data){ 
   
   direction, current_floor, destination := HeisInit()
   var cost_copy driver.Data
   var command_list_copy driver.Data
   var command_list_temp driver.Data
   var order_list_copy driver.Data
   var order_list_temp driver.Data
   var remove_orders driver.Data
   var remove_commands driver.Data
   order_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   command_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   
   go func(){
      for{
         select{
         case data := <- order_list:
            fmt.Printf("skal lese knapp to opp om den er lagt inn\n")
            order_list_temp = data
            order_list_copy = order_list_temp
         case data := <- command_list:
            command_list_temp = data
            command_list_copy = command_list_temp
          //  fmt.Println(data.Array)
         case data := <- cost:
            cost_copy = data
           // fmt.Println(data.Array)
         default:
            time.Sleep(1)   
         }
      }
   }()
   
   for{
      //order_list_copy, command_list_copy := ReadIo(order_list, command_list)
      //direction is initialized to zero, this function returns the first found destination if the direction
      // is zero, and optimalizes the destination if the direction is positive or negative. 
      destination = GetDestination(direction, current_floor, order_list_copy.Array, command_list_copy.Array)
      fmt.Println(destination, " denne skal være 1")
      // decides direction required to reach destination from current floor.
      direction = GetDirection(destination, current_floor)
      
      driver.SetSpeed(direction*300)
      for(destination != -1){
         
            
         //   cost_temp = <- cost
          //  cost_temp.Array = CostFunction(current_floor, direction, destination)
          //  cost <- cost_temp
         floor := driver.GetFloor() 
         if(floor != -1){
            current_floor = floor
         }
         
         //order_list_copy, command_list_copy = ReadIo(order_list, command_list)
         if( (direction==-1 && order_list_copy.Array[2*current_floor]==1) || (direction==1 && order_list_copy.Array[2*current_floor+1]==1) || command_list_copy.Array[current_floor] == 1 || (destination == current_floor)){
            driver.SetSpeed(0)
            driver.SetDoorLamp(1)
            
            
            
            remove_orders.Array, remove_commands.Array = RemoveOrders(direction, destination)
            fmt.Printf("Heis Deadlock yo?\n")
                   
            fmt.Println("heis skriver")
            remove_order <- remove_orders
            fmt.Println("heis skriver2")
            remove_command <- remove_commands
                                                              
            fmt.Printf("Heis yo no\n")
            
            time.Sleep(3*time.Second)
            driver.SetDoorLamp(0)
            if current_floor == destination{
               destination = -1
            }
         }
         time.Sleep(1*time.Millisecond)
      }
      // direction
      // 2) else if there are orders in order list. Complete them until
      time.Sleep(1*time.Millisecond) 
   }
}

func GetDirection(destination int, current_floor int)(int){
   direction := 0
   if (destination == -1){
      return direction
   }else if(destination > current_floor){
      direction = 1  
   }else if(destination < current_floor){
      direction = -1
   }
   return direction
}

func GetDestination(direction int, current_floor int, order_list [8]int, command_list [8]int)(int){
   var i int
   if(direction == 1){
      i = 3
      for(i >= current_floor){
         if (order_list[i*2+1] == 1 || command_list[i] == 1){
            return i
         }
         i -= 1 
      }
      return -1
         
   }else if (direction == -1){
      i = 0
      for(i <= current_floor){
         if (order_list[i*2] == 1 || command_list[i] == 1){
            return i
         }
         i += 1
      }
      return -1
         
         //hvis ikke behold det som destination
   }else{
      i = 0
      for(i < 4){
         if (order_list[i*2] == 1 || order_list[i*2+1] == 1 || command_list[i] == 1){
            return i
         }
         i += 1
      }
      return -1
      
         //sjekk, command lista, så order lista(?) sett første som finnes til destination.
         //hvis ikke, sett destination til eller 0 eller noe. 
   }

}


func CostFunction(current_floor int,direction int, destination int)([8]int){
   i := 0
   var cost [8]int
   for i<8{
      if (direction == 0){
         cost[i] = int(math.Abs(float64(i/2 - current_floor)))    //ABSOLUTT VERDI
      }else if(direction == 1){
         if(i%2 == 1 && i/2 > current_floor){
            cost[i] = i/2 - current_floor - 1
         }else if (i%2 == 1 && i/2 <= current_floor || i%2 == 0){
            cost[i] =int (math.Abs(float64(i/2 - destination)) + math.Abs(float64(destination - current_floor - 1)))
         }else{
            cost[i] = 6
         }
      }else{
         if(i%2 == 0 && i/2 < current_floor){
            cost[i] =  current_floor - i/2 - 1
         }else if (i%2 == 0 && i/2 >= current_floor || i%2 ==  1){
            cost[i] = int(math.Abs(float64(i/2 - destination)) + math.Abs(float64(current_floor- destination - 1)))
         }else{
            cost[i] = 6
         }
      }
      i += 1
   }
   return cost
}

func RemoveOrders(direction int, destination int)([8]int,[8]int){
   remove_order := [8]int{0,0,0,0,0,0,0,0}
   remove_command :=[8]int{0,0,0,0,0,0,0,0}
   i := 0
   for (i < 4){
      if (driver.GetFloor() == i){
         remove_command[i] = 1
         if (destination == i){
         	remove_command[i] = 1
         	remove_order[i*2] = 1
         	remove_order[i*2+1] = 1
         	
         }else if (direction == 1){
            remove_order[i*2+1] = 1
         } else if (direction == 0){
            remove_order[i*2] = 1
         } 
      }
      i +=1
   }
   return remove_order, remove_command
}
/*
func ReadIo(order_list chan driver.Data, command_list chan driver.Data)(driver.Data, driver.Data){
   var order_list_copy driver.Data
   order_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var command_list_copy driver.Data
   command_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   i:= 0
   for i < 4{
      i += 1
      select {       
      case data := <- order_list:
         order_list_copy = data   
      case data := <- command_list:
         command_list_copy = data
      default:
      }
   }
   return  order_list_copy, command_list_copy
}

*/

