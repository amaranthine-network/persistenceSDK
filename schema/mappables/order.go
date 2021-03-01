/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package mappables

import (
	"github.com/persistenceOne/persistenceSDK/schema/helpers"
	"github.com/persistenceOne/persistenceSDK/schema/traits"
	"github.com/persistenceOne/persistenceSDK/schema/types"
)

type Order interface {
	GetID() types.ID
	GetClassificationID() types.ID
	GetRateID() types.ID
	GetCreationID() types.ID
	GetMakerOwnableID() types.ID
	GetTakerOwnableID() types.ID
	GetMakerID() types.ID
	GetKey() helpers.Key

	GetCreation() types.MetaProperty
	GetExchangeRate() types.MetaProperty

	GetTakerID() types.Property
	GetExpiry() types.Property
	GetMakerOwnableSplit() types.Property

	traits.HasMutables
	traits.HasImmutables
	helpers.Mappable
}
