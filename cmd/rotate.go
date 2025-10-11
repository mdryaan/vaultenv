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

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Change the vault master password",
	Long: `Rotate the master password: decrypt with the old password and
re-encrypt the entire vault with the new password.

Example:
  vaultenv rotate`,
	RunE: runRotate,
}

func init() {
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Resolve(vaultPath)
	if err != nil {
		return err
	}

	var oldPassword []byte
	if cfg.Password != "" {
		oldPassword = []byte(cfg.Password)
	} else {
		oldPassword, err = prompt.AskPassword("Current master password: ")
		if err != nil {
			return err
		}
	}
	defer crypto.ZeroBytes(oldPassword)

	v, err := vault.Open(cfg.VaultPath, oldPassword)
	if err != nil {
		return fmt.Errorf("opening vault (wrong password?): %w", err)
	}

	newPassword, err := prompt.AskPasswordConfirm(
		"New master password: ",
		"Confirm new master password: ",
	)
	if err != nil {
		return err
	}
	defer crypto.ZeroBytes(newPassword)

	if len(newPassword) == 0 {
		return fmt.Errorf("new master password must not be empty")
	}

	if err := v.Rotate(newPassword); err != nil {
		return err
	}

	output.Success(os.Stdout, "Master password rotated successfully")
	return nil
}
