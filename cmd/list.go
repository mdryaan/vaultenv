package cmd

import (
	"os"
	"strings"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var (
	listShowValues bool
	listTags       []string
	listOutput     string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in the vault",
	Long: `List all secrets stored in the vault. Values are masked by default.

Examples:
  vaultenv list
  vaultenv list --show-values
  vaultenv list --tags production
  vaultenv list --output json`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVar(&listShowValues, "show-values", false, "display full secret values (use with caution)")
	listCmd.Flags().StringSliceVar(&listTags, "tags", nil, "filter by tag(s)")
	listCmd.Flags().StringVar(&listOutput, "output", "table", "output format: table, json")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
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

	tags := normalizeTags(listTags)
	entries := v.List(tags)

	if len(entries) == 0 {
		output.Info(os.Stderr, "No secrets found%s", tagSuffix(tags))
		return nil
	}

	formatter := output.Get(output.Format(strings.ToLower(listOutput)))
	return formatter.WriteEntries(os.Stdout, entries, listShowValues)
}

func tagSuffix(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return " with tags: " + strings.Join(tags, ", ")
}
