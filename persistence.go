package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"encoding/csv"
	"strconv"
	"net"
	"strings"
)

var confFile = "sgjp/ElasticProxy/conf"

var taskDurationFile = "sgjp/ElasticProxy/TaskDuration.csv"

//Create the required files if they don't exist
func initPersistence() {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		writeToFile(confFile, "")
	}
	if _, err := os.Stat(taskDurationFile); os.IsNotExist(err) {
		writeToFile(taskDurationFile, "")
	}
}
func getConfiguration() Configuration {
	confValue := readFromFile(confFile)
	var configuration Configuration

	json.Unmarshal([]byte(confValue), &configuration)
	return configuration

}


func addResource(path string, a *net.UDPAddr) {
	splittedAddress := strings.Split(a.String(),":")
	conf := getConfiguration()

	for _, res :=range conf.Resources{
		if res.Path == path && res.Host == splittedAddress[0]+":5683"{
			return
		}
	}

	var resource Resource

	resource.Host = splittedAddress[0]+":5683"
	resource.Path = path
	//The following conditional is in order to set the Id to the new resource
	if len(conf.Resources) > 0 {
		resource.Id = (conf.Resources[len(conf.Resources)-1].Id)+1
	} else {
		resource.Id = 1

	}


	resources := append(conf.Resources, resource)
	conf.Resources = resources

	confPlain, _ := json.Marshal(conf)
	writeToFile(confFile, string(confPlain))

}


func readFromFile(fileName string) string {

	stream, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	readString := string(stream)

	return readString
}

func writeToFile(fileName string, data string) {

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(data)


	file.Close()
}

func appendToFile(fileName string, data string){

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(data)
	f.WriteString("\n")
	f.Close()

}

func getResourceByName(resourceName string) Resource {
	/*conf := getConfiguration()

	var resource Resource

	for _, res := range conf.Resources {
		if res.Path == resourceName {
			resource = res
		}
	}*/


	var resource Resource

	var key int
	for k,res := range Conf.Resources{
		if res.Path == resourceName{
			if resource.Id<1{
				resource=res
				key = k
			}else if res.Used<resource.Used{
				resource=res
				key=k
			}
		}

	}
	resource.Used++
	Conf.Resources[key]=resource
	return resource
}

func removeResource (resource Resource){

	for k,res := range Conf.Resources{
		if res == resource{
			log.Printf("DELETING RESOURCE: %v", resource)
			Conf.Resources = append(Conf.Resources[:k], Conf.Resources[k+1:]...)
		}

	}
}



func saveTaskDuration(elapsed time.Duration, resource Resource){
	record := []string{
		elapsed.String(),strconv.Itoa(resource.Id),resource.Host,resource.Path,
	}

	file, er := os.OpenFile(taskDurationFile, os.O_RDWR|os.O_APPEND, 0666)

	if er != nil {
		log.Fatal(er)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	err := writer.Write(record)


	if err != nil {
		log.Fatal(er)
	}

	defer writer.Flush()
}

type Configuration struct {
	Resources []Resource
}
type Resource struct {
	Id   int
	Host string
	Path string
	Used int
}
