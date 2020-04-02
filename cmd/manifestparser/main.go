package main

import (
	"github.com/smm-goddess/manifestparser/amxl"
	"io/ioutil"
)

var manifestPath = "/home/neal/work/apks/AndroidManifest.xml"

func main() {
	bs, _ := ioutil.ReadFile(manifestPath)
	amxl.Parse(bs)
}
