package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type vars struct {
	fileCount   int
	hashCount   int
	hashList    []string
	hashListAll []string
	home        string
	num         string
	mem         map[string]int
	newPath     string
	myPath      string
	sep         string
}

func (v *vars) makeNewDir() {

	if v.myPath != "" {
		return
	}

	v.refactorPath()

	v.myPath = strings.TrimSuffix(v.myPath, v.sep)
	arr := strings.Split(v.myPath, v.sep)
	v.newPath = fmt.Sprintf("%s/new-%s", v.myPath, arr[len(arr)-1])

	os.Mkdir(v.newPath, 0755)
}

func (v *vars) refactorPath() string {

	if len(os.Args) > 1 {
		v.myPath = os.Args[1]
	} else {
		fmt.Print("Choose a directory to look through: ")
		fmt.Scanf("%s", &v.myPath)
	}

	v.myPath = strings.Replace(v.myPath, "~", v.home, -1)

	return v.myPath
}

func (v *vars) getMD5Checksum(filePath string, fileName string) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(file)

	if err != nil {
		log.Fatal(err)
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	v.saveHash(fileName, file, checksum)
}

func (v *vars) saveHash(fileName string, file []byte, md5Sum string) {
	value, isThere := v.mem[fileName]
	look4Hash := v.look4Hashes(fileName, md5Sum)

	if look4Hash == 1 {
		v.hashCount += 1

		if isThere {
			v.mem[fileName] += 1
			v.num = fmt.Sprintf("%d", value+1)
		} else {
			v.num = ""
		}

		v.hashList = append(v.hashList, fileName, md5Sum)
		arr := strings.Split(fileName, ".")
		filePath := fmt.Sprintf("%s%s%s(%s).%s", v.newPath, v.sep, arr[0], v.num, arr[len(arr)-1])
		err := ioutil.WriteFile(filePath, file, 0644)

		if err != nil {
			log.Fatal(err)
		}

	} else if look4Hash == 2 {
		v.hashCount += 1
		v.mem[fileName] = 0

		v.hashList = append(v.hashList, fileName, md5Sum)
		filePath := fmt.Sprintf("%s%s%s", v.newPath, v.sep, fileName)
		err := ioutil.WriteFile(filePath, file, 0644)

		if err != nil {
			log.Fatal(err)
		}

	}

	v.hashListAll = append(v.hashListAll, fileName, md5Sum)

}

func (v *vars) look4Hashes(fileName string, md5Sum string) int {

	for _, element := range v.hashList {

		if md5Sum == element {
			return 0
		}

	}

	for _, element := range v.hashList {

		if fileName == element {
			return 1
		}

	}

	return 2
}

func (v *vars) look4Files() {
	arr := strings.Split(v.newPath, v.sep)
	newPathDir := arr[len(arr)-1]

	err := filepath.Walk(v.myPath, func(path string, info os.FileInfo, err error) error {
		arr = strings.Split(path, v.sep)
		checkDir := arr[len(arr)-2]

		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() && checkDir != newPathDir {
			v.fileCount += 1
			v.getMD5Checksum(path, info.Name())
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if v.verifyFiles() {
		fmt.Print("\n[+] SUCCESS\n")
		fmt.Print("[+] ALL FILES HAVE BEEN COPIED\n")
		fmt.Printf("\n>>> Old Files: %d\n", v.fileCount)
		fmt.Printf(">>> New Files: %d\n\n", v.hashCount)
	} else {
		fmt.Print("\n[-] FAIL\n")
		fmt.Print("[-] SOMETHING WENT WRONG\n\n")
	}

}

func (v *vars) verifyFiles() bool {
	var integrity bool

start:

	for i := 1; i < len(v.hashListAll); i += 2 {

		for j := 1; j < len(v.hashList); j += 2 {

			if v.hashListAll[i] == v.hashList[j] {
				integrity = true
				continue start
			} else {
				integrity = false
			}

		}

	}

	return integrity
}

func main() {
	home, _ := os.UserHomeDir()

	iv := vars{
		fileCount:   0,
		hashCount:   0,
		hashList:    []string{},
		hashListAll: []string{},
		home:        home,
		num:         "",
		mem:         map[string]int{},
		newPath:     "",
		myPath:      "",
		sep:         fmt.Sprintf("%c", os.PathSeparator),
	}

	iv.makeNewDir()
	iv.look4Files()
}
