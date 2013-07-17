package web

import (
	hosieweb "github.com/hoisie/web"
	"io/ioutil"
)

func PostInput(ctx *hosieweb.Context) {
	data, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.WriteHeader(400)
		return
	}

}
