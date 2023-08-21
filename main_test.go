package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func Test_compress(t *testing.T) {
	dirName := "test"
	os.Mkdir(dirName, os.ModePerm)

	fileName := "test.txt"
	fPath := filepath.Join(dirName, fileName)
	log.Println(fPath)
	makeTextFile(t, fPath)
	dirName2 := filepath.Join(dirName, "child")
	os.Mkdir(dirName2, os.ModePerm)
	c := Config{
		OutputName: "ziptest",
		Files: []DocumentInfomation{
			DocumentInfomation{InputName: fPath},
			DocumentInfomation{InputName: dirName2},
		},
	}
	compress(&c, "ver")
}

func Test_compress2(t *testing.T) {
	dirName := "../Zipper/test"
	os.Mkdir(dirName, os.ModePerm)

	fileName := "test.txt"
	fPath := filepath.Join(dirName, fileName)
	log.Println(fPath)
	makeTextFile(t, fPath)
	dirName2 := filepath.Join(dirName, "child")
	os.Mkdir(dirName2, os.ModePerm)
	c := Config{
		OutputName: "ziptest2",
		Files: []DocumentInfomation{
			DocumentInfomation{InputName: fPath},
			DocumentInfomation{InputName: dirName2},
		},
	}
	compress(&c, "ver")
}

func makeTextFile(t *testing.T, fPath string) {
	fp, err := os.Create(fPath)
	defer fp.Close()
	fp.WriteString("hoge")
	if err != nil {
		t.Log("oi")
		t.Fatal("file does not exist:" + fPath)
	}
}
