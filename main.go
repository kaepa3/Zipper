package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func writeZip(w *zip.Writer, done chan struct{}, ch chan FInfo) {
Loop:
	for {
		select {
		case v := <-ch:
			fmt.Println(v.Path)
			hdr, _ := zip.FileInfoHeader(v.Info)
			hdr.Name = "files/" + v.Path
			f, err := w.CreateHeader(hdr)
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadFile(v.Path)
			if err != nil {
				panic(err.Error() + ":" + v.Path)
			}
			f.Write(body)
			break
		case <-done:
			break Loop
		}
	}
}

func compress(c *Config, ver string) {
	filepath := c.OutputName + "_" + ver + ".zip"
	zf, _ := os.Create(filepath)
	defer zf.Close()
	w := zip.NewWriter(zf)

	for _, file := range c.Files {
		ch := make(chan FInfo)
		done := make(chan struct{})
		go func() {
			inputFiles(file, ch)
			close(done)
		}()
		writeZip(w, done, ch)
	}

	w.Close()
	return
}

type FInfo struct {
	Path string
	Info os.FileInfo
}

func inputFiles(path string, ch chan<- FInfo) {
	info, err := os.Stat(path)
	if err != nil {
		s := strings.Join([]string{err.Error(), path}, ",")
		panic(s)
	}
	if info.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			s := strings.Join([]string{err.Error(), ":", path}, ":")
			panic(s)
		}
		for _, v := range files {
			cPath := filepath.Join(path, v.Name())
			log.Println(cPath + ":start")
			inputFiles(cPath, ch)
		}
	} else {
		ch <- FInfo{Path: path, Info: info}
	}
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
