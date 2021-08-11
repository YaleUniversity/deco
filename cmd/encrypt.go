package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/YaleUniversity/deco/control"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(encryptionCmd)
	encryptionCmd.AddCommand(encryptionGenKey)
	encryptionCmd.AddCommand(encryptValue)
	encryptionCmd.AddCommand(decryptValue)
}

var encryptionCmd = &cobra.Command{
	Use:   "encryption",
	Short: "Manage encryption",
	Long:  `Manage encryption mechanisms for deco values.`,
}

var encryptionGenKey = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a new encryption key",
	Long:  `Generates a new random 256bit encrytion key.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := control.NewEncryptionKey()
		if err != nil {
			return fmt.Errorf("failed to generate key: %s", err)
		}

		fmt.Printf("%x\n", *key)

		return nil
	},
}

var encryptValue = &cobra.Command{
	Use:   "encrypt [value]",
	Short: "Encrypt a value",
	Long:  `Generates cyphertext of the passed value using the encryption key found in the environment or passed via flags.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("unexpected number of arguments: %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := getKey()
		if err != nil {
			return err
		}

		cipherText, err := control.Encrypt([]byte(args[0]), key)
		if err != nil {
			return fmt.Errorf("failed to encrypt to ciphertext: %s", err)
		}

		fmt.Printf("%x\n", cipherText)

		return nil
	},
}

var decryptValue = &cobra.Command{
	Use:   "decrypt [cipher]",
	Short: "Decrypt a value",
	Long:  `Decryptes passed ciphertext using the encryption key found in the environment or passed via flags.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := getKey()
		if err != nil {
			return err
		}

		cipherBytes, err := hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("invalid ciphertext encoding: %s", err)
		}

		plainText, err := control.Decrypt(cipherBytes, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt ciphertext: %s", err)
		}

		fmt.Printf("%s\n", plainText)

		return nil
	},
}

func getKey() (*[32]byte, error) {
	if encryptionKey == "" {
		return nil, errors.New("missing encryption key")
	}

	keyBytes, err := hex.DecodeString(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key encoding: %s", err)
	}

	var key [32]byte
	copy(key[:], keyBytes)

	return &key, nil
}
