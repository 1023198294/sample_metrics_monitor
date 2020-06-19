package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"example/metrics"
	"github.com/shirou/gopsutil/mem" //用于监控CPU和内存的包
	"time"
)

func main(){
	log.Println("hyj monitor")
	http.HandleFunc("/abc", index)
	http.HandleFunc("/usage", usage)
	http.Handle("/metrics", promhttp.Handler())
	metrics.Register()
	go loop()
	err := http.ListenAndServe(":5565", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func loop(){
	for true {
		time.Sleep(time.Second * 1)
		v,_ := mem.VirtualMemory()
		metrics.RequestResourceUpdate(v.UsedPercent)
		log.Println("v.UsedPercent"+strconv.FormatFloat(v.UsedPercent,'E',-1,64))
	}
	
}
func usage(w http.ResponseWriter,r *http.Request){
	timer:=metrics.NewAdmissionLatency()
	metrics.RequestIncrease()
	v,_ := mem.VirtualMemory() //获取CPU占用率情况
	_,err := w.Write([]byte("Total memory:"+strconv.Itoa(int(v.Total))+", Freed memory: "+strconv.Itoa(int(v.Free))+", UsedPercent:"+strconv.FormatFloat(v.UsedPercent,'E',-1,64)+"% \n"))
	if err!=nil{
		log.Println("err:"+err.Error()+" Yes\n")
	}
	timer.Observe()
}
func index(w http.ResponseWriter, r *http.Request) {
	timer:=metrics.NewAdmissionLatency()
	metrics.RequestIncrease()
	num:=os.Getenv("Num")
	if num==""{
		Fibonacci(10)
		_,err:=w.Write([]byte("there is no env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" No\n")
		}
	}else{
		numInt,_:=strconv.Atoi(num)
		Fibonacci(numInt)
		_,err:=w.Write([]byte("there is env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" Yes\n")
		}
	}
	timer.Observe()
}

func Fibonacci(n int)int{
	if n<=2{
		return 1
	}else{
		return Fibonacci(n-1)+Fibonacci(n-2)
	}
}
