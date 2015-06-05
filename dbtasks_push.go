package goshare

import (
	"strings"

	golkeyvalNS "github.com/abhishekkr/gol/golkeyvalNS"
	golkeyvalTSDS "github.com/abhishekkr/gol/golkeyvalTSDS"
)

/*
PushKeyVal pushes a given set of Key-Val.
*/
func PushKeyVal(key string, val string) bool {
	return db.PushKeyVal(key, val)
}

/*
PushKeyValNS pushes a given Namespace-Key and its value.
*/
func PushKeyValNS(key string, val string) bool {
	return golkeyvalNS.PushNS(key, val, db)
}

/*
PushKeyValNowTSDS pushes a key namespace-d with current time.
*/
func PushKeyValNowTSDS(key string, val string) bool {
	return golkeyvalTSDS.PushNowTSDS(key, val, db)
}

/*
PushKeyValTSDS pushes a key namespace-d with goltime.Timestamp.
*/
func PushKeyValTSDS(packet Packet) bool {
	status := true
	_time := packet.TimeDot.Time()
	for _key, _val := range packet.HashMap {
		_val = strings.Replace(_val, "\n", " ", -1)
		status = status && golkeyvalTSDS.PushTSDS(_key, _val, _time, db)
	}
	return status
}

/*
PushFuncByKeyType returns func handle according to KeyType.
*/
func PushFuncByKeyType(keyType string) FunkAxnParamKeyVal {
	switch keyType {
	case "now":
		return PushKeyValNowTSDS

	case "ns":
		return PushKeyValNS

	default:
		return PushKeyVal

	}
}

/*
PushFromPacket handles push task based on provided Packet.
*/
func PushFromPacket(packet Packet) bool {
	status := true
	switch packet.KeyType {
	case "tsds":
		PushKeyValTSDS(packet)

	default:
		axnFunk := PushFuncByKeyType(packet.KeyType)
		for _key, _val := range packet.HashMap {
			_val = strings.Replace(_val, "\n", " ", -1)
			status = status && axnFunk(_key, _val)
		}
	}

	return status
}
