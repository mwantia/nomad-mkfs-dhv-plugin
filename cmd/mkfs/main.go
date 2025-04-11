package main

import "github.com/spf13/cobra"

var (
	Root = &cobra.Command{
		Use: "nomad-mkfs-host-volume-plugin",
		Short: "",
	}
	Fingerprint = &cobra.Command{
		Use: "fingerprint",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		}
	}
	Create = &cobra.Command{
		Use: "create",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		}
	}
	Delete = &cobra.Command{
		Use: "delete",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		}
	}
)

func main() {
	Root.AddCommand(Fingerprint, Create, Delete)

	if err := Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}