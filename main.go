// logbuzz project main.go
package main

import (
	hosieweb "github.com/hoisie/web"
	myweb "logbuzz/web"
)

func main() {
	hosieweb.Post("/api/input", myweb.PostInput)
	hosieweb.Run("0.0.0.0:1212")
}
