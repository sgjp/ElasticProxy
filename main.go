package main

import "github.com/sgjp/go-coap"

func main() {

	Conf = getConfiguration()





	c := make(chan *coap.Message, 100000)
	go startClient(c)
	startServer(c)


}
