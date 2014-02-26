package main

import (
    "./driver" 
    "fmt"
    "time"
)

func main(){
    fmt.Printf("du står i: ", driver.GetFloor(), " etasje")
    driver.SetFloorIndicator(driver.GetFloor())
    driver.SetSpeed(20)
    time.Sleep(1*time.Second)
    driver.SetSpeed(0)
    driver.SetButtonLamp("down", 1, 1)
    driver.SetButtonLamp("down", 0, 1)
    driver.SetDoorLamp(1)

}

