package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.

import(
//    "math"
    "fmt"
)


var last_speed = 1
const N_FLOORS = 4

type Data struct {                        //Burde deklareres i en fil som alle modulene bruker. (typ driver)
    Array [8]int
}

func DataInit()(Data){
   var data Data
   data.Array = [8]int{0,0,0,0,0,0,0,0}
   return data
}


var button_adress = [12]int{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1, FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2, FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3, FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4}
// 0-2 = first floor, 3-5 = second floor, 6-8 = third floor, 9-11 = fourth floor

var lamp_matrix = [12]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1, LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2, LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3, LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}
// 0-2 = first floor, 3-5 = second floor, 6-8 = third floor, 9-11 = fourth floor


func Init() int {
    if (IoInit() == 0){
        return 0
    }
    i := 0
    for i < N_FLOORS { // fiks variabel
        if i != 0 {
            SetButtonLamp("down", i, 0)
        }
        if (i != N_FLOORS-1){
            SetButtonLamp("up", i, 0)
        }
        SetButtonLamp("command", i, 0);
        i += 1
    }
    
    SetStopLamp(0)
    SetDoorLamp(0)
    SetFloorIndicator(0)
    return 1
}



func SetSpeed(speed int){
    
    if (speed > 0){
        IoClearBit(MOTORDIR)
    } else if (speed < 0) {
        IoSetBit(MOTORDIR)
    } else if (last_speed < 0){
        IoClearBit(MOTORDIR)
    } else if (last_speed > 0){
        IoSetBit(MOTORDIR)
    }
    
    last_speed = speed
    
    if (speed > 0){
        IoWriteAnalog(MOTOR, 2048 + 4*speed)
    } else {
        IoWriteAnalog(MOTOR, 2048 - 4*speed)
   }
}

func GetFloor() (int){
    if (IoReadBit(SENSOR1) == 1){
        return 0
    } else if (IoReadBit(SENSOR2) == 1){
        return 1
    }else if (IoReadBit(SENSOR3)== 1){
        return 2
    }else if (IoReadBit(SENSOR4)== 1){
        return 3
    }else{
        return -1
    }
}

func GetButtonSignal(button string, floor int) (int){
    if (button == "down" && floor >= 1 && floor < N_FLOORS){
        if (IoReadBit(button_adress[3*floor + 1])==1){
            return 1
        }
    }else if(button == "up" && floor < (N_FLOORS-1) && floor >= 0){
         if (IoReadBit(button_adress[3*floor])==1 ){
            return 1
        }
    }else if(button == "command" && floor < N_FLOORS && floor >= 0){
        if (IoReadBit(button_adress[3*floor + 2])==1 ){
            return 1
        }
    }
    return 0
}

func SetButtonLamp(button string, floor int, value int){
    if button == "up" || button == "down" || button == "command"{
        if  floor < N_FLOORS && floor >= 0{
            if(button == "down" && floor == 0){
                fmt.Printf("Button nonexisting for floor (down) \n")  
                return
            }
            if(button == "up" && floor == N_FLOORS -1){
                fmt.Println("Button nonexisting for floor (up) \ni = ", floor)  
                return
            }
            if (value == 1){
                if (button == "up"){
                    IoSetBit(lamp_matrix[3*floor])
                }else if (button == "down"){
                    IoSetBit(lamp_matrix[3*floor+1])
                }else if (button == "command"){
                    IoSetBit(lamp_matrix[3*floor+2])
                }
            }else{
            if (button == "up"){
                    IoClearBit(lamp_matrix[3*floor])
                }else if (button == "down"){
                    IoClearBit(lamp_matrix[3*floor+1])
                }else if (button == "command"){
                    IoClearBit(lamp_matrix[3*floor+2])
                }
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
        IoSetBit(FLOOR_IND2)
        IoClearBit(FLOOR_IND1)
    }else if floor == 2{
        IoSetBit(FLOOR_IND1)
        IoClearBit(FLOOR_IND2)
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




