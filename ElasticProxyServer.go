package main

import (
	//"fmt"
	"github.com/sgjp/go-coap"
	"log"
	"net"
)

var Conf Configuration

func startServer(cIn chan *coap.Message) {


	log.Fatal(coap.ListenAndServeMulticast("udp", "224.0.1.187:5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			log.Printf("Got message path=%q: PayLoad: %#v from %v proxy-uri: %v", m.Path(), string(m.Payload), a, m.Option(coap.ProxyURI))
			if len(m.Path()) > 0 {

				switch m.Path()[0] {

				case "fwd":
					res := fwdHandler(m,cIn)
					return res


				default:
					res := notFoundHandler(m)
					return res

				}
			} else {
				res := notFoundHandler(m)
				return res
			}
			return nil
		}),"en1"))
	log.Printf("FINISHED SERVER")

}


func notFoundHandler(m *coap.Message) *coap.Message {

	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.NotFound,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   []byte("4.05"),
	}
	res.SetOption(coap.ContentFormat, coap.TextPlain)
	return res

}



func fwdHandler(m *coap.Message, cIn chan *coap.Message) *coap.Message{
	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.Valid,
		MessageID: GenerateMessageID(),
		Token:     m.Token,
		Payload:   []byte("2.05"),
	}
	//log.Printf("WRITING TO CHANNEL")
	cIn<-m
	//log.Printf("WROTE, SIZE: %v",len(cIn))

	return res
}





