package app

import "github.com/hashicorp/go-tfe"

type Config struct {
	TFE     *tfe.Config
	OrgName string
}
