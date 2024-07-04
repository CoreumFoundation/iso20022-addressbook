package iso20022_addressbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/CoreumFoundation/iso20022-client/iso20022/addressbook"
)

func TestValidateAddresses(t *testing.T) {
	ctx := context.Background()

	chainIds := []string{"coreum-mainnet-1", "coreum-testnet-1", "coreum-devnet-1"}
	for _, chainId := range chainIds {
		ab := addressbook.NewWithRepoAddress(fmt.Sprintf("file://./%s/addressbook.json", chainId))

		err := ab.Update(ctx)
		if err != nil {
			t.Fatalf("could not update %s address book: %v", chainId, err)
		}

		err = ab.Validate()
		if err != nil {
			t.Fatalf(
				"%s addressbook error: %v",
				chainId,
				err,
			)
		}
	}
}
