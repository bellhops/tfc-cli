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
		Usage:    "Update variable set variable.",
		Category: "variable-set variables",
		Action:   tfc.varSetVariableUpdate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "var-set-id",
				Usage:    "id of the variable set containing the variable to modify. See tfc-client var-sets to query variable sets.",
				Aliases:  []string{"set"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "var-id",
				Usage:    "id of the variable to modify.",
				Aliases:  []string{"id"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Usage:    "The name of the variable.",
				Aliases:  []string{"k"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "value",
				Usage:    "The name of the variable.",
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

	vsv, err := tfc.Client.VariableSetVariables.Update(ctx.Context, ctx.String("var-set-id"), ctx.String("var-id"), opts)
	if err != nil {
		return err
	}

	fmt.Println("updated variable set variable: ")
	r, err := json.MarshalIndent(vsv, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))
	return nil
}
