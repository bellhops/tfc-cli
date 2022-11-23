package app

import "github.com/urfave/cli/v2"

func ptrString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func ptrBool(v bool) *bool {
	return &v
}

func getIfSetBool(ctx *cli.Context, key string) *bool {
	if !ctx.IsSet(key) {
		return nil
	}
	r := ctx.Bool(key)
	return &r
}

func getIfSetString(ctx *cli.Context, key string) *string {
	if !ctx.IsSet(key) {
		return nil
	}
	r := ctx.String(key)
	return &r
}
