package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

type DocumentInfomation struct {
	InputName  string
	OutputName string
}

type Config struct {
	OutputName  string
	ProjectFile string
	Files       []DocumentInfomation
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

func CreateFileNAme(dInfo DocumentInfomation) string {

	outName := dInfo.OutputName
	if outName == "" {
		outName = dInfo.InputName
	}
	outName = strings.Replace(outName, "../", "", -1)
	outName = strings.Replace(outName, "..\\", "", -1)
	return outName
}

func writeZip(w *zip.Writer, done chan struct{}, ch chan FInfo) {
Loop:
	for {
		select {
		case v := <-ch:
			fmt.Println(v.Doc.InputName)
			hdr, _ := zip.FileInfoHeader(v.Info)
			hdr.Name = CreateFileNAme(v.Doc)
			f, err := w.CreateHeader(hdr)
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadFile(v.Doc.InputName)
			if err != nil {
				panic(err.Error() + ":" + v.Doc.InputName)
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

	for _, docInfo := range c.Files {
		ch := make(chan FInfo)
		done := make(chan struct{})
		go func() {
			inputFiles(docInfo, ch)
			close(done)
		}()
		writeZip(w, done, ch)
	}

	w.Close()
	return
}

type FInfo struct {
	Doc  DocumentInfomation
	Info os.FileInfo
}

func inputFiles(docInfo DocumentInfomation, ch chan<- FInfo) {
	info, err := os.Stat(docInfo.InputName)
	if err != nil {
		s := strings.Join([]string{err.Error(), docInfo.InputName}, ",")
		panic(s)
	}
	if info.IsDir() {
		files, err := ioutil.ReadDir(docInfo.InputName)
		if err != nil {
			s := strings.Join([]string{err.Error(), ":", docInfo.InputName}, ":")
			panic(s)
		}
		for _, v := range files {
			iPath := filepath.Join(docInfo.InputName, v.Name())
			oPath := filepath.Join(docInfo.OutputName, v.Name())
			inputFiles(DocumentInfomation{InputName: iPath, OutputName: oPath}, ch)
		}
	} else {
		ch <- FInfo{Doc: docInfo, Info: info}
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
