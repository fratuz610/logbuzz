package web

import (
	"encoding/json"
	"fmt"
	hosieweb "github.com/hoisie/web"
	"io/ioutil"
	"logbuzz/logitem"
	"strconv"
	"time"
)

func PostInput(ctx *hosieweb.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		fmt.Println("Error reading request body: " + err.Error())
		ctx.WriteHeader(400)
		return
	}

	logItem := logitem.LogItem{}

	err = json.Unmarshal(data, &logItem)

	if err != nil {
		fmt.Printf("unable to parse '%v' into a logItem because %v\n", string(data), err.Error())
		ctx.WriteHeader(400)
		return
	}

	go logitem.AddLogItem(&logItem)
}

func GetQuery(ctx *hosieweb.Context) {

	inputID := ctx.Params["inputID"]

	if inputID == "" {
		fmt.Println("No inputID specified")
		ctx.WriteHeader(400)
		return
	}

	env := ctx.Params["env"]
	level := ctx.Params["level"]

	var err error

	skip := 0
	if ctx.Params["skip"] != "" {
		skip, err = strconv.Atoi(ctx.Params["skip"])
		if err != nil {
			fmt.Println("Unable to parse skip '%v' into a valid number because %v\n", ctx.Params["skip"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	fromTS := int64(0)
	if ctx.Params["fromTS"] != "" {
		fromTS, err = strconv.ParseInt(ctx.Params["fromTS"], 10, 64)
		if err != nil {
			fmt.Println("Unable to parse fromTS '%v' into a valid timestamp because %v\n", ctx.Params["fromTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	toTS := time.Now().UnixNano() / 1000

	if ctx.Params["toTS"] != "" {
		toTS, err = strconv.ParseInt(ctx.Params["toTS"], 10, 64)
		if err != nil {
			fmt.Println("Unable to parse toTS '%v' into a valid timestamp because %v\n", ctx.Params["toTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	logItemList := logitem.Search(inputID, level, env, fromTS, toTS, skip)

	byteArray, err := json.Marshal(logItemList)

	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(byteArray)
}
