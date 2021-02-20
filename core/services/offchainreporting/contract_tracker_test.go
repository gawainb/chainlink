package offchainreporting_test

import (
	"context"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/require"
)

func Test_OCRContractTracker_HandleLog_OCRContractLatestRoundRequested(t *testing.T) {
	fixtureLogAddress := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(fixtureLogAddress, nil)
	require.NoError(t, err)

	t.Run("does not update if contract address doesn't match", func(t *testing.T) {
		tracker, err := offchainreporting.NewOCRContractTracker(
			cltest.NewAddress(),
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
		)
		require.NoError(t, err)
		logBroadcast := new(mocks.Broadcast)

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log.json")
		logBroadcast.On("RawLog").Return(rawLog)
		logBroadcast.On("MarkConsumed").Return(nil)
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)

		configDigest, round, epoch, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		tracker.HandleLog(logBroadcast, nil)

		configDigest, round, epoch, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
	})

	t.Run("does nothing if log has already been consumed", func(t *testing.T) {
		tracker, err := offchainreporting.NewOCRContractTracker(
			cltest.NewAddress(),
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
		)
		require.NoError(t, err)
		logBroadcast := new(mocks.Broadcast)

		logBroadcast.On("WasAlreadyConsumed").Return(true, nil)

		configDigest, round, epoch, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		tracker.HandleLog(logBroadcast, nil)

		configDigest, round, epoch, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
	})

	t.Run("for new round requested log", func(t *testing.T) {
		tracker, err := offchainreporting.NewOCRContractTracker(
			fixtureLogAddress,
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
		)
		require.NoError(t, err)
		logBroadcast := new(mocks.Broadcast)

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log.json")
		logBroadcast.On("RawLog").Return(rawLog)
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast.On("MarkConsumed").Return(nil)

		configDigest, round, epoch, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		tracker.HandleLog(logBroadcast, nil)

		configDigest, round, epoch, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, "92eed6d902b18bdbf881402d47631bcd", configDigest.Hex())
		require.Equal(t, 1, int(round))
		require.Equal(t, 1, int(epoch))

		logBroadcast.AssertExpectations(t)
	})
}
