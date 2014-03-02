package main

import (
    "./driver" 
    "fmt"
    //"time"
)

func CostFunction(current_floor int,direction int, destination int)([8]int){
   i := 0
   var cost [8]int
   for i<8{
      if (direction == 0){
         cost[i] = i/2 - current_floor    //ABSOLUTT VERDI
      }else if(direction == 1){
         if(i%2 == 1 && i/2 > current_floor){
            cost[i] = i/2 - current_floor - 1
         }else{
            cost[i] = 6 
         }
      }else{
         if(i%2 == 0 && i/2 < current_floor){
            cost[i] = current_floor - i/2 - 1
         }else{
            cost[i] = 6
         }
      }
      i += 1
   }
   i=0
   for i<8{
      if (cost[i] < (current_floor- destination) - (destination - i/2) ){  //ABSOLUTT VERDI I KLAMMENE
         cost[i] = (current_floor- destination) - (destination - i/2)
      }
      i += 1
   }
   
   return cost

}

func main(){
    fmt.Printf("hei")
    cost := CostFunction(0, 1, 3)
    fmt.Printf("YO")
    fmt.Println(cost)
    
    
    
    
    fmt.Println("du står i: ", driver.GetFloor(), " etasje \n")
    
/*
    driver.Init()
    fmt.Printf("du står i: ", driver.GetFloor(), " etasje")
    driver.SetFloorIndicator(driver.GetFloor())
    driver.SetSpeed(20)
    time.Sleep(1*time.Second)
    driver.SetSpeed(0)
    driver.SetButtonLamp("down", 1, 1)
    driver.SetButtonLamp("down", 0, 1)
    driver.SetDoorLamp(1)
*/
}

// GetFloor funker ikke
// SetSpeed funker
// GetButtonSignal funker
// Init funker
// SetButtonLamp funker
// GetFloor funker
// SetFloorIndicator funker
