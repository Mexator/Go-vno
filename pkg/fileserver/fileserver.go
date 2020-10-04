package main

import (
	"fmt"
	"os"

	"../config"
)

// FileServer implements FileServerServer
type FileServer struct {
	storageDir *os.File
}

type fileServerConfig struct {
	FilesDir string `json:"files_dir"`
}

/*
MakeFileServer creates FileServer reading its configuration from
`config.json` file.
*/
func MakeFileServer() (FileServer, error) {
	conf := new(fileServerConfig)
	err := config.ReadConfig(&conf, "config.json")
	if err == nil {
		file, err := os.Open(conf.FilesDir)
		if err == nil {
			return FileServer{file}, nil
		}
	}
	return FileServer{}, err
}

func main() {
	fmt.Println(MakeFileServer())
}
