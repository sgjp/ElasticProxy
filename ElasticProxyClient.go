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
			if resource.Id==0{
				log.Fatalf("No available resources found for path: %v. Finishing",m.Option(coap.ProxyURI).(string))
			}

			_, success:=sendCoapMsg(resource.Path,resource.Host,string(m.Payload),m.Code,m.Type,m.MessageID,m.Token)

			if !success{
				removeResource(resource)
				channel <- m
			}
		}

	}



}

func sendCoapMsg(path string, host string, payload string, coapCode coap.COAPCode, coapType coap.COAPType, messageId uint16, token []byte) (*coap.Message, bool){
	req := coap.Message{
		Type:      coapType,
		Code:      coapCode,
		MessageID: messageId,
		Token:     token,
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
		//log.Fatalf("Error sending request: %v", err)
		log.Printf("Error calling resource: %v",err)
		return rv, false
	}

	if rv != nil {
		//log.Printf("Response payload: %s", rv.Payload)
	}
	return rv, true
}

