package assetFactory

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/queries/asset"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/transactions/mint"
)

func RegisterCodec(codec *codec.Codec) {
	mapper.RegisterCodec(codec)

	asset.Query.RegisterCodec(codec)
	mint.Transaction.RegisterCodec(codec)
}