package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/BurntSushi/toml"
)

func main() {
	conf := initConfig()
	fmt.Println(conf)
	zf, _ := os.Create(conf.OutputName)
	defer zf.Close()
	compress(zf, conf.Files)
}

type Config struct {
	OutputName string
	Files      []string
}

func initConfig() *Config {
	c := Config{}
	toml.DecodeFile("./config.toml", &c)
	return &c

}

func compress(buf io.Writer, files []string) {
	w := zip.NewWriter(buf)

	for _, file := range files {
		info, _ := os.Stat(file)
		val, _ := info.Sys().(*syscall.Stat_t)
		fmt.Println(val)
		fmt.Println(info.Name())

		hdr, _ := zip.FileInfoHeader(info)
		hdr.Name = "files/" + file
		f, err := w.CreateHeader(hdr)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		f.Write(body)
	}

	w.Close()

	return
}

func save(b *bytes.Buffer) error {
	zf, err := os.Create("sample.zip")
	if err != nil {
		return err
	}
	zf.Write(b.Bytes())
	zf.Close()
	return nil
}
