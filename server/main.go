package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	if ginMode, _ := os.LookupEnv("GIN_MODE"); ginMode == "release" {

	} else {
		parcel := exec.Command("parcel", "index.html")
		parcel.Stdout = os.Stdout
		parcel.Stderr = os.Stderr
		parcel.Dir, _ = filepath.Abs("ui/")
		parcel.Start()
		defer parcel.Process.Kill()
		r.Use(proxy("127.0.0.1:1234"))
	}

	r.Run()
}
