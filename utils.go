package main

import (
)


var CurrentMessageID = 0

func GenerateMessageID() uint16 {
	if CurrentMessageID != 65535 {
		CurrentMessageID++
	} else {
		CurrentMessageID = 1
	}
	return uint16(CurrentMessageID)
}



