package goshare

import (
	"fmt"
	"strings"

	golhashmap "github.com/abhishekkr/gol/golhashmap"
	gollist "github.com/abhishekkr/gol/gollist"
	goltime "github.com/abhishekkr/gol/goltime"
)

/*
Packet for modelling data passed to GoShare into a structure of possible fields.
*/
type Packet struct {
	DBAction string
	TaskType string

	KeyType string // key: default, namespace key: ns, timeseries key: tsds, timeseries for goshare time: now
	ValType string // single: default, csv, json

	HashMap golhashmap.HashMap
	KeyList []string

	ParentNS string // allowed for ns|tsds|now
	TimeDot  goltime.Timestamp
}

/*
FunkAxnParamKeyVal is a function type which get passed two string parameters
and returns one boolean. Like Push Key-Val calls.
*/
type FunkAxnParamKeyVal func(key string, val string) bool

/*
FunkAxnParamKey is a function type which get passed one string parameter
and returns one boolean. Like Del Key tasks.
*/
type FunkAxnParamKey func(key string) bool

/*
FunkAxnParamKeyReturnMap is a function type which get passed one string parameters
and returns one hashmap. Like Get Key tasks.
*/
type FunkAxnParamKeyReturnMap func(key string) golhashmap.HashMap

/*
CreatePacket formulates Packet structure from passed message array.
*/
func CreatePacket(packetArray []string) Packet {
	packet := Packet{}
	packet.HashMap = make(golhashmap.HashMap)

	lenPacketArray := len(packetArray)
	if lenPacketArray < 3 {
		packet.DBAction = "ERROR"
		return packet
	}

	packet.DBAction = packetArray[0]
	packet.TaskType = packetArray[1]
	dataStartsFrom := 2

	taskTypeTokens := strings.Split(packet.TaskType, "-")
	packet.KeyType = taskTypeTokens[0]
	if packet.KeyType == "tsds" && packet.DBAction == "push" {
		if lenPacketArray < 9 {
			packet.DBAction = "ERROR"
			return packet
		}
		packet.TimeDot = goltime.CreateTimestamp(packetArray[2:8])
		dataStartsFrom += 6
	}

	if len(taskTypeTokens) > 1 {
		packet.ValType = taskTypeTokens[1]

		if len(taskTypeTokens) == 3 {
			// if packet requirement grows more than 3, that's the limit
			// go get 'msgpack' to handle it instead...
			thirdTokenFeature(&packet, packetArray, &dataStartsFrom, taskTypeTokens[2])
		}
	}

	decodeData(&packet, packetArray[dataStartsFrom:])
	return packet
}

/* thirdTokenFeature handles special 3rd token feature, to populate Packet data. */
func thirdTokenFeature(packet *Packet, packetArray []string, dataStartsFrom *int, token string) {
	switch token {
	case "parent":
		packet.ParentNS = packetArray[*dataStartsFrom]
		(*dataStartsFrom)++
	}
}

/*
decodeData handles Packet formation based on DBAction.
Handles TimeDot and pre-pending of ParentNS.
*/
func decodeData(packet *Packet, messageArray []string) {
	switch packet.DBAction {
	case "read", "delete":
		packet.KeyList = decodeKeyData(packet.ValType, messageArray)
		if packet.ParentNS != "" {
			PrefixKeyParentNamespace(packet)
		}

	case "push":
		packet.HashMap = decodeKeyValData(packet.ValType, messageArray)
		if packet.ParentNS != "" {
			PrefixKeyValParentNamespace(packet)
		}

	default:
		packet.DBAction = "ERROR"
	}
}

/*
decodeKeyData handles Packet formation based on valType for GET, DELETE.
*/
func decodeKeyData(valType string, messageArray []string) []string {
	switch valType {
	case "csv", "json":
		multiValue := strings.Join(messageArray, "\n")
		listEngine := gollist.GetListEngine(valType)
		return listEngine.ToList(multiValue)

	default:
		return []string{messageArray[0]}
	}
}

/*
decodeKeyValData handles Packet formation based on valType for PUSH.
*/
func decodeKeyValData(valType string, messageArray []string) golhashmap.HashMap {
	var hashmap golhashmap.HashMap
	hashmap = make(golhashmap.HashMap)

	switch valType {
	case "csv", "json":
		multiValue := strings.Join(messageArray, "\n")
		hashmapEngine := golhashmap.GetHashMapEngine(valType)
		hashmap = hashmapEngine.ToHashMap(multiValue)

	default:
		key := messageArray[0]
		value := strings.Join(messageArray[1:], " ")
		hashmap[key] = value
	}
	return hashmap
}

/*
PrefixKeyParentNamespace prefixes Parent Namespaces to all keys in List
if val for 'parentNamespace'.
*/
func PrefixKeyParentNamespace(packet *Packet) {
	var newList []string
	newList = make([]string, len(packet.KeyList))

	parentNamespace := packet.ParentNS
	for idx, key := range packet.KeyList {
		newList[idx] = fmt.Sprintf("%s:%s", parentNamespace, key)
	}
	packet.KeyList = newList
}

/*
PrefixKeyValParentNamespace prefixes Parent Namespaces to all key-val in HashMap
if it has val for 'parentNamespace'.
*/
func PrefixKeyValParentNamespace(packet *Packet) {
	var newHmap golhashmap.HashMap
	newHmap = make(golhashmap.HashMap)

	parentNamespace := packet.ParentNS
	for key, val := range packet.HashMap {
		newHmap[fmt.Sprintf("%s:%s", parentNamespace, key)] = val
	}
	packet.HashMap = newHmap
}
