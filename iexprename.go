package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	root     *string = flag.String("root", "/Users/jum/Desktop/Zum Platzieren", "root dir for")
	// XX2008-11XX30.11.2008
	nameRegex = regexp.MustCompile(`XX([\d]{4})-([\d]{2})XX(.*)`)
	monthMap = map[string]string{
		"01": "01 Januar",
		"02": "02 Februar",
		"03": "03 März",
		"04": "04 April",
		"05": "05 Mai",
		"06": "06 Juni",
		"07": "07 Juli",
		"08": "08 August",
		"09": "09 September",
		"10": "10 Oktober",
		"11": "11 November",
		"12": "12 Dezember",
	}
)

const DEBUG = false

func debug(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

func main() {
	flag.Parse()
	debug("dir %#v\n", *root)
	fd, err := os.Open(*root)
	if err != nil {
		panic(err.Error())
	}
	defer fd.Close()
	fi, err := fd.Readdir(-1)
	if err != nil {
		panic(err.Error())
	}
	//debug("fi = %#v\n", fi)
	for _, e := range fi {
		if !e.IsDir() {
			continue
		}
		if m := nameRegex.FindStringSubmatch(e.Name()); m != nil {
			oldName := filepath.Join(*root, e.Name())
			debug("old = %#v\n", oldName)
			//debug("m = %#v\n", m)
			newName := filepath.Join(*root, m[1], monthMap[m[2]], m[3])
			debug("new = %#v\n", newName)
			newDir := filepath.Join(*root, m[1])
			debug("dir = %#v\n", newDir)
			if err = os.Mkdir(newDir, 0775); err != nil {
				debug("mkdir %v: %v\n", newDir, err.Error())
			}
			newDir = filepath.Join(*root, m[1], monthMap[m[2]])
			debug("dir = %#v\n", newDir)
			if err = os.Mkdir(newDir, 0775); err != nil {
				debug("mkdir %v: %v\n", newDir, err.Error())
			}
			if err = os.Rename(oldName, newName); err != nil {
				panic(err.Error())
			}
		}
	}
}
