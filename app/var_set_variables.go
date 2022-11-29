package app

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

// VariableSetVariables describes all variable set variable related methods within the scope of
// Variable Sets that the Terraform Enterprise API supports
//
// TFE API docs: https://www.terraform.io/cloud-docs/api-docs/variable-sets#variable-relationships
//
//	// List all variables in the variable set.
//	List(ctx context.Context, variableSetID string, options *VariableSetVariableListOptions) (*VariableSetVariableList, error)
//
//	// Create is used to create a new variable within a given variable set
//	Create(ctx context.Context, variableSetID string, options *VariableSetVariableCreateOptions) (*VariableSetVariable, error)
//
//	// Read a variable by its ID
//	Read(ctx context.Context, variableSetID string, variableID string) (*VariableSetVariable, error)
//
//	// Update valuse of an existing variable
//	Update(ctx context.Context, variableSetID string, variableID string, options *VariableSetVariableUpdateOptions) (*VariableSetVariable, error)
//
//	// Delete a variable by its ID
//	Delete(ctx context.Context, variableSetID string, variableID string) error

func (tfc *TFCClient) VarSetVariablesListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Aliases:  []string{"ls"},
		Usage:    "List all variables in the variable set.",
		Category: "variable-set variables",
		Action: func(ctx *cli.Context) error {
			vsl, err := tfc.Client.VariableSetVariables.List(
				ctx.Context,
				ctx.String("var-set-id"), &tfe.VariableSetVariableListOptions{
					ListOptions: tfe.ListOptions{
						PageNumber: ctx.Int("page-num"),
						PageSize:   ctx.Int("page-size"),
					},
				},
			)
			if err != nil {
				return err
			}

			r, err := json.MarshalIndent(vsl, "", "    ")
			if err != nil {
				return nil
			}

			fmt.Println(string(r))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "var-set-id",
				Usage:    "id of the variable set to query. See tfc-client var-sets to query variable sets.",
				Aliases:  []string{"id"},
				Required: true,
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

func (tfc *TFCClient) VarSetVariablesUpdateCmd() *cli.Command {
	return &cli.Command{
		Name:     "update",
		Usage:    "Update variable set variable.\none of set-id, var-id is required",
		Category: "variable-set variables",
		Action:   tfc.varSetVariableUpdate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "set-id",
				Usage: "id of the variable set containing the variable to modify. See tfc-client var-sets to query variable sets.",
			},
			&cli.StringFlag{
				Name:  "set-name",
				Usage: "name of the variable set containing the variable to modify. See tfc-client var-sets to query variable sets. IGNORED IF var-set-id IS SET.",
			},
			&cli.StringFlag{
				Name:     "key",
				Aliases:  []string{"k"},
				Usage:    "(Required) key of the variable to modify.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "value",
				Usage:    "(Required) The value of the variable.",
				Aliases:  []string{"v"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "description",
				Usage:   "The description of the variable.",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:  "hcl",
				Usage: "Whether to evaluate the value of the variable as a string of HCL code.",
			},
			&cli.BoolFlag{
				Name:    "sensitive",
				Usage:   "Whether the value is sensitive.",
				Aliases: []string{"s"},
			},
			// TODO: accept json input
		},
	}
}

func (tfc *TFCClient) varSetVariableUpdate(ctx *cli.Context) error {
	if ctx.IsSet("set-id") && ctx.IsSet("set-name") {
		return fmt.Errorf("only one of \"--var-set-id\" or \"--var-set-name\" can be used")
	}

	var (
		varSet  *tfe.VariableSet
		varID   *string
		verbose = ctx.Bool("verbose")
	)

	// Fetch the variable set with all related variables
	if ctx.IsSet("set-id") {
		var err error
		opts := &tfe.VariableSetReadOptions{Include: &[]tfe.VariableSetIncludeOpt{tfe.VariableSetVars}}

		if verbose {
			fmt.Printf("recieved set-id: %s\nReading variable set with options: %+v\n", ctx.String("set-id"), *opts)
		}

		if varSet, err = tfc.Client.VariableSets.Read(ctx.Context, ctx.String("set-id"), opts); err != nil {
			return err
		}

		if verbose {
			fmt.Printf("successfully read variable set: %+v\n", varSet)
		}
	} else {
		// If we get a name and not an ID we need to query for the ID
		var p *tfe.Pagination
		for p == nil || p.CurrentPage < p.TotalPages {
			p = &tfe.Pagination{
				NextPage: 0,
			}
			opts := &tfe.VariableSetListOptions{ListOptions: tfe.ListOptions{
				PageNumber: p.NextPage,
				PageSize:   50,
			}, Include: string(tfe.VariableSetVars)}

			if verbose {
				fmt.Printf("listing variable sets with options: %+v\n", opts)
			}

			lr, err := tfc.Client.VariableSets.List(ctx.Context, tfc.Cfg.OrgName, opts)
			if err != nil {
				return err
			}

			fmt.Println("Total variable sets returned: ", lr.TotalCount)

			// look for matching set on this page
			fmt.Println("looking for matching set: ", ctx.String("set-name"))
			for _, vs := range lr.Items {
				if vs.Name == ctx.String("set-name") {
					if verbose {
						// OMG why did I write so much logging??
						var result struct {
							ID          string                    `json:"id"`
							Name        string                    `json:"name"`
							Description string                    `json:"description"`
							Global      bool                      `json:"global"`
							Variables   []tfe.VariableSetVariable `json:"vars,omitempty"`
						}

						result.ID = vs.ID
						result.Global = vs.Global
						result.Name = vs.Name
						result.Description = vs.Description
						result.Variables = make([]tfe.VariableSetVariable, len(vs.Variables))

						for i := range vs.Variables {
							result.Variables[i] = *vs.Variables[i]
						}

						fmt.Printf("Variable set name match: %+v\n", result)
					}
					varSet = vs
					break
				}

				if verbose {
					fmt.Printf("Variable set name doesn't match: %s\n", vs.Name)
				}
			}

			// If we found it, stop looking
			if varSet != nil {
				break
			}

			p = lr.Pagination
		}

		if varSet == nil {
			return fmt.Errorf("variable set with matching name not found")
		}
	}

	for _, v := range varSet.Variables {
		if v.Key == ctx.String("key") {
			if verbose {
				fmt.Printf("Variable key match: %+v\n", *v)
			}
			varID = ptrString(v.ID)
			if err := ctx.Set("key", v.Key); err != nil {
				return err
			}
		}
	}

	if varID == nil {
		return fmt.Errorf("matching variable not found\nkey: %s", ctx.String("key"))
	}

	opts := &tfe.VariableSetVariableUpdateOptions{
		Key:         ptrString(ctx.String("key")),
		Value:       ptrString(ctx.String("value")),
		Description: ptrString(ctx.String("description")),
	}

	if ctx.IsSet("hcl") {
		opts.HCL = ptrBool(ctx.Bool("hcl"))
	}

	if ctx.IsSet("sensitive") {
		opts.Sensitive = ptrBool(ctx.Bool("sensitive"))
	}

	if verbose {
		fmt.Printf("Updating Variable: %+v\n", struct {
			variableSetID, variableID string
			options                   tfe.VariableSetVariableUpdateOptions
		}{
			variableSetID: varSet.ID,
			variableID:    *varID,
			options:       *opts,
		})
	}

	vsv, err := tfc.Client.VariableSetVariables.Update(ctx.Context, varSet.ID, *varID, opts)
	if err != nil {
		return err
	}

	r, err := json.MarshalIndent(vsv, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Printf("updated variable set variable:\n%s\n", string(r))
	return nil
}
