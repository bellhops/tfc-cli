package app

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/urfave/cli/v2"
)

// ConfigurationVersions describes all the configuration version related
// methods that the Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/enterprise/api/configuration-versions.html
//	// List returns all configuration versions of a workspace.
//	List(ctx context.Context, workspaceID string, options *ConfigurationVersionListOptions) (*ConfigurationVersionList, error)
//
//	// Create is used to create a new configuration version. The created
//	// configuration version will be usable once data is uploaded to it.
//	Create(ctx context.Context, workspaceID string, options ConfigurationVersionCreateOptions) (*ConfigurationVersion, error)
//
//	// Read a configuration version by its ID.
//	Read(ctx context.Context, cvID string) (*ConfigurationVersion, error)
//
//	// ReadWithOptions reads a configuration version by its ID using the options supplied
//	ReadWithOptions(ctx context.Context, cvID string, options *ConfigurationVersionReadOptions) (*ConfigurationVersion, error)
//
//	// Upload packages and uploads Terraform configuration files. It requires
//	// the upload URL from a configuration version and the full path to the
//	// configuration files on disk.
//	Upload(ctx context.Context, url string, path string) error
//
//	// Archive a configuration version. This can only be done on configuration versions that
//	// were created with the API or CLI, are in an uploaded state, and have no runs in progress.
//	Archive(ctx context.Context, cvID string) error
//
//	// Download a configuration version.  Only configuration versions in the uploaded state may be downloaded.
//	Download(ctx context.Context, cvID string) ([]byte, error)

var configVerIncludeOpts = map[string]tfe.ConfigVerIncludeOpt{
	"ingress_attributes": tfe.ConfigVerIngressAttributes,
	"run":                tfe.ConfigVerRun,
}

func (tfc *TFCClient) ConfigVersionsListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Aliases:  []string{"ls"},
		Usage:    "List all the terraform configuration versions within an organization.",
		Category: "configuration versions",
		Action:   tfc.configVersionsList,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "include",
				Usage:   "A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api/workspaces.html#available-related-resources",
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

func (tfc *TFCClient) configVersionsList(ctx *cli.Context) error {
	opts := &tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: ctx.Int("page-num"),
			PageSize:   ctx.Int("page-size"),
		},
		Include: nil,
	}

	include := ctx.StringSlice("include")

	for _, r := range include {
		if opt, ok := configVerIncludeOpts[r]; !ok {
			return fmt.Errorf("include opt not recognized: %s", r)
		} else {
			if opts.Include == nil {
				opts.Include = []tfe.ConfigVerIncludeOpt{}
			}
			opts.Include = append(opts.Include, opt)
		}
	}

	cvl, err := tfc.Client.ConfigurationVersions.List(ctx.Context, tfc.Cfg.OrgName, opts)
	if err != nil {
		return err
	}

	r, err := json.MarshalIndent(cvl, "", "    ")
	if err != nil {
		return nil
	}

	fmt.Println(string(r))

	return nil
}
