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

	fromTS := int64(0)
	var err error

	if ctx.Params["fromTS"] != "" {
		fromTS, err = strconv.ParseInt(ctx.Params["fromTS"], 10, 64)
		if err != nil {
			fmt.Println("Unable to parse fromTS '%v' into a valid timestamp because %v\n", ctx.Params["fromTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	toTS := time.Now().Unix()

	if ctx.Params["toTS"] != "" {
		toTS, err = strconv.ParseInt(ctx.Params["toTS"], 10, 64)
		if err != nil {
			fmt.Println("Unable to parse toTS '%v' into a valid timestamp because %v\n", ctx.Params["toTS"], err.Error())
			ctx.WriteHeader(400)
			return
		}
	}

	logItemList := logitem.Search(inputID, level, env, fromTS, toTS)

	byteArray, err := json.Marshal(logItemList)

	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(byteArray)
}
