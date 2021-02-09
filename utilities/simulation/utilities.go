package simulation

import (
	"math/rand"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func RandomBool(r *rand.Rand) bool {
	return r.Int()%2 == 0
}

func GenerateRandomAddresses(r *rand.Rand) []sdkTypes.AccAddress {
	randomAccounts := simulation.RandomAccounts(r, r.Int())

	addresses := make([]sdkTypes.AccAddress, len(randomAccounts))
	for i, account := range randomAccounts {
		addresses[i] = account.Address
	}

	return addresses
}
