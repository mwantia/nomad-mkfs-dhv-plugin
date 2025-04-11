package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mwantia/nomad-mkfs-host-volume-plugin/pkg/plugin"
	"github.com/spf13/cobra"
)

type StdErrResponse struct {
	Error string `json:"error"`
}

var (
	Root = &cobra.Command{
		Use:           "mkfs",
		Short:         "",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	Fingerprint = &cobra.Command{
		Use:   "fingerprint",
		Short: "Displays the version; Is also used to validate the plugin during startup",
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.Fingerprint()
		},
	}
	Create = &cobra.Command{
		Use:   "create",
		Short: "Creates a new mount with the provided nomad host volume configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.Create()
		},
	}
	Delete = &cobra.Command{
		Use:   "delete",
		Short: "Deletes the mount defined during nomad host volume deletion",
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.Delete()
		},
	}
)

func main() {
	Root.AddCommand(Fingerprint, Create, Delete)

	if err := Root.Execute(); err != nil {
		resp := StdErrResponse{
			Error: err.Error(),
		}
		if json, err := json.Marshal(resp); err == nil {
			fmt.Print(string(json))
		}

		os.Exit(1)
	}
}
