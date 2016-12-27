//All calls to other services originated from the proxy should be done in this file

package main

import (
	"github.com/sgjp/go-coap"
	"log"
)



func startClient(channel chan *coap.Message ) *coap.Message{

	for{



		m:=<-channel
		//log.Printf("READ DONE. NEW SIZE: %v",len(channel))

		if(m.Option(coap.ProxyURI)==nil){
			log.Printf("Proxy-URI not set for request: %v. PAYLOAD: %v",m,string(m.Payload))
		}else{
			//TODO Avoid getting the same resource as the requester
			resource := getResourceByName(m.Option(coap.ProxyURI).(string))
			log.Printf("Sending")
			sendCoapMsg(resource.Path,resource.Host,string(m.Payload),m.Code,m.Type)
			log.Printf("SENT")
			//cacheEntry := genNewCacheEntry(req, rv, resource.Host, resource.Path)

			//addEntry(cacheEntry)

			//return rv
		}

	}



}

func sendCoapMsg(path string, host string, payload string, coapCode coap.COAPCode, coapType coap.COAPType) *coap.Message{
	req := coap.Message{
		Type:      coapType,
		Code:      coapCode,
		MessageID: GenerateMessageID(),
		Payload:   []byte(payload),
	}


	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)

	c, err := coap.Dial("udp", host)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		//log.Printf("Response payload: %s", rv.Payload)
	}
	return rv
}

