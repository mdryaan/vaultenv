package cmd

import (
	"fmt"
	"os"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long: `Initialize a new encrypted vault at the default location or a custom path.
The vault is encrypted with AES-256-GCM using a key derived from your
master password via Argon2id.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cfg, err := config.Resolve(vaultPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfg.VaultPath); err == nil {
		return fmt.Errorf("vault already exists at %s\nUse a different path with --vault or delete the existing file", cfg.VaultPath)
	}

	var password []byte
	if cfg.Password != "" {
		password = []byte(cfg.Password)
	} else {
		password, err = prompt.AskPasswordConfirm(
			"Enter master password: ",
			"Confirm master password: ",
		)
		if err != nil {
			return err
		}
	}
	defer crypto.ZeroBytes(password)

	if len(password) == 0 {
		return fmt.Errorf("master password must not be empty")
	}

	if _, err := vault.Create(cfg.VaultPath, password); err != nil {
		return fmt.Errorf("creating vault: %w", err)
	}

	output.Success(os.Stdout, "Vault initialized at %s", cfg.VaultPath)
	return nil
}
