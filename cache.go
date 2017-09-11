package main

import(
	"github.com/sgjp/go-coap"
	"time"
	"sync"
	"log"
)


var cache []CacheEntry

var requestList []*coap.Message

var mutex = &sync.Mutex{}

func addEntry(cacheEntry CacheEntry){
	cache = append(cache,cacheEntry)
	log.Printf("ENTRY ADDED! REQ PY: %v, RESP PY: %v, new LEN: %v, PATH: %v",string(cacheEntry.RequestPayload),string(cacheEntry.ResponsePayload),len(cache),cacheEntry.RequestPath)

}

func getEntry(req *coap.Message) (CacheEntry, bool){
	var cacheEntry CacheEntry

	for _, entry := range cache {
		//log.Printf("PATHS: %v, %v, %v",entry.RequestPath,req.Option(coap.ProxyURI),entry.RequestPath == req.Option(coap.ProxyURI))
		//log.Printf("CODES: %v, %v, %v",entry.RequestCode ,req.Code,entry.RequestCode == req.Code)
		if entry.RequestPath == req.Option(coap.ProxyURI) && entry.RequestCode == req.Code && entry.RequestType == req.Type && entry.RequestPayload == string(req.Payload){
			//log.Printf("Matches !")
			//maxLiveDuration := entry.ResponseMaxLive.(uint32)
			//t := time.Now()
			//maxLive := time.Minute * time.Duration(maxLiveDuration)
			//if t.Before(entry.Timestamp.Add(maxLive)){
				cacheEntry = entry
				return cacheEntry, true
			//}

		}
	}
	return cacheEntry, false


}


func addOriginalRequest (req *coap.Message){

	//log.Printf("--ADDING ORIGINAL REQUEST WITH ID:%v, CURRENT SIXE: %v",req.MessageID,len(requestList))
	mutex.Lock()
	requestList = append(requestList,req)
	mutex.Unlock()
	//log.Printf("--ADDING ORIGINAL REQUEST WITH ID:%v, NEW SIXE: %v",req.MessageID,len(requestList))
}

func getOriginalRequest(req *coap.Message) (*coap.Message, bool){
	mutex.Lock()
	for i, request := range requestList {
		if request.MessageID==req.MessageID &&
			request.Token[0]==req.Token[0] &&
			request.Token[1]==req.Token[1] &&
			request.Token[2]==req.Token[2] &&
			request.Token[3]==req.Token[3] &&
			request.Token[4]==req.Token[4] &&
			request.Token[5]==req.Token[5] &&
			request.Token[6]==req.Token[6] &&
			request.Token[7]==req.Token[7] {
			requestList = append(requestList[:i], requestList[i+1:]...)
			//originalRequest = request
			//log.Printf("----GETTING ORIGINAL REQUEST FOR IDDD:%v, NEW SIZE: %v",req.MessageID,len(requestList))
			mutex.Unlock()
			return request, true
		}else{
			//log.Printf("----NOT GETTING ORIGINAL REQUEST for IDD:%v, NEW SIZE: %v",req.MessageID,len(requestList))
		}

	}
	mutex.Unlock()
	var request coap.Message
	return &request,false
}



func genNewCacheEntry(req *coap.Message, res *coap.Message, path string) CacheEntry{
	cacheEntry := CacheEntry{
		RequestType:req.Type,
		RequestCode:req.Code,
		RequestPath:path,
		RequestPayload: string(req.Payload),
		ResponseType: res.Type,
		ResponseCode: res.Code,
		ResponseMaxLive: res.Option(coap.MaxAge),
		ResponsePayload: string(res.Payload),
		Timestamp: time.Now(),
	}

	return cacheEntry
}

type CacheEntry struct{
	RequestType coap.COAPType
	RequestCode coap.COAPCode
	RequestPath string
	RequestPayload string

	ResponseType coap.COAPType
	ResponseCode coap.COAPCode
	ResponseMaxLive interface{}
	ResponsePayload string

	Timestamp time.Time
}