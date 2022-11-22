package app

import (
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

func CreateCtrl(app *cli.App, cfg *Config, c *tfe.Client) *Ctrl {
	if app == nil {
		app = cli.NewApp()
	}

	return &Ctrl{
		App:    app,
		Client: c,
		Cfg:    cfg,
	}
}

type Ctrl struct {
	App    *cli.App
	Client *tfe.Client
	Cfg    *Config
}
