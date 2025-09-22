package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/utils"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var (
	getCopy bool
	getMask bool
)

var getCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Retrieve a secret",
	Long: `Retrieve a secret from the vault by key.

Examples:
  vaultenv get DATABASE_URL
  vaultenv get API_KEY --copy
  vaultenv get API_KEY --mask`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	getCmd.Flags().BoolVar(&getCopy, "copy", false, "copy value to clipboard instead of printing")
	getCmd.Flags().BoolVar(&getMask, "mask", false, "print masked value (****1234)")
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	if err := utils.ValidateKey(key); err != nil {
		return err
	}

	cfg, err := config.Resolve(vaultPath)
	if err != nil {
		return err
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

	entry, err := v.Get(key)
	if err != nil {
		return err
	}

	if getCopy {
		if err := clipboard.WriteAll(entry.Value); err != nil {
			return fmt.Errorf("writing to clipboard: %w", err)
		}
		output.Success(os.Stdout, "Copied %s to clipboard (clears in 30s)", key)

		// Clear clipboard after 30 seconds in the background
		go clearClipboardAfter(entry.Value, 30*time.Second)
		return nil
	}

	if getMask {
		fmt.Println(utils.MaskValue(entry.Value))
		return nil
	}

	fmt.Println(entry.Value)
	return nil
}

func clearClipboardAfter(originalValue string, d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	<-ctx.Done()
	current, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	// Only clear if clipboard still has our value
	if current == originalValue {
		clipboard.WriteAll("") //nolint:errcheck
	}
}
