package web

import (
	"encoding/csv"
	"encoding/json"
	hosieweb "github.com/hoisie/web"
	"io/ioutil"
	"log"
	"logbuzz/config"
	"logbuzz/data"
	"logbuzz/list"
	"logbuzz/persist"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func PostInput(ctx *hosieweb.Context) {

	rawData, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		log.Println("Error reading request body: " + err.Error())
		ctx.WriteHeader(400)
		return
	}

	logItem := data.LogItem{}

	err = json.Unmarshal(rawData, &logItem)

	if err != nil {
		log.Printf("unable to parse '%v' into a logItem because %v\n", string(rawData), err.Error())
		ctx.WriteHeader(400)
		return
	}

	if logItem.TagList == nil {
		logItem.TagList = make([]string, 1)
	}
	logItem.TagList = append(logItem.TagList, "ip:"+strings.Split(ctx.Request.RemoteAddr, ":")[0])

	go func() {
		list.AddLogItem(&logItem)

		// we also persist the log item
		persist.PersistLogItem(&logItem)
	}()
}

func GetQuery(ctx *hosieweb.Context) {

	var err error

	var tagList []string = nil

	if ctx.Params["tagList"] != "" {
		tagList, err = csv.NewReader(strings.NewReader(ctx.Params["tagList"])).Read()
		if err != nil {
			log.Println("Unable to parse tagList '%v' into a valid CSV list because %v\n", ctx.Params["tagList"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	level := ctx.Params["level"]

	skip := 0
	if ctx.Params["skip"] != "" {
		skip, err = strconv.Atoi(ctx.Params["skip"])
		if err != nil {
			log.Println("Unable to parse skip '%v' into a valid number because %v\n", ctx.Params["skip"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	fromTS := int64(0)
	if ctx.Params["fromTS"] != "" {
		fromTS, err = strconv.ParseInt(ctx.Params["fromTS"], 10, 64)
		if err != nil {
			log.Println("Unable to parse fromTS '%v' into a valid timestamp because %v\n", ctx.Params["fromTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	toTS := time.Now().UnixNano() / 1000

	if ctx.Params["toTS"] != "" {
		toTS, err = strconv.ParseInt(ctx.Params["toTS"], 10, 64)
		if err != nil {
			log.Println("Unable to parse toTS '%v' into a valid timestamp because %v\n", ctx.Params["toTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	logItemList := list.Search(level, tagList, fromTS, toTS, skip)

	byteArray, _ := json.Marshal(logItemList)

	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(byteArray)
}

func GetStats(ctx *hosieweb.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ret := make(map[string]interface{})
	ret["itemCount"] = list.GetNumItems()
	ret["systemMemory"] = m.HeapSys - m.HeapReleased

	startTS, endTS := list.GetTimestampRange()

	ret["firstTimestamp"] = startTS
	ret["lastTimestamp"] = endTS

	byteArray, _ := json.Marshal(ret)

	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(byteArray)
}

func PostConfigHttpPort(ctx *hosieweb.Context) {
	rawData, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		log.Println("Error reading request body: " + err.Error())
		ctx.WriteHeader(400)
		return
	}

	var httpPort int = 0

	err = json.Unmarshal(rawData, &httpPort)

	if err != nil {
		log.Println("unable to parse '", string(rawData), "' into a valid httpPort value because:", err.Error())
		ctx.WriteHeader(400)
		return
	}

	if httpPort < 1023 || httpPort > 65534 {
		log.Println("Invalid port number:", httpPort)
		ctx.WriteHeader(400)
		return
	}

	config.UpdateHttpPort(httpPort)
}

func PostConfigMemoryLimit(ctx *hosieweb.Context) {
	rawData, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		log.Println("Error reading request body: " + err.Error())
		ctx.WriteHeader(400)
		return
	}

	var memoryLimit uint64 = 0

	err = json.Unmarshal(rawData, &memoryLimit)

	if err != nil {
		log.Println("unable to parse '", string(rawData), "' into a valid memoryLimit value because:", err.Error())
		ctx.WriteHeader(400)
		return
	}

	if memoryLimit < 5*1024*1024 {
		log.Println("Invalid memory limit:", memoryLimit)
		ctx.WriteHeader(400)
		return
	}

	config.UpdateMemoryLimit(memoryLimit)
}

func GetConfig(ctx *hosieweb.Context) {
	byteArray, _ := json.Marshal(config.GetConfig())

	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(byteArray)
}
