package goshare

/* TBD-or-MayBeDone
		   	   make it match all "<default|ns|tsds>(-<csv|json>)(-<now>)"
		   	   1. pull out <default|ns|tsds>(-<now>) only as done currently
		   	   2. <default|ns|tsds>-<csv|json>(-<now>) driven by encoding converted to hashmap, looped for call from (1)
		        2.1 encoding for tsds (not now) should contain root time config, overridden by any key-val-time config
            2.2 encoding for ns may contain a root parent-ns config to be used, overridden by any key-val-parent config
*/

import (
	"strings"

	golhashmap "github.com/abhishekkr/gol/golhashmap"
	goltime "github.com/abhishekkr/gol/goltime"
	levigoNS "github.com/abhishekkr/levigoNS"
	abkleveldb "github.com/abhishekkr/levigoNS/leveldb"
	levigoTSDS "github.com/abhishekkr/levigoTSDS"
)

/* Push a given set of Key-Val */
func PushKeyVal(key string, val string) bool {
	return abkleveldb.PushKeyVal(key, val, db)
}

/* Push a given Namespace-Key and its value */
func PushKeyValNS(key string, val string) bool {
	return levigoNS.PushNS(key, val, db)
}

/* Push a key namespace-d with goltime.Timestamp  */
func PushKeyValTSDS(key string, val string, timestamp goltime.Timestamp) bool {
	if levigoTSDS.PushTSDS(key, val, timestamp.Time(), db) == "" {
		return false
	}
	return true
}

/* Push a key namespace-d with current time */
func PushKeyValNowTSDS(key string, val string) bool {
	if levigoTSDS.PushNowTSDS(key, val, db) == "" {
		return false
	}
	return true
}

/* handles single-Item; delegates multi-item */
func PushKeyValSolo(task_type string, key string, value string, message_array *[]string) bool {

	switch task_type {
	case "tsds":
		timestamp := goltime.CreateTimestamp((*message_array)[2:8])
		_key := (*message_array)[8]
		_value := strings.Join((*message_array)[9:], " ")
		return PushKeyValTSDS(_key, _value, timestamp)

	case "now":
		return PushKeyValNowTSDS(key, value)

	case "ns":
		return PushKeyValNS(key, value)

	default:
		return PushKeyVal(key, value)

	}
}

/* handles multi-item */
func PushKeyValMulti(task_type string, multi_type string, message_array *[]string) bool {
	var hashmap_key_value golhashmap.HashMap
	timestamp := goltime.CreateTimestamp((*message_array)[2:8])
	multi_value := strings.Join((*message_array)[8:], "\n")

	switch multi_type {
	case "csv":
		hashmap_key_value = golhashmap.Csv_to_hashmap(multi_value)

		/*make multi_type sent to golhashmap and get converter, pass multi_value and get hashmap*/

	//case "json":
	//	hashmap_key_value = golhashmap.Json_to_hashmap(multi_value)

	default:
		return false
	}

	status := true
	_time := timestamp.Time()
	for _key, _val := range hashmap_key_value {
		_val = strings.Replace(_val, "\n", " ", -1)
		PushKeyValSolo(task_type, _key, _val, message_array)
	}
	return status
}

/* Push a key-val based on task-type and multi|solo-type*/
func PushKeyValByType(task_type string, message_array []string) bool {
	_key := message_array[2]
	_value := strings.Join(message_array[3:], " ")

	task_type_tokens := strings.Split(task_type, "-")

	if len(task_type_tokens) == 2 {
		return PushKeyValMulti(task_type_tokens[0], task_type_tokens[1], &message_array)
	} else if len(task_type_tokens) == 1 {
		return PushKeyValSolo(task_type_tokens[0], _key, _value, &message_array)
	} else {
		return false
	}
}
