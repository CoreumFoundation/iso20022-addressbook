package iso20022_addressbook

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"github.com/CoreumFoundation/iso20022-client/iso20022/addressbook"
)

func TestValidateAddresses(t *testing.T) {
	ctx := context.Background()

	chainIds := []string{"coreum-mainnet-1", "coreum-testnet-1", "coreum-devnet-1"}
	for _, chainId := range chainIds {
		ab := addressbook.NewWithRepoAddress(chainId, "file://./%s/addressbook.json")

		err := ab.Update(ctx)
		if err != nil {
			t.Fatalf("could not update %s address book: %v", chainId, err)
		}

		localAddressBook := make(map[string]struct{})
		ab.ForEach(func(address addressbook.Address) bool {
			if _, alreadyExists := localAddressBook[address.Bech32EncodedAddress]; alreadyExists {
				t.Fatalf(
					"duplicate entries with bech32 encoded address %q in %s",
					address.Bech32EncodedAddress,
					chainId,
				)
			}

			publicKeyBytes, err := base64.StdEncoding.DecodeString(address.PublicKey)
			if err != nil {
				t.Fatalf(
					"public key of %q is not a valid base64 encoded string in %s: %v",
					address.Bech32EncodedAddress,
					chainId,
					err,
				)
			}

			if _, err = secp256k1.ParsePubKey(publicKeyBytes); err != nil {
				t.Fatalf(
					"public key of %q is not a valid secp256k1 public key in %s: %v",
					address.Bech32EncodedAddress,
					chainId,
					err,
				)
			}
			localAddressBook[address.Bech32EncodedAddress] = struct{}{}
			matches := make([]string, 0, 1)
			ab.ForEach(func(address2 addressbook.Address) bool {
				if address.BranchAndIdentification.Equal(address2.BranchAndIdentification) {
					matches = append(matches, address2.Bech32EncodedAddress)
					if len(matches) > 1 {
						return false
					}
				}
				return true
			})
			if len(matches) > 1 {
				t.Fatalf(
					"ISO20022 branch and identification data of entry %q and %q conflicts in %s.",
					matches[0],
					matches[1],
					chainId,
				)
			}
			return true
		})
	}
}
