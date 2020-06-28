package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/godcong/fate/config"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
)

func cmdInit() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "output the init config",
		Run: func(cmd *cobra.Command, args []string) {
			absPath, e := filepath.Abs(path)
			if e != nil {
				log.Fatalw("wrong path", "error", e, "path", path)
			}
			fmt.Printf("config will output to %s\n", filepath.Join(absPath, config.JSONName))
			config.DefaultJSONPath = path

			e = config.OutputConfig(config.DefaultConfig())
			if e != nil {
				log.Fatalw("config wrong", "error", e)
			}

			e = zoneCheck()
			if e != nil {
				log.Fatalw("zoneinfo fix failed", "error", e)
			}
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "set the output path")
	return cmd
}

func zoneCheck() error {
	fmt.Println("GOROOT:", runtime.GOROOT())
	path := runtime.GOROOT() + "/lib/time"
	info, e := os.Stat(path)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.MkdirAll(path, 0755)
			if e != nil {
				return fmt.Errorf("could not make dir for copy zoneinfo:%w", e)
			}
			//fix:https://github.com/godcong/fate/issues/68
			info, e = os.Stat(path)
			if e != nil {
				return fmt.Errorf("error after make dir:%v", e)
			}
		} else {
			return fmt.Errorf("the target file is exist(%s):%w", path, e)
		}
	}
	if !info.IsDir() {
		return errors.New("destination file is not a directory")
	}

	filename := "zoneinfo.zip"

	src, e := os.Open(filename)
	if e != nil {
		return e
	}
	target := filepath.Join(path, filename)
	fmt.Println("copy zoneinfo to:", target)
	dst, e := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0755)
	if e != nil {
		return e
	}
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
