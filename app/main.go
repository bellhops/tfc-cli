package app

import (
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

func CreateTFCClient(app *cli.App, cfg *Config, c *tfe.Client) *TFCClient {
	//if app == nil {
	//	app = cli.NewApp()
	//}

	return &TFCClient{
		//App:    app,
		Client: c,
		Cfg:    cfg,
	}
}

type TFCClient struct {
	//App    *cli.App
	Client *tfe.Client
	Cfg    *Config
}
