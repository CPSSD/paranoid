package network

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
)

type jsonStruct struct {
	Sender string `json:"sender"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length *int   `json:"length,omitempty"`
	Offset *int   `json:"offset,omitempty"`
	Data   string `json:"data,omitempty"`
	Target string `json:"target,omitempty"`
}

func jsonEncode(structure jsonStruct) []byte {
	structure.Data = base64.StdEncoding.EncodeToString([]byte(structure.Data))
	json, err := json.Marshal(structure)
	if err != nil {
		log.Fatalln(err)
	}
	return []byte(json)
}

func Write(directory, filename string, offset, length *int, data string) {
	machineID, address, port := getMetaInfo(directory)

	writeStruct := jsonStruct{
		Sender: machineID,
		Type:   "write",
		Name:   filename,
		Offset: offset,
		Length: length,
		Data:   data,
	}

	writeMessage := jsonEncode(writeStruct)
	sendMessage(writeMessage, address, port)
}

func Creat(directory, filename string) {
	machineID, address, port := getMetaInfo(directory)

	creatStruct := jsonStruct{
		Sender: machineID,
		Name:   filename,
		Type:   "creat",
	}

	createMessage := jsonEncode(creatStruct)
	sendMessage(createMessage, address, port)
}

func Link(directory, filename, targetName string) {
	machineID, address, port := getMetaInfo(directory)

	linkStruct := jsonStruct{
		Sender: machineID,
		Type:   "link",
		Name:   filename,
		Target: targetName,
	}

	linkMessage := jsonEncode(linkStruct)
	sendMessage(linkMessage, address, port)
}

func Unlink(directory, filename string) {
	machineID, address, port := getMetaInfo(directory)

	uLinkStruct := jsonStruct{
		Sender: machineID,
		Name:   filename,
		Type:   "unlink",
	}

	uLinkMessage := jsonEncode(uLinkStruct)
	sendMessage(uLinkMessage, address, port)
}

func Truncate(directory, filename string, offset int) {
	machineID, address, port := getMetaInfo(directory)

	truncateStruct := jsonStruct{
		Sender: machineID,
		Type:   "truncate",
		Name:   filename,
		Offset: &offset,
	}

	truncateMessage := jsonEncode(truncateStruct)
	sendMessage(truncateMessage, address, port)
}

func getMeta(directory, file string) string {
	fileData, err := ioutil.ReadFile(path.Join(directory, "meta", file))
	if err != nil {
		log.Fatalln(err)
	}
	return string(fileData)
}

func getMetaInfo(directory string) (string, string, string) {
	return getMeta(directory, "uuid"), getMeta(directory, "ip"), getMeta(directory, "port")
}
