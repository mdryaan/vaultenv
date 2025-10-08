package cmd

import (
	"fmt"
	"os"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/dotenv"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var (
	importOverwrite bool
	importTags      []string
	importDryRun    bool
)

var importCmd = &cobra.Command{
	Use:   "import FILE",
	Short: "Import secrets from a .env file",
	Long: `Import secrets from a .env file into the vault.

Examples:
  vaultenv import .env
  vaultenv import .env --overwrite
  vaultenv import .env --tags staging
  vaultenv import .env --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "overwrite existing keys")
	importCmd.Flags().StringSliceVar(&importTags, "tags", nil, "tag all imported secrets")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "show what would be imported without making changes")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	entries, err := dotenv.Parse(f)
	if err != nil {
		return fmt.Errorf("parsing .env file: %w", err)
	}

	if importDryRun {
		fmt.Printf("Would import %d secret(s) from %s:\n", len(entries), filePath)
		for _, e := range entries {
			fmt.Printf("  %s\n", e.Key)
		}
		return nil
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

	tags := normalizeTags(importTags)
	imported := 0
	skipped := 0

	for _, e := range entries {
		_, exists := v.Data().Get(e.Key)
		if exists && !importOverwrite {
			skipped++
			continue
		}
		if err := v.Set(e.Key, e.Value, tags); err != nil {
			return fmt.Errorf("setting %s: %w", e.Key, err)
		}
		imported++
	}

	output.Success(os.Stdout, "Imported %d secret(s) from %s", imported, filePath)
	if skipped > 0 {
		output.Warn(os.Stderr, "Skipped %d existing secret(s) (use --overwrite to replace)", skipped)
	}
	return nil
}
