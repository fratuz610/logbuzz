package vars

import (
	"log"
	"os"
	"os/user"
)

const (
	persistFile string = "logbuzz-data.bin"
)

func GetPersistFile() string {
	return getUserHomeFolder() + string(os.PathSeparator) + persistFile
}

func getUserHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
