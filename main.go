package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	conf := initConfig()
	fmt.Println(conf)
	version, _ := getVersion(conf)
	fmt.Println(version)
	compress(conf, version)
}

type Config struct {
	OutputName  string
	ProjectFile string
	Files       []string
}

func initConfig() *Config {
	c := Config{}
	toml.DecodeFile("./config.toml", &c)
	return &c

}

type Properties struct {
	Version string
}
type Project struct {
	PropertyGroup Properties
}

func getVersion(c *Config) (string, error) {
	var origin Project
	data, _ := ioutil.ReadFile(c.ProjectFile)
	err := xml.Unmarshal(data, &origin)
	if err != nil {
		return "", err
	}
	return origin.PropertyGroup.Version, nil

}

func compress(c *Config, ver string) {
	filepath := c.OutputName + "_" + ver + ".zip"
	zf, _ := os.Create(filepath)
	defer zf.Close()
	w := zip.NewWriter(zf)

	for _, file := range c.Files {
		info, err := os.Stat(file)
		if err != nil {
			panic(err)
		}
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
