package main

import (
    "iomodule" 
    "fmt"
    "heis"
    "driver"
    "time"
    "runtime"
)


/*
func ChannelInit( c1 chan driver.Data, c2 chan driver.Data, c3 chan driver.Data, c4 chan driver.Data ){
    c1 <- driver.Data{[8]int{0,0,0,0,0,0,0,0}}
    c2 <- *new(driver.Data)
    c3 <- *new(driver.Data)
    c3 <- *new(driver.Data)
}
*/

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU()) 
    driver.Init()  
    cost := make(chan driver.Data)
    order_list := make(chan driver.Data)
    order_queue := make(chan driver.Data)
    command_list := make(chan driver.Data)
    
    
    go heis.Heis(order_list, command_list, cost)
    go iomodule.IoManager(order_queue, command_list, order_list, cost)
    
    
    var init driver.Data
    init.Array = [8]int{0,0,0,0,0,0,0,0}
    order_list <- init
    command_list <- init
    time.Sleep(1000*time.Millisecond)
    
    
 //   ChannelInit(order_list, order_queue, command_list, cost)
//    command_list <- *init
//    order_queue <-  *init
//    command_list <- *init
//    cost <- *init
    
   // go iomodule.PanelLights(order_list, command_list)
    
    /*i := 0
    for i < 15{
      time.Sleep(1000*time.Millisecond)
      fmt.Println("og: ", i)  
      i +=1
    }*/
    
    fmt.Printf("\nwut slutten \n")
    //fmt.Println(command_list.Array)
  
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

// GetFloor funker
// SetSpeed funker
// GetButtonSignal funker
// Init funker
// SetButtonLamp funker
// GetFloor funker
// SetFloorIndicator funker
