package asset

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/persistenceOne/persistenceSDK/modules/assets/constants"
	"github.com/persistenceOne/persistenceSDK/types"
)

type queryRequest struct {
	ID types.ID
}

var _ types.QueryRequest = (*queryRequest)(nil)

func (QueryRequest queryRequest) FromCLI(cliCommand types.CLICommand, cliContext context.CLIContext) types.QueryRequest {
	return &queryRequest{ID: types.NewID(cliCommand.ReadString(constants.AssetID))}
}

func (QueryRequest queryRequest) FromMap(vars map[string]string) types.QueryRequest {
	return &queryRequest{
		ID: types.NewID(vars[constants.AssetID.GetName()]),
	}
}

func queryRequestPrototype() types.QueryRequest {
	return &queryRequest{}
}

func NewQueryRequest(id types.ID) types.QueryRequest {
	return &queryRequest{ID: id}
}