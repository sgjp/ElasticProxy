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
			log.Printf("Got message path=%q: PayLoad: %#v from %v proxy-uri: %v. Code: %v", m.Path(), string(m.Payload), a, m.Option(coap.ProxyURI),m.Code)
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

	if m.Code==coap.GET{
		entry, cached := getEntry(m)

		//Valid is non chacheable, here we use it to not cache the response of calling calcPrimeNumberResult
		if cached {
			log.Printf("CACHED !")
			res := &coap.Message{
				Type:                coap.Acknowledgement,
				Code:                coap.Content,
				MessageID:           m.MessageID,
				Token:               m.Token,
				Payload:             []byte(entry.ResponsePayload),
			}
			res.SetOption(coap.MaxAge, 60)
			res.SetOption(coap.ContentFormat, coap.TextPlain)

			return res
		}else{
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Valid,
				MessageID: GenerateMessageID(),
				Token:     m.Token,
				Payload:   []byte("2.05"),
			}
			addOriginalRequest(m)
			cIn<-m

			return res
		}

	}else if m.Code == coap.PUT{
		//If the code is PUT, the request is a response to another request
		originalRequest, exists:=getOriginalRequest(m)
		if exists{
			//The response belongs to a request, must be cached
			cacheEntry := genNewCacheEntry(originalRequest,m,originalRequest.Option(coap.ProxyURI).(string))
			addEntry(cacheEntry)
		}else{
			log.Printf("-----ENTRY NOT ADDED----")
		}

		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: GenerateMessageID(),
			Token:     m.Token,
			Payload:   []byte("2.05"),
		}


		cIn<-m
		return res

	}else{
		//Its not GET or PUT, respond and forward
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: GenerateMessageID(),
			Token:     m.Token,
			Payload:   []byte("2.05"),
		}

		cIn<-m
		return res
	}

}





