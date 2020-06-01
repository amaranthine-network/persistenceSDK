package burn

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/constants"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/mapper"
)

type Keeper interface {
	transact(sdkTypes.Context, Message) error
}

type keeper struct {
	mapper mapper.Mapper
}

func NewKeeper(mapper mapper.Mapper) Keeper {
	return keeper{mapper: mapper}
}

var _ Keeper = (*keeper)(nil)

func (keeper keeper) transact(context sdkTypes.Context, message Message) error {
	assets := keeper.mapper.Assets(context, message.assetID)
	asset := assets.Get(message.assetID)
	if asset == nil {
		return constants.EntityNotFoundCode
	}
	return assets.Remove(asset)
}
