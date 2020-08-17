package main

import (
	"errors"
	"fmt"
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

		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "set the output path")
	return cmd
}

func zoneCheck() error {
	log.Info("GOROOT:", runtime.GOROOT())
	path := runtime.GOROOT() + "/lib/time"
	_, e := os.Stat(filepath.Join(path, "zoneinfo.zip"))
	if e != nil && os.IsNotExist(e) {
		_, e1 := os.Stat(filepath.Join(getCurrentPath(), "zoneinfo.zip"))
		if e1 != nil && os.IsNotExist(e) {
			return errors.New("zoneinfo file not found")
		} else if e1 != nil {
			return fmt.Errorf("found error in current path:%w", e1)
		} else {
			return nil
		}
	} else if e != nil {
		return fmt.Errorf("found error in go path:%w", e)
	} else {
		return nil
	}
}
