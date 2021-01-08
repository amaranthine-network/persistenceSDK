/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package mappables

import (
	"github.com/persistenceOne/persistenceSDK/schema/helpers"
	"github.com/persistenceOne/persistenceSDK/schema/types"
)

type Maintainer interface {
	GetClassificationID() types.ID
	GetIdentityID() types.ID

	CanAddMaintainer() bool
	CanRemoveMaintainer() bool
	CanMutateMaintainer() bool

	MaintainsTrait(types.ID) bool
	helpers.Mappable
}
