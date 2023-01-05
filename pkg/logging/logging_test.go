package logging_test

import (
	"testing"

	"github.com/fionahiklas/simple-dynamodb-go-code/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {

	var mockConfig *Mockconfig

	resetTest := func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfig = NewMockconfig(ctrl)
	}

	t.Run("create new logger", func(t *testing.T) {
		resetTest(t)
		mockConfig.EXPECT().LogLevel().Return("Debug")
		mockConfig.EXPECT().EnableJsonLogFormat().Return(false)
		result := logging.NewLogger(mockConfig)
		require.NotNil(t, result)
		require.Equal(t, logrus.DebugLevel, result.GetLevel())
	})

	t.Run("create new logger default to info", func(t *testing.T) {
		resetTest(t)
		mockConfig.EXPECT().LogLevel().Return("wibble")
		mockConfig.EXPECT().EnableJsonLogFormat().Return(false)
		result := logging.NewLogger(mockConfig)
		require.NotNil(t, result)
		require.Equal(t, logrus.InfoLevel, result.GetLevel())
	})

	t.Run("enable json output", func(t *testing.T) {
		resetTest(t)
		mockConfig.EXPECT().LogLevel().Return("Debug")
		mockConfig.EXPECT().EnableJsonLogFormat().Return(true)
		result := logging.NewLogger(mockConfig)
		require.NotNil(t, result)
		require.Equal(t, logrus.DebugLevel, result.GetLevel())
		require.IsType(t, &logrus.JSONFormatter{}, result.Formatter)
	})

	t.Run("create new test logger", func(t *testing.T) {
		resetTest(t)
		result := logging.NewTestLogger()
		require.Equal(t, logrus.DebugLevel, result.GetLevel())
		require.NotNil(t, result)
	})

}
