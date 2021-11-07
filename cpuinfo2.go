package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	//"time"
)

// cpu info
func getCpuInfo() {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
	}

	/*for _, ci := range cpuInfos {
		fmt.Println(ci)
	}*/

	//fmt.Println(cpuInfos.ModelName)

	fmt.Println("CPU 具体型号:\t", cpuInfos[0].ModelName)      //d
	c, _ := cpu.Counts(false)                              //true: cpu逻辑数量,     false:CPU物理数量
	fmt.Println("CPU 物理数量:\t", c)                          //0
	c1, _ := cpu.Counts(true)                              //true: cpu逻辑数量,     false:CPU物理数量
	fmt.Println("CPU 逻辑数量:\t", c1)                         //
	raminfo, _ := mem.VirtualMemory()                      //	raminfo, _ := mem.SwapMemory()
	fmt.Println("内存大小(M):\t", raminfo.Total/1024/1024)     //
	swapraminfo, _ := mem.SwapMemory()                     //	raminfo, _ := mem.SwapMemory()
	fmt.Println("交换内存(M):\t", swapraminfo.Total/1024/1024) //

	// CPU使用率
	//for {
	//percent, _ := cpu.Percent(time.Second, false)
	//fmt.Printf("cpu percent:%v\n", percent)
	//}

	//c, _ := cpu.Counts(false)  //true: cpu逻辑数量,     false:CPU物理数量
	//fmt.Println("cpu逻辑数量:", c) //4

}

func main() {
	getCpuInfo()
}
