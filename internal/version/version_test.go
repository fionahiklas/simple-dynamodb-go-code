package version_test

import (
	"os"
	"strings"
	"testing"

	"github.com/fionahiklas/simple-dynamodb-go-code/internal/version"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {

	versionBytes, err := os.ReadFile("../../VERSION")
	require.NoError(t, err)

	versionString := string(versionBytes)
	require.NotEmpty(t, versionString)

	require.Equal(t, strings.TrimSpace(versionString), version.CodeVersion())

	commitHash := version.CommitHash()
	require.NotEmpty(t, commitHash)
	require.Greater(t, len(commitHash), 5)

	appName := version.AppName()
	require.Equal(t, "product-service", appName)
}
