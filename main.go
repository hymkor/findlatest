package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

var (
	flagAll        = flag.Bool("a", false, "Include dot files")
	flagVerboseDir = flag.Bool("vd", false, "Display the name of the directory currently being processed.")
	flagQuiet      = flag.Bool("q", false, "Be quiet")
	flagUntil      = flag.String("until", "2999-01-02 15:04:05", "")
)

type Latest struct {
	Path  string
	Stamp time.Time
	Until time.Time
	All   bool
}

func checkDir(path1 string, latest *Latest) error {
	if *flagVerboseDir {
		fmt.Println(path1)
	}
	entries, err := os.ReadDir(path1)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		full := filepath.Join(path1, name)
		if !latest.All && len(name) > 0 && name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			if name == "." || name == ".." {
				continue
			}
			if err := checkDir(full, latest); err != nil {
				return err
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		stamp := info.ModTime()
		if stamp.After(latest.Stamp) {
			if stamp.After(latest.Until) {
				continue
			}
			if !*flagQuiet {
				fmt.Println(stamp.Format(dateFormat), full)
			}
			latest.Stamp = stamp
			latest.Path = full
		}
	}
	return nil
}

func check(path1 string, latest *Latest) error {
	stat, err := os.Stat(path1)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return checkDir(path1, latest)
	}
	stamp := stat.ModTime()
	if stamp.After(latest.Stamp) {
		fmt.Println(stamp.Format(dateFormat), path1)
		latest.Stamp = stamp
		latest.Path = path1
	}
	return nil
}

func mains(args []string) error {
	latest := &Latest{
		Stamp: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if *flagQuiet {
		defer func() {
			fmt.Println(latest.Stamp.Format(dateFormat), latest.Path)
		}()
	}
	if stat, err := os.Stat(*flagUntil); err == nil {
		latest.Until = stat.ModTime()
	} else if t, err := time.Parse("2006-01-02 15:04:05", *flagUntil); err == nil {
		latest.Until = t
	} else {
		return err
	}
	latest.All = *flagAll
	if len(args) <= 0 {
		return checkDir(".", latest)
	}
	for _, arg := range args {
		if matches, err := filepath.Glob(arg); err == nil {
			for _, path1 := range matches {
				if err := check(path1, latest); err != nil {
					return err
				}
			}
		} else {
			if err := check(arg, latest); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
