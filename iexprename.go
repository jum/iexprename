package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	root            *string = flag.String("root", "/Volumes/Bilder/Jens-Uwe/iPhoto_Export", "root dir for")
	phoshare_compat *bool   = flag.Bool("phoshare-compat", true, "suffix compatibilty with phoshare exports")
	monthMap                = map[time.Month]string{
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
				if err == io.EOF {
					continue
				}
				panic(err.Error())
			}
			f.Close()
			debug("x %#v\n", x)
			dateTime, err := x.DateTime()
			if err != nil {
				panic(err.Error())
			}
			debug("dt %v", dateTime.Format(time.RFC3339))
			newDir := filepath.Join(*root, strconv.Itoa(dateTime.Year()), monthMap[dateTime.Month()])
			newName := filepath.Join(newDir, e.Name())
			debug("oldName %v, newName %v\n", subPath, newName)
			if err = os.MkdirAll(newDir, 0755); err != nil {
				panic(err.Error())
			}
			if err = os.Rename(subPath, newName); err != nil {
				panic(err.Error())
			}
			if *phoshare_compat {
				fd, err := os.Open(newName)
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
					if e.IsDir() {
						continue
					}
					fileName := e.Name()
					for i := len(fileName) - 1; i >= 0 && !os.IsPathSeparator(fileName[i]); i-- {
						if fileName[i] == '.' {
							fileName = fileName[:i] + strings.ToLower(fileName[i:])
							break
						}
					}
					if e.Name() != fileName {
						origPath := filepath.Join(newName, e.Name())
						newPath := filepath.Join(newName, fileName)
						debug("origPath %v, newPath %v\n", origPath, newPath)
						if err = os.Rename(origPath, newPath); err != nil {
							panic(err.Error())
						}
					}
				}
			}
			break
		}
	}
}
