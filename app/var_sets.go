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
			&cli.StringSliceFlag{
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
			&cli.StringFlag{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "exact name to filter results by",
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

	r := vsl.Items

	if ctx.IsSet("search") {
		r := []*tfe.VariableSet{}

		for _, vs := range vsl.Items {
			if vs.Name == ctx.String("search") {
				r = append(r, vs)
			}
		}
	}

	pp, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(pp))

	return nil
}

func (tfc *TFCClient) VarSetsListForWorkspaceCmd() *cli.Command {
	return &cli.Command{
		Name:     "list-for-workspace",
		Aliases:  []string{"ws-ls"},
		Usage:    "List all the variable sets within a workspace.",
		Category: "variable-sets",
		Action:   tfc.varSetsListForWorkspace,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "workspace-id",
				Usage:   "ID of workspace to query.",
				Aliases: []string{"id", "ws"},
			},
			&cli.StringSliceFlag{
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

type varSetResponse struct {
	ID          string
	Name        string
	Description string
	Global      bool
	OrgName     string
	Workspaces  []*tfe.Workspace           `json:",omitempty"`
	Variables   []*tfe.VariableSetVariable `json:",omitempty"`
}

func (tfc *TFCClient) varSetsListForWorkspace(ctx *cli.Context) error {
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

	vsl, err := tfc.Client.VariableSets.ListForWorkspace(ctx.Context, ctx.String("workspace-id"), opts)
	if err != nil {
		return err
	}

	items := make([]varSetResponse, len(vsl.Items))

	for i := range items {
		raw := vsl.Items[i]
		items[i] = varSetResponse{
			ID:          raw.ID,
			Name:        raw.Name,
			Description: raw.Description,
			Global:      raw.Global,
			OrgName:     raw.Organization.Name,
		}

		if strings.Contains(opts.Include, string(tfe.VariableSetWorkspaces)) {
			items[i].Workspaces = raw.Workspaces
		}

		if strings.Contains(opts.Include, string(tfe.VariableSetVars)) {
			items[i].Variables = raw.Variables
		}
	}

	r, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))

	return nil
}

func (tfc *TFCClient) VarSetsReadCmd() *cli.Command {
	return &cli.Command{
		Name:     "get",
		Usage:    "get variable set by id or name",
		Category: "variable-sets",
		Action:   tfc.varSetsRead,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "include",
				Usage:   "A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:  "id",
				Usage: "id to query",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "name to query",
			},
		},
	}
}
func (tfc *TFCClient) varSetsRead(ctx *cli.Context) error {
	if !ctx.IsSet("id") && !ctx.IsSet("name") {
		return fmt.Errorf("one of --id or --name is required")
	}

	if ctx.IsSet("id") && ctx.IsSet("name") {
		return fmt.Errorf("only one of --id or --name can be passed")
	}

	opts := &tfe.VariableSetListOptions{}

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

	r := vsl.Items

	if ctx.IsSet("id") {
		r = []*tfe.VariableSet{}

		for _, vs := range vsl.Items {
			if vs.ID == ctx.String("id") {
				r = append(r, vs)
			}
		}
	}

	if ctx.IsSet("name") {
		r = []*tfe.VariableSet{}

		for _, vs := range vsl.Items {
			if vs.Name == ctx.String("name") {
				r = append(r, vs)
			}
		}
	}

	pp, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(pp))

	return nil
}
