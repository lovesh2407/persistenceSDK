package scrub

import (
	"github.com/persistenceOne/persistenceSDK/constants/errors"
	"github.com/persistenceOne/persistenceSDK/schema/types"
	"github.com/persistenceOne/persistenceSDK/schema/types/base"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Super_Response(t *testing.T) {

	metaProperty := base.NewMetaProperty(base.NewID("id"), base.NewMetaFact(base.NewStringData("Data")))
	metaPropertyList := base.NewMetaProperties([]types.MetaProperty{metaProperty})

	testAuxiliaryResponse := newAuxiliaryResponse(metaPropertyList, nil)
	require.Equal(t, auxiliaryResponse{Success: true, Error: nil, Properties: metaPropertyList}, testAuxiliaryResponse)
	require.Equal(t, true, testAuxiliaryResponse.IsSuccessful())
	require.Equal(t, nil, testAuxiliaryResponse.GetError())

	testAuxiliaryResponse2 := newAuxiliaryResponse(metaPropertyList, errors.IncorrectFormat)
	require.Equal(t, auxiliaryResponse{Success: false, Error: errors.IncorrectFormat, Properties: nil}, testAuxiliaryResponse2)
	require.Equal(t, false, testAuxiliaryResponse2.IsSuccessful())
	require.Equal(t, errors.IncorrectFormat, testAuxiliaryResponse2.GetError())

	propertiesFromResponse, Error := GetPropertiesFromResponse(testAuxiliaryResponse)
	require.Equal(t, metaPropertyList, propertiesFromResponse)
	require.Equal(t, nil, Error)

	propertiesFromResponse2, Error := GetPropertiesFromResponse(testAuxiliaryResponse2)
	require.Equal(t, nil, propertiesFromResponse2)
	require.Equal(t, errors.IncorrectFormat, Error)

	propertiesFromResponse3, Error := GetPropertiesFromResponse(nil)
	require.Equal(t, nil, propertiesFromResponse3)
	require.Equal(t, errors.NotAuthorized, Error)
}