package information

import "github.com/godcong/fate/config"

type Information interface {
	Write()
}

func NewWithConfig(config config.Config) Information {
	return nil
}
