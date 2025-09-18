package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdryaan/vaultenv/internal/config"
	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/output"
	"github.com/mdryaan/vaultenv/internal/prompt"
	"github.com/mdryaan/vaultenv/internal/utils"
	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/spf13/cobra"
)

var (
	setTags     []string
	setGenerate bool
)

var setCmd = &cobra.Command{
	Use:   "set KEY [VALUE]",
	Short: "Add or update a secret",
	Long: `Add or update a secret in the vault. If VALUE is not provided,
you will be prompted to enter it securely (no echo).

Examples:
  vaultenv set DATABASE_URL postgres://localhost/mydb
  vaultenv set API_KEY
  vaultenv set JWT_SECRET --generate
  vaultenv set REDIS_URL redis://localhost --tags production,backend`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSet,
}

func init() {
	setCmd.Flags().StringSliceVar(&setTags, "tags", nil, "comma-separated tags (e.g. production,backend)")
	setCmd.Flags().BoolVar(&setGenerate, "generate", false, "generate a random 32-byte hex secret")
	rootCmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
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

	var value string
	switch {
	case setGenerate:
		value, err = utils.GenerateSecret(32)
		if err != nil {
			return fmt.Errorf("generating secret: %w", err)
		}
	case len(args) == 2:
		value = args[1]
	default:
		value, err = prompt.AskValue(key)
		if err != nil {
			return err
		}
	}

	// Normalize tags
	tags := normalizeTags(setTags)

	if err := v.Set(key, value, tags); err != nil {
		return err
	}

	if setGenerate {
		output.Success(os.Stdout, "Secret %s generated and stored", key)
		fmt.Println(value)
	} else {
		output.Success(os.Stdout, "Secret %s saved", key)
	}
	return nil
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	var result []string
	for _, t := range tags {
		parts := strings.Split(t, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
	}
	return result
}
