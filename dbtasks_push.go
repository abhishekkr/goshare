package goshare

import (
	"strings"

	golkeyvalNS "github.com/abhishekkr/gol/golkeyvalNS"
	golkeyvalTSDS "github.com/abhishekkr/gol/golkeyvalTSDS"
)

/* Push a given set of Key-Val */
func PushKeyVal(key string, val string) bool {
	return db.PushKeyVal(key, val)
}

/* Push a given Namespace-Key and its value */
func PushKeyValNS(key string, val string) bool {
	return golkeyvalNS.PushNS(key, val, db)
}

/* Push a key namespace-d with current time */
func PushKeyValNowTSDS(key string, val string) bool {
	return golkeyvalTSDS.PushNowTSDS(key, val, db)
}

/* Push a key namespace-d with goltime.Timestamp  */
func PushKeyValTSDS(packet Packet) bool {
	status := true
	_time := packet.TimeDot.Time()
	for _key, _val := range packet.HashMap {
		_val = strings.Replace(_val, "\n", " ", -1)
		status = status && golkeyvalTSDS.PushTSDS(_key, _val, _time, db)
	}
	return status
}

/* return func handle according to KeyType */
func PushFuncByKeyType(key_type string) FunkAxnParamKeyVal {
	switch key_type {
	case "now":
		return PushKeyValNowTSDS

	case "ns":
		return PushKeyValNS

	default:
		return PushKeyVal

	}
}

/* handles multi-item */
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
