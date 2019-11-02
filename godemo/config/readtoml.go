package godemo

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

func testtoml() {
	var f *os.File
	var err error
	if f, err = os.Open("gDNS.toml"); err != nil {
		log.Fatal("Open file failed,", err)
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(f); err != nil {
		log.Fatal("read from file failed,", err)
	}

	tree, err := toml.LoadReader(buf)
	if err != nil {
		log.Fatal("cannot load file content,", err)
	}

	res := tree.ToMap()
	fmt.Println(res)
	fmt.Println(res["type"])
	fmt.Println(res["server"])
	fmt.Println(res["hosts"])
	fmt.Println(res["log"])
}
