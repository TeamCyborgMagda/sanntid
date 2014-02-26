package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.

import(
    "math"
    "fmt"
)


var last_speed = 0
const N_FLOORS = 4

var button_adress = [...]int{
    {FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
    {FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
    {FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
    {FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
}

var lamp_matrix = [...]int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

func Init() int {
    if (IoInit() == 0){
        return 0
    }
    i := 0
    for i < N_FLOORS { // fiks variabel
        if i != 0 {
            SetButtonLamp(BUTTON_CALL_DOWN, i, 0)
        }
        if (i != N_FLOORS)
            SetButtonLamp(BUTTON_CALL_UP, i, 0)
        }
        SetButtonLamp(BUTTON_COMMAND, i, 0);
    }
    
    SetStorLamp(0)
    SetDoorLamp(0)
    SetFloorIndicator(0)
}



func SetSpeed(speed int){
    
    if (speed > 0){
        IoClearBit(MOTORDIR)
    } else if (speed < 0) {
        IoSetBit(MOTORDIR)
    } else if (lastSpeed < 0){
        IoClearBit(MOTORDIR)
    } else if (lastSpeed > 0){
        IoSetBit(MOTORDIR)
    }
    
    lastSpeed = speed
    
    if (speed > 0){
        IoWriteAnalog(MOTOR, 2048 + 4*speed)
    } else {
        IoWriteAnalog(MOTOR, 2048 - 4*speed)
}

func GetFloor() int {
    if (IoReadBit(SENSOR1)){
        return 0
    } else if (IoReadBit(SENSOR2)){
        return 1
    }else if (IoReadBit(SENSOR3)){
        return 2
    }else if (IoReadBit(SENSOR4)){
        return 3
    }else{
        return -1
    }
}

func GetButtonSignal(button string, floor int){
    if (button == "down" && floor >= 1 && floor < N_FLOORS){
        if IoReadBit(button_adress[floor][1] ){
            return 1
        }else{
            return 0
        }
    }else if(button == "up" && floor < (N_FLOORS-1) && floor >= 0){
         if IoReadBit(button_adress[floor][0] ){
            return 1
        }else{
            return 0
        }
    }else if(button == "command" && floor < N_FLOORS && floor >= 0){
        if IoReadBit(button_adress[floor][2] ){
            return 1
        }else{
            return 0
        }
    }

}

func SetButtonLamp(button string, floor int, value int){
    if button == "up" || button == "down" || button == "command"{
    
        if  floor < N_FLOORS && floor >= 0{
            if(button == "down" && floor == 0){
                fmt.Printf("Button nonexisting for floor")  
                return
            }
            if(button == "up" && floor == N_FLOORS -1){
                fmt.Printf("Button nonexisting for floor")  
                return
            }
            if (value == 1){
                IoSetBit(lamp_matrix[floor][button]);
            }else{
                IoClearBit(lamp_matrix[floor][button]);
            }
    
        }
    }
}


func SetFloorIndicator(floor int){
    if floor < 0  || floor > N_FLOORS{
        fmt.Printf("Floor value error")
    }
    // Binary encoding. One light must always be on.
    if floor == 0{
        IoClearBit(FLOOR_IND1)
        IoClearBit(FLOOR_IND2)
    }else if floor ==1 {
        IoSetBit(FLOOR_IND1)
        IoClearBit(FLOOR_IND2)
    }else if floor == 2{
        IoSetBit(FLOOR_IND2)
        IoClearBit(FLOOR_IND1)
    }else if floor ==3 {
        IoSetBit(FLOOR_IND1)
        IoSetBit(FLOOR_IND2)
    }
}

func SetDoorLamp(value int){
    if value == 1{
        IoSetBit(DOOR_OPEN)
    }else{
        IoClearBit(DOOR_OPEN)
    }
}

func SetStopLamp(value int){
    if value == 1{
        IoSetBit(LIGHT_STOP)
    }else{
        IoClearBit(LIGHT_STOP)
    }
}


