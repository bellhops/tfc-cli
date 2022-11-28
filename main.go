package main

import (
	"flag"
	"github.com/bellhops/tfc-cli/app"
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

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

	tfc := app.CreateTFCClient(nil, cfg, tfeClient)

	a := cli.NewApp()
	a.Name = "tfc-cli"
	a.Usage = "tfc-client [resource] [action] [options]; tfc-client workspaces list --search=\"api\""
	a.UsageText = "Interact with Terraform Cloud via CLI\nRequires Terraform Cloud API Access Token. " +
		"Set TFC_TOKEN or pass --token\nRequires Terraform Cloud Organization Name. Set TFC_ORG or pass --org" +
		"\nSee https://developer.hashicorp.com/terraform/cloud-docs/api-docs for reference"

	// The App.Commands field contains the top level resource commands:
	// 		tfc-client [resource]; tfc-client workspaces
	// The Subcommands field for each resource command contain the action subcommands:
	//		tfc-client [resource] [action] [options]; tfc-client workspaces list --search="api"
	a.Commands = []*cli.Command{
		{
			Name:        "workspaces",
			Usage:       "Query Workspaces via cli options",
			UsageText:   "Query Workspaces via cli options\nReference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces",
			Subcommands: []*cli.Command{tfc.WorkspaceListCmd()},
		},
		{
			Name:        "config-versions",
			Usage:       "Query Terraform Workspace Configuration Versions",
			UsageText:   "Query Terraform Workspace Configuration Versions\nReference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/configuration-versions",
			Subcommands: []*cli.Command{tfc.ConfigVersionsListCmd()},
		},
		{
			Name:        "var-sets",
			Usage:       "Interact Terraform Variable Sets",
			UsageText:   "Interact Terraform Variable Sets\nReference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-sets",
			Subcommands: []*cli.Command{tfc.VarSetsListCmd(), tfc.VarSetsListForWorkspaceCmd(), tfc.VarSetsReadCmd()},
		},
		{
			Name:        "var-set-variables",
			Usage:       "Interact Terraform Variable Set Variables",
			UsageText:   "Interact Terraform Variable Set Variables\nReference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/variable-set-variables",
			Subcommands: []*cli.Command{tfc.VarSetVariablesListCmd(), tfc.VarSetVariablesUpdateCmd()},
		},
		{
			Name:        "runs",
			Usage:       "Interact with Terraform Cloud runs",
			UsageText:   "Interact with Terraform Cloud runs\nReference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/runs",
			Subcommands: []*cli.Command{tfc.RunsCreateCmd()},
		},
	}

	err = a.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
