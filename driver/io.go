package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
*/
import "C"

C.static C.comedi_t* it_g := nil


func io_init()(int){
    status := 0

    it_g := comedi_open("/dev/comedi0")
  
    if (it_g == nil)
        return 0

    i := 0
    for i < 8 {
        status |= comedi_dio_config(it_g, PORT1, i,     COMEDI_INPUT)
        status |= comedi_dio_config(it_g, PORT2, i,     COMEDI_OUTPUT)
        status |= comedi_dio_config(it_g, PORT3, i+8,   COMEDI_OUTPUT)
        status |= comedi_dio_config(it_g, PORT4, i+16,  COMEDI_INPUT)
        i += 1 
    }

    return (status == 0)            //kanskje må gjøres om til if
}



func io_set_bit(channel int){
    comedi_dio_write(it_g, channel >> 8, channel & 0xff, 1)
}



func io_clear_bit(int channel){
    comedi_dio_write(it_g, channel >> 8, channel & 0xff, 0)
}



func io_write_analog(channel int ,value int){
    C.comedi_data_write(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, value)
}



func io_read_bit(channel int)(int){
    data := 0
    C.comedi_dio_read(it_g, channel >> 8, channel & 0xff, &data)
    return int(data)
}



func io_read_analog(int channel)(int){
    var data C.lsampl_t 
    C.comedi_data_read(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, &data)
    return int(data)
}

