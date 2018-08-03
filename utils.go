package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateUnixNanoId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func JsonEncode(nodes interface{}) string {
	body, err := json.Marshal(nodes)
	if err != nil {
		panic(err.Error())
	}
	return string(body)
}

func JsonDecode(jsonStr string) interface{} {
	jsonStr = strings.Replace(jsonStr, "\n", "", -1)
	var f interface{}
	err := json.Unmarshal([]byte(jsonStr), &f)
	if err != nil {
		panic(err)
	}
	return f
}

func UptimeFormat(secs uint32, section int) string {
	timeUnits := map[uint32][2]string{
		1:     [2]string{"second", "seconds"},
		60:    [2]string{"minute", "minutes"},
		3600:  [2]string{"hour", "hours"},
		86400: [2]string{"day", "days"},
	}
	timeSeq := [4]uint32{86400, 3600, 60, 1}
	timeSeqLen := len(timeSeq)
	result := make([]string, timeSeqLen)
	if section < 1 {
		section = 1
	} else if section > timeSeqLen {
		section = timeSeqLen
	}

	i := 0
	for _, index := range timeSeq {
		if v, prs := timeUnits[index]; prs {
			if secs >= index {
				num := secs / index
				secs = secs % index
				unit := v[0]
				if num > 1 {
					unit = v[1]
				}
				result[i] = fmt.Sprintf("%d %s", num, unit)
				i++
			}
		}
	}
	sliceLen := i
	if sliceLen > section {
		sliceLen = section
	}
	return strings.Join(result[0:sliceLen], ", ")
}
