package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

func (tfc *TFCClient) RunsCreateCmd() *cli.Command {
	return &cli.Command{
		Name:     "create",
		Usage:    "Create run",
		Category: "runs",
		Action:   tfc.runCreate,
		Flags: []cli.Flag{
			// Required
			&cli.StringFlag{
				Name:     "workspace-id",
				Aliases:  []string{"workspace", "ws"},
				Usage:    "(Required) The workspace where the run will be executed.",
				Required: true,
			},
			// Optional
			&cli.BoolFlag{
				Name:  "allow-empty-apply",
				Usage: "Whether Terraform can apply the run even when the plan contains no changes. Often used to upgrade state after upgrading a workspace to a new terraform version.",
			},
			&cli.StringFlag{
				Name:  "terraform-version",
				Usage: "The Terraform version to use in this run. Only valid for plan-only runs; must be a valid Terraform version available to the organization.",
			},
			&cli.BoolFlag{
				Name:  "plan-only",
				Usage: "This is a speculative, plan-only run that Terraform cannot apply. Often used in conjunction with terraform-version in order to test whether an upgrade would succeed.",
			},
			&cli.BoolFlag{
				Name:  "is-destroy",
				Usage: "This plan is a destroy plan, which will destroy all provisioned resources.",
			},
			&cli.BoolFlag{
				Name:  "refresh",
				Usage: "The run should update the state prior to checking for differences",
			},
			&cli.BoolFlag{
				Name:  "refresh-only",
				Usage: "The run should ignore config changes and refresh the state only",
			},
			&cli.BoolFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "Message to be associated with this run.",
			},
			&cli.StringFlag{
				Name:    "configuration-version",
				Aliases: []string{"config-version"},
				Usage:   "The configuration version to use for this run. If the configuration version object is omitted, the run will be created using the workspace's latest configuration version.",
			},
			&cli.BoolFlag{
				Name:  "auto-apply",
				Usage: "The run should be applied automatically without user confirmation. It defaults to the Workspace.AutoApply setting.",
			},
			&cli.StringSliceFlag{
				Name:  "var",
				Usage: "key=value; Terraform input variables for a particular run, prioritized over variables defined on the workspace. All values must be expressed as an HCL literal in the same syntax you would use when writing terraform code.",
			},
		},
	}
}

type runCreateResponse struct {
	ID              string
	CreatedAt       time.Time
	AutoApply       bool
	HasChanges      bool
	Status          string
	PositionInQueue int
}

func (tfc *TFCClient) runCreate(ctx *cli.Context) error {
	opts := tfe.RunCreateOptions{
		AllowEmptyApply:  getIfSetBool(ctx, "allow-empty-apply"),
		TerraformVersion: getIfSetString(ctx, "terraform-version"),
		PlanOnly:         getIfSetBool(ctx, "plan-only"),
		IsDestroy:        getIfSetBool(ctx, "is-destroy"),
		Refresh:          getIfSetBool(ctx, "refresh"),
		RefreshOnly:      getIfSetBool(ctx, "refresh-only"),
		Message:          getIfSetString(ctx, "message"),
		Workspace:        &tfe.Workspace{ID: ctx.String("workspace-id")},
		AutoApply:        getIfSetBool(ctx, "auto-apply"),
	}

	if ctx.IsSet("configuration-version") {
		opts.ConfigurationVersion = &tfe.ConfigurationVersion{ID: ctx.String("configuration-version")}
	}

	if ctx.IsSet("var") {
		pairs := strings.Split(ctx.String("var"), ",")

		v := make([]*tfe.RunVariable, len(pairs))

		for i := range pairs {
			kv := strings.Split(pairs[i], "=")
			if len(kv) != 2 {
				return fmt.Errorf("invalid variable format: %s", pairs[i])
			}

			v[i] = &tfe.RunVariable{
				Key:   kv[0],
				Value: kv[1],
			}
		}

		opts.Variables = v
	}

	run, err := tfc.Client.Runs.Create(ctx.Context, opts)
	if err != nil {
		return err
	}

	fmt.Println("updated variable set variable: ")
	r, err := json.MarshalIndent(runCreateResponse{
		ID:              run.ID,
		CreatedAt:       run.CreatedAt,
		AutoApply:       run.AutoApply,
		HasChanges:      run.HasChanges,
		Status:          string(run.Status),
		PositionInQueue: run.PositionInQueue,
	}, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))
	return nil
}
