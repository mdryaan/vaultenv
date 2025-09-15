package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var vaultPath string

var rootCmd = &cobra.Command{
	Use:   "vaultenv",
	Short: "A local encrypted vault for developer secrets",
	Long: `vaultenv is a local 1Password for developer environment variables.
It encrypts your secrets with AES-256-GCM using a master password
derived via Argon2id, keeping your credentials safe at rest.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultPath, "vault", "", "path to vault file (overrides VAULTENV_PATH)")
}
