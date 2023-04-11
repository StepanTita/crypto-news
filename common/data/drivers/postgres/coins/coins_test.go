//go:build integration

package coins

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"common/data/drivers"
	"common/data/model"
)

func createTestCoins(t *testing.T, log *logrus.Entry, db *sqlx.DB) *coins {
	t.Helper()

	return &coins{
		log: log.WithField("service", "[INSERTER-INTEGRATION-TEST]"),
		ext: db,
	}
}

func UpsertCoinsBatch(t *testing.T, log *logrus.Entry, db *sqlx.DB) {
	testCoins := createTestCoins(t, log, db)

	ctx := context.Background()

	inCoins := []model.Coin{
		{
			Code:  "BTC",
			Title: "Bitcoin",
			Slug:  "bitcoin",
		},
		{
			Code:  "ETH",
			Title: "Ethereum",
			Slug:  "ethereum",
		},
	}

	outCoins := make([]model.Coin, len(inCoins))
	copy(outCoins, inCoins)

	testCases := []struct {
		name         string
		precondition func()
		in           []model.Coin
		out          []model.Coin
		err          error
	}{
		{
			name: "ok coins insert batch",
			in:   inCoins,
			out:  outCoins,
			err:  nil,
		},
		{
			name: "ok coins upsert batch",
			precondition: func() {
				for i := range inCoins {
					inCoins[i].Slug += "2"
					outCoins[i].Slug += "2"
				}
			},
			in:  inCoins,
			out: outCoins,
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.precondition != nil {
				tc.precondition()
			}
			err := testCoins.UpsertCoinsBatch(ctx, tc.in)
			require.ErrorIs(t, err, tc.err)
			for i := range tc.in {
				require.EqualValues(t, tc.out[i], tc.in[i])
			}
		})
	}
}

func TestCoins(t *testing.T) {
	suite := drivers.NewSuite(t)
	suite.AddTests(UpsertCoinsBatch)

	suite.SetupSuite()
	defer suite.CleanupSuite()

	suite.TestRunIntegration()
}
