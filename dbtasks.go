package goshare

/*
[PATTERN]
action {read, push, delete}
type {default, ns, tsds, now}

## message_array here is devoided of axn and key_type
non-tsds {key&val, :type-data}
tsds(-*) {tdot&key&val, tdot&:type-data}
*/

// gotta refactor more and make this front for these tasks
func DBTasks(axn string, key_type string, message_array []string) ([]byte, bool) {
	response := ""
	axn_status := false

	key := message_array[0]

	switch axn {
	case "read":
		response = GetValTask(key_type, key)
		if response != "" {
			axn_status = true
		}

	case "push":
		if PushKeyValByType(key_type, message_array) {
			axn_status = true
		}

	case "delete":
		if DelKeyTask(key_type, key) {
			axn_status = true
		}
	}

	return []byte(response), axn_status
}
