package main

import (
	"flag"
	"github.com/cnorman/tfc-cli/app"
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

/*
todo:
1. list/search workspaces
2. list/search workspace variables
3. Add/Update Workspace Variables
4. List runs?
5. Create run (may have some configuration dependencies)

Objects:
- Workspace ( https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces )
- Variable-Set ( https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets )
- Configuration Versions ( https://developer.hashicorp.com/terraform/cloud-docs/api-docs/configuration-versions )
- Runs ( https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run )
*/

func main() {
	var t, o string
	flag.StringVar(&t, "token", os.Getenv("TFC_TOKEN"), "terraform cloud api token")
	flag.StringVar(&o, "org", os.Getenv("TFC_ORG"), "terraform cloud org")
	flag.Parse()

	if t == "" {
		log.Fatal("missing access token. Set TFC_TOKEN in the environment or pass --token")
	}

	if o == "" {
		log.Fatal("missing terraform cloud organization. Set TFC_ORG in the environment or pass --org")
	}

	// Get org and token from env or flag
	tfeCFG := &tfe.Config{
		Token: t,
	}

	tfeClient, err := tfe.NewClient(tfeCFG)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &app.Config{
		TFE:     tfeCFG,
		OrgName: o,
	}

	a := cli.NewApp()
	a.Name = "tfc-cli"
	a.Usage = "Interact with Terraform Cloud via CLI"
	a.UsageText = "Interact with Terraform Cloud via CLI\nRequires Terraform Cloud API Access Token. " +
		"Set TFC_TOKEN or pass --token\nRequires Terraform Cloud Organization Name. Set TFC_ORG or pass --org" +
		"\nSee https://developer.hashicorp.com/terraform/cloud-docs/api-docs for reference"

	ctrl := app.CreateCtrl(a, cfg, tfeClient)

	ctrl.InitWorkspaces()

	err = ctrl.App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
