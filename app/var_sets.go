package app

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
	"strings"
)

// VariableSets describes all the Variable Set related methods that the
// Terraform Enterprise API supports.
// TFE API docs: https://www.terraform.io/cloud-docs/api-docs/variable-sets
//
// List all the variable sets within an organization.
//	List(ctx context.Context, organization string, options *VariableSetListOptions) (*VariableSetList, error)
//
//	// ListForWorkspace gets the associated variable sets for a workspace.
//	ListForWorkspace(ctx context.Context, workspaceID string, options *VariableSetListOptions) (*VariableSetList, error)
//
//	// Create is used to create a new variable set.
//	Create(ctx context.Context, organization string, options *VariableSetCreateOptions) (*VariableSet, error)
//
//	// Read a variable set by its ID.
//	Read(ctx context.Context, variableSetID string, options *VariableSetReadOptions) (*VariableSet, error)
//
//	// Update an existing variable set.
//	Update(ctx context.Context, variableSetID string, options *VariableSetUpdateOptions) (*VariableSet, error)
//
//	// Delete a variable set by ID.
//	Delete(ctx context.Context, variableSetID string) error
//
//	// Apply variable set to workspaces in the supplied list.
//	ApplyToWorkspaces(ctx context.Context, variableSetID string, options *VariableSetApplyToWorkspacesOptions) error
//
//	// Remove variable set from workspaces in the supplied list.
//	RemoveFromWorkspaces(ctx context.Context, variableSetID string, options *VariableSetRemoveFromWorkspacesOptions) error
//
//	// Update list of workspaces to which the variable set is applied to match the supplied list.
//	UpdateWorkspaces(ctx context.Context, variableSetID string, options *VariableSetUpdateWorkspacesOptions) (*VariableSet, error)

var varSetIncludeOpts = map[string]tfe.VariableSetIncludeOpt{
	"workspaces": tfe.VariableSetWorkspaces,
	"vars":       tfe.VariableSetVars,
}

func (tfc *TFCClient) VarSetsListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Aliases:  []string{"ls"},
		Usage:    "List all the variable sets within an organization.",
		Category: "variable-sets",
		Action:   tfc.varSetsList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "include",
				Usage:   "A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api",
				Aliases: []string{"i"},
			},
			&cli.IntFlag{
				Name:  "page-num",
				Usage: "The page number to request. The results vary based on the PageSize.",
			},
			&cli.IntFlag{
				Name:  "page-size",
				Usage: "The number of elements returned in a single page.",
			},
		},
	}
}

func (tfc *TFCClient) varSetsList(ctx *cli.Context) error {
	opts := &tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: ctx.Int("page-num"),
			PageSize:   ctx.Int("page-size"),
		},
	}

	include := ctx.StringSlice("include")

	for _, r := range include {
		if opt, ok := varSetIncludeOpts[r]; !ok {
			return fmt.Errorf("include opt not recognized: %s", r)
		} else {
			if opts.Include == "" {
				opts.Include = string(opt)
				continue
			}
			opts.Include = strings.Join([]string{opts.Include, string(opt)}, ",")
		}
	}

	vsl, err := tfc.Client.VariableSets.List(ctx.Context, tfc.Cfg.OrgName, opts)
	if err != nil {
		return err
	}

	r, err := json.MarshalIndent(vsl, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))

	return nil
}
