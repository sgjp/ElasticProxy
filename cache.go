package main

import(
	"github.com/sgjp/go-coap"
	"time"
)


var cache []CacheEntry



func addEntry(cacheEntry CacheEntry){
	cache = append(cache,cacheEntry)
}

func getEntry(req coap.Message, host string, path string) (CacheEntry, bool){
	var cacheEntry CacheEntry

	for _, entry := range cache {
		if entry.RequestPath == path && entry.RequestCode == req.Code && entry.RequestHost == host && entry.RequestType == req.Type && entry.RequestPayload == string(req.Payload){

			maxLiveDuration := entry.ResponseMaxLive.(uint32)
			t := time.Now()
			maxLive := time.Minute * time.Duration(maxLiveDuration)
			if t.Before(entry.Timestamp.Add(maxLive)){
				cacheEntry = entry
				return cacheEntry, true
			}

		}
	}
	return cacheEntry, false


}






func genNewCacheEntry(req coap.Message, res *coap.Message, host string, path string) CacheEntry{
	cacheEntry := CacheEntry{
		RequestType:req.Type,
		RequestCode:req.Code,
		RequestHost:host,
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
	RequestHost string
	RequestPath string
	RequestPayload string

	ResponseType coap.COAPType
	ResponseCode coap.COAPCode
	ResponseMaxLive interface{}
	ResponsePayload string

	Timestamp time.Time
}