package main

import (
    //"iomodule" 
    "fmt"
    "heis"
    //"driver"
    //"encoding/json"
    
    //"time"
)



func main(){

	current_floor := 1
	direction := 1
	destination := 3
	
	fmt.Println(heis.CostFunction(current_floor,direction,destination))

/*
	nr_of_slaves := 1
	order_list_array := [8]int{0,0,0,0,0,0,0,0}
	var (
		remove_order_array [2]driver.Data
		order_queue_array [2]driver.Data
		cost_array [2]driver.Data
		remove1 driver.Data
		remove2 driver.Data
		cost1 driver.Data
		cost2 driver.Data
		queue1 driver.Data
		queue2 driver.Data
	)
	remove1.Array = [8]int{0,0,0,0,0,0,0,0}
	remove2.Array = [8]int{0,0,0,0,0,0,0,0}
	cost1.Array = [8]int{4,4,3,1,2,0,1,1}
	cost2.Array = [8]int{0,0,1,1,2,2,3,3}
	queue1.Array = [8]int{0,0,0,0,0,0,0,0}
	queue2.Array = [8]int{0,0,1,1,0,0,0,0}
	order_queue_array[0] = queue1
	order_queue_array[1] = queue2
	remove_order_array[0] = remove1
	remove_order_array[1] = remove2
	cost_array[0] = cost1
	cost_array[1] = cost2
	
	
	i := 0
	j := 0
	for i <= nr_of_slaves{
		j = 0
		for j<8{
			if(order_queue_array[i].Array[j] == 1){
				order_list_array[j] = 1
			}
			j += 1
		}
		j = 0
		for j <8{
			if(remove_order_array[i].Array[j] == 1){
				order_list_array[j] = 0
			}
			j += 1
		}
		i += 1
	}
	j = 0
	var lowest_cost [8]int
	for j<8{
		i = 0
		lowest_cost[j] = 1
		for i<nr_of_slaves{
			if(cost_array[i].Array[j] < cost_array[lowest_cost[j]].Array[j]){
				lowest_cost[j] = i
			}
			i += 1
		}
		j += 1
	}
	j = 0
	for j<8{ 	
		order_list_array[j] = order_list_array[j]*(lowest_cost[j] + 1) // 0 - (nr_of_slaves +1) based on the elevator with lowest cost. 
		j += 1
	}
	fmt.Println(order_list)
*/
}

