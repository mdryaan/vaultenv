package cmd

import (
	"fmt"
	"os"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/utils"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Remove a secret from the vault",
	Long: `Remove a secret from the vault by key.

Examples:
  vaultenv delete DATABASE_URL
  vaultenv delete DATABASE_URL --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	key := args[0]
	if err := utils.ValidateKey(key); err != nil {
		return err
	}

	cfg, err := config.Resolve(vaultPath)
	if err != nil {
		return err
	}

	if !deleteForce {
		ok, err := prompt.Confirm(fmt.Sprintf("Delete secret %q?", key))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(os.Stderr, "Aborted.")
			return nil
		}
	}

	var password []byte
	if cfg.Password != "" {
		password = []byte(cfg.Password)
	} else {
		password, err = prompt.AskPassword("Master password: ")
		if err != nil {
			return err
		}
	}
	defer crypto.ZeroBytes(password)

	v, err := vault.Open(cfg.VaultPath, password)
	if err != nil {
		return err
	}

	if err := v.Delete(key); err != nil {
		return err
	}

	output.Success(os.Stdout, "Deleted %s", key)
	return nil
}
