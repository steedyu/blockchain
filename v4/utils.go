package main

import (
	"encoding/binary"
	"bytes"
	"os"
	"fmt"
)

func IntToBytes(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	CheckErr("IntToBytes", err)
	return buffer.Bytes()
}

func CheckErr(pos string, err error) {
	if err != nil {
		fmt.Println("error, pos:", err, pos)
		os.Exit(1)
	}
}
