package goshare

import golhashmap "github.com/abhishekkr/gol/golhashmap"

/*
ReadKey gets value of given key.
*/
func ReadKey(key string) golhashmap.HashMap {
	var hashmap golhashmap.HashMap
	hashmap = make(golhashmap.HashMap)
	val := tsds.GetVal(key)
	if val == "" {
		return hashmap
	}
	hashmap[key] = val
	return hashmap
}

/*
ReadKeyNS gets value for all descendents of given key's namespace.
*/
func ReadKeyNS(key string) golhashmap.HashMap {
	return tsds.ReadNSRecursive(key)
}

/*
ReadKeyTSDS gets value for the asked time-frame key, aah same NS.
*/
func ReadKeyTSDS(key string) golhashmap.HashMap {
	return tsds.ReadTSDS(key)
}

/*
ReadFuncByKeyType calls a read task on task-type.
*/
func ReadFuncByKeyType(keyType string) FunkAxnParamKeyReturnMap {
	switch keyType {
	case "tsds":
		return ReadKeyTSDS

	case "ns":
		return ReadKeyNS

	default:
		return ReadKey

	}
}

/*
ReadFromPacket calls ReadFuncByKeyType for multi-keys based on provided packet.
*/
func ReadFromPacket(packet Packet) string {
	var response string
	var hashmap golhashmap.HashMap
	hashmap = make(golhashmap.HashMap)

	axnFunk := ReadFuncByKeyType(packet.KeyType)
	for _, _key := range packet.KeyList {
		hashmap = axnFunk(_key)
		if len(hashmap) == 0 {
			continue
		}
		response += responseByValType(packet.ValType, hashmap)
	}

	return response
}

/* transform response by ValType, if none default:csv */
func responseByValType(valType string, responseMap golhashmap.HashMap) string {
	var response string

	switch valType {
	case "csv", "json":
		hashmapEngine := golhashmap.GetHashMapEngine(valType)
		response = hashmapEngine.FromHashMap(responseMap)

	default:
		response = golhashmap.HashMapToCSV(responseMap)
	}
	return response
}
