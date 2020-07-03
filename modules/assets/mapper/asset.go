package mapper

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

var _ types.InterNFT = (*asset)(nil)

type asset struct {
	ID         types.ID
	Mutables   types.Mutables
	Immutables types.Immutables
	Lock       types.Height
	Burn       types.Height
}

func (asset asset) GetID() types.ID {
	return asset.ID
}

func (asset asset) GetChainID() types.ID {
	return assetIDFromInterface(asset.ID).ChainID
}

func (asset asset) GetClassificationID() types.ID {
	return assetIDFromInterface(asset.ID).ClassificationID
}

func (asset asset) GetMaintainersID() types.ID {
	return assetIDFromInterface(asset.ID).MaintainersID
}

func (asset asset) GetHashID() types.ID {
	return asset.Immutables.GetHashID()
}

func (asset asset) GetMutables() types.Mutables {
	return asset.Mutables
}

func (asset asset) GetImmutables() types.Immutables {
	return asset.Immutables
}

func (asset asset) GetLock() types.Height {
	return asset.Lock
}

func (asset asset) CanSend(currentHeight types.Height) bool {
	return currentHeight.IsGreaterThan(asset.Lock)
}

func (asset asset) GetBurn() types.Height {
	return asset.Burn
}

func (asset asset) CanBurn(currentHeight types.Height) bool {
	return currentHeight.IsGreaterThan(asset.Burn)
}

func NewAsset(assetID types.ID, mutables types.Mutables, immutables types.Immutables, lock types.Height, burn types.Height) types.InterNFT {
	return asset{
		ID:         assetID,
		Mutables:   mutables,
		Immutables: immutables,
		Lock:       lock,
		Burn:       burn,
	}
}