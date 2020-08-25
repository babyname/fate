package main

import (
	"errors"
	"fmt"
	"github.com/goextension/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
)

func cmdCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check the env does correct",
		Run: func(cmd *cobra.Command, args []string) {
			root := runtime.GOROOT()
			log.Infow("check", "GOROOT", runtime.GOROOT())
			if e := zoneCheck(root); e != nil {
				log.Fatalw("zoneinfo check failed", "error", e)
			}
			fmt.Println("check all done!!!")
		},
	}
	return cmd
}
func zoneCheck(root string) error {
	path := filepath.Join(root, "lib", "time")
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
