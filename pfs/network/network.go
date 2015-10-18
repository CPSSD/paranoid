package network

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
)

type jsonStruct struct {
	Sender string `json:"sender"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length int    `json:"length"`
	Offset int    `json:"offset"`
	Data   string `json:"data"`
	Target string `json:"target"`
}

func jsonEncode(structure jsonStruct) []byte {
	json, err := json.Marshal(structure)
	if err != nil {
		log.Fatalln(err)
	}
	return []byte(json)
}

func Write(directory, file_name string, offset int, length int, data, address, port string) {
	machineID := getUUID(directory)

	creatStruct := jsonStruct{
		Sender: machineID,
		Type:   "write",
		Name:   file_name,
		Offset: offset,
		Length: length,
		Data:   data,
	}

	createMessage := jsonEncode(creatStruct)
	SendMessage(createMessage, address, port)
}

func Creat(directory, filename, address, port string) {
	//TODO validation on data passed to func
	machineID := getUUID(directory)

	creatStruct := jsonStruct{
		Sender: machineID,
		Name:   filename,
		Type:   "creat",
	}

	createMessage := jsonEncode(creatStruct)
	SendMessage(createMessage, address, port)
}

func Link(directory, filename, targetName, address, port string) {
	machineID := getUUID(directory)

	creatStruct := jsonStruct{
		Sender: machineID,
		Type:   "link",
		Name:   filename,
		Target: targetName,
	}

	createMessage := jsonEncode(creatStruct)
	SendMessage(createMessage, address, port)
}

func Unlink(directory, filename, address, port string) {
	machineID := getUUID(directory)

	uLinkStruct := jsonStruct{
		Sender: machineID,
		Name:   filename,
		Type:   "unlink",
	}

	uLinkMessage := jsonEncode(uLinkStruct)
	SendMessage(uLinkMessage, address, port)
}

func Truncate(directory, address, port, filename string, offset int) {
	machineID := getUUID(directory)

	truncateStruct := jsonStruct{
		Sender: machineID,
		Type:   "truncate",
		Name:   filename,
		Offset: offset,
	}

	truncateMessage := jsonEncode(truncateStruct)
	SendMessage(truncateMessage, address, port)
}

func getUUID(directory string) string {
	uuid, err := ioutil.ReadFile(path.Join(directory, "meta", "uuid"))
	if err != nil {
		log.Fatalln(err)
	}
	return string(uuid)
}
