package main

import (
	"encoding/json"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
	"io/ioutil"
	"net/http"
	"bytes"
)


func postJson(urlstr string, params map[string]interface{}) (string, error) {
	var err error
	var resp *http.Response
	jsonPost := JsonEncode(params)
	requestBody := bytes.NewBuffer([]byte(jsonPost))
	request, err := http.NewRequest("POST", urlstr, requestBody)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	resp, err = client.Do(request)

	if err != nil || resp == nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}
	return string(body), nil
}

func sha1Encode(input string) string{
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

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
