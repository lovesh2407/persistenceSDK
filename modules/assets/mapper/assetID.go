package mapper

import (
	"bytes"
	"github.com/persistenceOne/persistenceSDK/modules/assets/constants"
	"github.com/persistenceOne/persistenceSDK/types"
	"strings"
)

var _ types.ID = (*assetID)(nil)

type assetID struct {
	ChainID          types.ID
	MaintainersID    types.ID
	ClassificationID types.ID
	HashID           types.ID
}

func (assetID assetID) Bytes() []byte {
	return append(append(append(
		assetID.ChainID.Bytes(),
		assetID.MaintainersID.Bytes()...),
		assetID.ClassificationID.Bytes()...),
		assetID.HashID.Bytes()...)
}

func (assetID assetID) String() string {
	var values []string
	values = append(values, assetID.ChainID.String())
	values = append(values, assetID.MaintainersID.String())
	values = append(values, assetID.ClassificationID.String())
	values = append(values, assetID.HashID.String())
	return strings.Join(values, constants.IDSeparator)
}

func (assetID assetID) Compare(id types.ID) int {
	return bytes.Compare(assetID.Bytes(), id.Bytes())
}

func assetIDFromInterface(id types.ID) assetID {
	idList := strings.Split(id.String(), constants.IDSeparator)
	if len(idList) == 4 {
		return assetID{
			ChainID:          types.NewID(idList[0]),
			MaintainersID:    types.NewID(idList[1]),
			ClassificationID: types.NewID(idList[2]),
			HashID:           types.NewID(idList[3]),
		}
	}
	return assetID{ChainID: types.NewID(""), MaintainersID: types.NewID(""), ClassificationID: types.NewID(""), HashID: types.NewID("")}
}

func NewAssetID(chainID types.ID, maintainersID types.ID, classificationID types.ID, hashID types.ID) types.ID {
	return &assetID{
		ChainID:          chainID,
		MaintainersID:    maintainersID,
		ClassificationID: classificationID,
		HashID:           hashID,
	}
}