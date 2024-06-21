package iso20022_addressbook

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"github.com/CoreumFoundation/iso20022-client/iso20022/addressbook"
)

func TestValidateAddresses(t *testing.T) {
	networkTypes := []string{"mainnet", "testnet", "devnet"}
	for _, networkType := range networkTypes {
		ab, err := addressbook.NewWithRepoAddress(networkType, "file://./%s/addressbook.json")
		if err != nil {
			t.Fatalf("could not get %s address book : %v", networkType, err)
		}

		localAddressBook := make(map[string]struct{})
		ab.ForEach(func(address addressbook.Address) bool {
			if _, alreadyExists := localAddressBook[address.Bech32EncodedAddress]; alreadyExists {
				t.Fatalf(
					"duplicate entries with bech32 encoded address %q in %s",
					address.Bech32EncodedAddress,
					networkType,
				)
			}

			publicKeyBytes, err := base64.StdEncoding.DecodeString(address.PublicKey)
			if err != nil {
				t.Fatalf(
					"public key of %q is not a valid base64 encoded string in %s : %v",
					address.Bech32EncodedAddress,
					networkType,
					err,
				)
			}

			switch ab.KeyAlgo() {
			case "secp256k1":
				if _, err = secp256k1.ParsePubKey(publicKeyBytes); err != nil {
					t.Fatalf(
						"public key of %q is not a valid secp256k1 public key in %s : %v",
						address.Bech32EncodedAddress,
						networkType,
						err,
					)
				}
			case "secp256r1":
				pbKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
				if err != nil {
					t.Fatalf(
						"public key of %q is not a valid secp256r1 public key in %s : %v",
						address.Bech32EncodedAddress,
						networkType,
						err,
					)
				}
				publicKey, ok := pbKey.(*ecdsa.PublicKey)
				if !ok {
					t.Fatalf(
						"public key of %q is not a valid secp256r1 public key in %s",
						address.Bech32EncodedAddress,
						networkType,
					)
				}
				_, err = publicKey.ECDH()
				if err != nil {
					t.Fatalf(
						"public key of %q is not a valid secp256r1 public key in %s : %v",
						address.Bech32EncodedAddress,
						networkType,
						err,
					)
				}
			default:
				t.Fatalf(
					"key algorithm %q is not supported in %s",
					ab.KeyAlgo(),
					networkType,
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
					networkType,
				)
			}
			return true
		})
	}
}
