package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	root     *string = flag.String("root", "/Volumes/Bilder/Jens-Uwe/iPhoto_Export", "root dir for")
	monthMap         = map[time.Month]string{
		time.January:   "01 Januar",
		time.February:  "02 Februar",
		time.March:     "03 März",
		time.April:     "04 April",
		time.May:       "05 Mai",
		time.June:      "06 Juni",
		time.July:      "07 Juli",
		time.August:    "08 August",
		time.September: "09 September",
		time.October:   "10 Oktober",
		time.November:  "11 November",
		time.December:  "12 Dezember",
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
	fi, err := fd.Readdir(-1)
	if err != nil {
		panic(err.Error())
	}
	fd.Close()
	//debug("fi = %#v\n", fi)
	for _, e := range fi {
		if !e.IsDir() {
			continue
		}
		//debug("subdir %#v\n", e)
		subPath := filepath.Join(*root, e.Name())
		subfd, err := os.Open(subPath)
		if err != nil {
			panic(err.Error())
		}
		subDir, err := subfd.Readdir(-1)
		subfd.Close()
		if err != nil {
			panic(err.Error())
		}
		for _, se := range subDir {
			//debug("ssubdir %#v\n", se)
			if se.IsDir() {
				continue
			}
			ext := filepath.Ext(se.Name())
			//debug("ext %v\n", ext)
			if !(strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg")) {
				continue
			}
			fname := filepath.Join(subPath, se.Name())
			debug("doing %v\n", fname)
			f, err := os.Open(fname)
			if err != nil {
				panic(err.Error())
			}
			x, err := exif.Decode(f)
			if err != nil {
				panic(err.Error())
			}
			f.Close()
			debug("x %#v\n", x)
			dt, err := x.Get(exif.DateTime)
			if err != nil {
				if _, ok := err.(exif.TagNotPresentError); ok {
					continue
				}
				panic(err.Error())
			}
			dateTime, err := time.Parse("2006:01:02 15:04:05", dt.StringVal())
			if err != nil {
				panic(err.Error())
			}
			debug("dt %#v, %v", dt, dateTime.Format(time.RFC3339))
			newDir := filepath.Join(*root, strconv.Itoa(dateTime.Year()), monthMap[dateTime.Month()])
			newName := filepath.Join(newDir, e.Name())
			debug("oldName %v, newName %v\n", subPath, newName)
			if err = os.MkdirAll(newDir, 0755); err != nil {
				panic(err.Error())
			}
			if err = os.Rename(subPath, newName); err != nil {
				panic(err.Error())
			}
			break
		}
	}
}
