package app

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

var WSIncludeOpts = map[string]tfe.WSIncludeOpt{
	"organization":                                         tfe.WSOrganization,
	"current_configuration_version":                        tfe.WSCurrentConfigVer,
	"current_configuration_version.ingress_attributes":     tfe.WSCurrentConfigVerIngress,
	"current_run":                                          tfe.WSCurrentRun,
	"current_run.plan":                                     tfe.WSCurrentRunPlan,
	"current_run.configuration_version":                    tfe.WSCurrentRunConfigVer,
	"current_run.configuration_version.ingress_attributes": tfe.WSCurrentrunConfigVerIngress,
	"locked_by":                                            tfe.WSLockedBy,
	"readme":                                               tfe.WSReadme,
	"outputs":                                              tfe.WSOutputs,
	"current-state-version":                                tfe.WSCurrentStateVer,
}

func (tfc *TFCClient) WorkspaceListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Aliases:  []string{"ls"},
		Usage:    "List all the workspaces within an organization.",
		Category: "workspace",
		Action:   tfc.workspacesList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "search",
				Usage:   "A search string (partial workspace name) used to filter the results.",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:    "tags",
				Usage:   "A search string (comma-separated tag names) used to filter the results.",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    "exclude-tags",
				Usage:   "A search string (comma-separated tag names to exclude) used to filter the results.",
				Aliases: []string{"et"},
			},
			&cli.StringFlag{
				Name:    "wildcard-name",
				Usage:   "A search on substring matching to filter the results.",
				Aliases: []string{"wn"},
			},
			&cli.StringFlag{
				Name:    "include",
				Usage:   "A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api/workspaces.html#available-related-resources",
				Aliases: []string{"i"},
			},
		},
	}
}

func (tfc *TFCClient) workspacesList(ctx *cli.Context) error {
	opts := &tfe.WorkspaceListOptions{
		ListOptions:  tfe.ListOptions{},
		Search:       ctx.String("search"),
		Tags:         ctx.String("tags"),
		ExcludeTags:  ctx.String("exclude-tags"),
		WildcardName: ctx.String("wildcard-name"),
	}

	include := ctx.StringSlice("include")

	for _, r := range include {
		if opt, ok := WSIncludeOpts[r]; !ok {
			return fmt.Errorf("include opt not recognized: %s", r)
		} else {
			if opts.Include == nil {
				opts.Include = []tfe.WSIncludeOpt{}
			}
			opts.Include = append(opts.Include, opt)
		}
	}

	wl, err := tfc.Client.Workspaces.List(ctx.Context, tfc.Cfg.OrgName, opts)
	if err != nil {
		return err
	}

	type listResponse struct {
		ID   string
		Name string
		Tags []string `json:",omitempty"`
	}

	response := make([]listResponse, len(wl.Items))

	for i := range wl.Items {
		response[i] = listResponse{
			ID:   wl.Items[i].ID,
			Name: wl.Items[i].Name,
			Tags: wl.Items[i].TagNames,
		}
	}

	r, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))

	return nil
}
