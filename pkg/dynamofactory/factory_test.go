package dynamofactory_test

import (
	"github.com/fionahiklas/simple-dynamodb-go-code/pkg/dynamofactory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFactory(t *testing.T) {

	t.Run("create factory returns non nil", func(t *testing.T) {
		result := dynamofactory.NewFactory()
		require.NotNil(t, result)
	})

}
