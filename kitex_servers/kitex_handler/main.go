package main

import (
	echosvc "Orbital_Hebao/kitex_servers/kitex_gen/echo/echosvc"
	sumsvc "Orbital_Hebao/kitex_servers/kitex_gen/sum/sumsvc"
	"log"
)

func main() {
	svr1 := echosvc.NewServer(new(EchoImpl))
	svr2 := sumsvc.NewServer(new(SumImpl))

	err1 := svr1.Run()
	if err1 != nil {
		log.Println(err1.Error())
	}

	err2 := svr2.Run()
	if err2 != nil {
		log.Println(err2.Error())
	}
}
