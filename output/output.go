package output

import "github.com/godcong/fate/config"

type Output interface {
}

func NewOutputWithConfig(config config.Config) Output {
	return nil
}
