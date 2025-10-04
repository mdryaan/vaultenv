package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/dotenv"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportTags   []string
	exportKeys   []string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export secrets to a .env file",
	Long: `Export secrets from the vault as a .env file. By default prints to stdout.

Examples:
  vaultenv export
  vaultenv export --output .env
  vaultenv export --tags production
  vaultenv export --keys DB_URL,API_KEY`,
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "write to file instead of stdout")
	exportCmd.Flags().StringSliceVar(&exportTags, "tags", nil, "export only secrets with these tags")
	exportCmd.Flags().StringSliceVar(&exportKeys, "keys", nil, "export only specific keys (comma-separated)")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
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

	tags := normalizeTags(exportTags)
	entries := v.List(tags)

	// Filter by explicit keys if provided
	keySet := buildKeySet(exportKeys)
	var dotentries []dotenv.Entry
	for _, e := range entries {
		if len(keySet) > 0 && !keySet[e.Key] {
			continue
		}
		dotentries = append(dotentries, dotenv.Entry{Key: e.Key, Value: e.Value})
	}

	if exportOutput == "" {
		return dotenv.Write(os.Stdout, dotentries)
	}

	f, err := os.OpenFile(exportOutput, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("opening output file: %w", err)
	}
	defer f.Close()

	if err := dotenv.Write(f, dotentries); err != nil {
		return err
	}

	output.Success(os.Stderr, "Exported %d secret(s) to %s", len(dotentries), exportOutput)
	return nil
}

func buildKeySet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	set := make(map[string]bool)
	for _, k := range keys {
		for _, part := range strings.Split(k, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				set[part] = true
			}
		}
	}
	return set
}
