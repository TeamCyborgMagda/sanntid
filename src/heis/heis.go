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

func Heis(order_list chan driver.Data, command_list chan driver.Data, cost chan driver.Data){ 
   direction, current_floor, destination := HeisInit()
 //  cost <- init   
   
   for{
      fmt.Printf("while loop heis: check \n")
    //  cost_temp := <- cost
    //  cost_temp.Array = CostFunction(current_floor, direction, destination)
    //  cost <- cost_temp
    //  fmt.Printf("Første kopiering i heis: check \n")
      
      order_list_copy := <- order_list
      fmt.Printf("Andre kopiering i heis: check \n")
      
      command_list_copy := <- command_list
      
      fmt.Printf("Tredje kopiering i heis: check \n")
      
      //direction is initialized to zero, this function returns the first found destination if the direction
      // is zero, and optimalizes the destination if the direction is positive or negative. 
      destination = GetDestination(direction, current_floor, order_list_copy.Array, command_list_copy.Array)
      
      // decides direction required to reach destination from current floor.
      direction = GetDirection(destination, current_floor)
      order_list <- order_list_copy
      command_list <- command_list_copy
      
      driver.SetSpeed(direction*300)
      for(destination != -1){
         if(driver.GetFloor() == -1){
            continue
         }else{
            current_floor = driver.GetFloor()
            
         //   cost_temp = <- cost
          //  cost_temp.Array = CostFunction(current_floor, direction, destination)
          //  cost <- cost_temp
         }
         
         order_list_copy = <- order_list
         order_list <- order_list_copy
         command_list_copy <- command_list
         command_list <- command_list_copy
         
         
         command_list <- command_list_copy
         if( (direction==-1 && order_list_copy.Array[2*current_floor-2]==1) || (direction==1 && order_list_copy.Array[2*current_floor-1]==1) || command_list_copy.Array[current_floor] == 1){
            driver.SetSpeed(0)
            driver.SetDoorLamp(1)
            //removeOrders(current_floor, direction)
            time.Sleep(3*time.Second)
            driver.SetDoorLamp(0)
            if current_floor == destination{
               destination = -1
            }
            break
         }
      }
      // direction
      // 2) else if there are orders in order list. Complete them until 
      
   
   
   
   }
}

func GetDirection(destination int, current_floor int)(int){
   direction := 0
   if(destination == -1){
      direction = 0
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
         
   }else if(direction ==  -1){
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

