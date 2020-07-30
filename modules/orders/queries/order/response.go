package order

import (
	"github.com/persistenceOne/persistenceSDK/schema/helpers"
	"github.com/persistenceOne/persistenceSDK/schema/mappers"
)

type queryResponse struct {
	Orders mappers.Orders
}

var _ helpers.QueryResponse = (*queryResponse)(nil)

func queryResponsePrototype() helpers.QueryResponse {
	return queryResponse{}
}

func newQueryResponse(orders mappers.Orders) helpers.QueryResponse {
	return queryResponse{Orders: orders}
}