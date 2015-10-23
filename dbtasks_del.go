package goshare

/*
DelKey deletes val for a given key, returns status.
*/
func DelKey(key string) bool {
	return tsds.DelKey(key)
}

/*
DelKeyNS deletes given key's namespace and all its values, returns status.
*/
func DelKeyNS(key string) bool {
	return tsds.DeleteNSRecursive(key)
}

/*
DelKeyTSDS deletes all keys under given namespace, same as NS.
As here TimeSeries is a NameSpace
*/
func DelKeyTSDS(key string) bool {
	return tsds.DeleteTSDS(key)
}

/*
DeleteFuncByKeyType calls a delete action for a key based on task-type.
*/
func DeleteFuncByKeyType(keyType string) FunkAxnParamKey {
	switch keyType {
	case "tsds":
		return DelKeyTSDS

	case "ns":
		return DelKeyNS

	default:
		return DelKey

	}
}

/*
DeleteFromPacket can handle multi-keys delete action,
it acts on packet data.
*/
func DeleteFromPacket(packet Packet) bool {
	status := true
	axnFunk := DeleteFuncByKeyType(packet.KeyType)
	for _, _key := range packet.KeyList {
		status = status && axnFunk(_key)
	}

	return status
}
