package uuid_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/stretchr/testify/require"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

var defaultConnTestRunner pgxtest.ConnTestRunner

func init() {
	defaultConnTestRunner = pgxtest.DefaultConnTestRunner()
	defaultConnTestRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxuuid.Register(conn.TypeMap())
	}
}

func TestCodecDecodeValue(t *testing.T) {
	defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		original, err := uuid.NewRandom()
		require.NoError(t, err)

		rows, err := conn.Query(ctx, `select $1::uuid`, original)
		require.NoError(t, err)

		for rows.Next() {
			values, err := rows.Values()
			require.NoError(t, err)

			require.Len(t, values, 1)
			v0, ok := values[0].(uuid.UUID)
			require.True(t, ok)
			require.Equal(t, original, v0)
		}

		require.NoError(t, rows.Err())

		rows, err = conn.Query(ctx, `select $1::uuid`, nil)
		require.NoError(t, err)

		for rows.Next() {
			values, err := rows.Values()
			require.NoError(t, err)

			require.Len(t, values, 1)
			require.Equal(t, nil, values[0])
		}

		require.NoError(t, rows.Err())
	})
}

func TestArray(t *testing.T) {
	defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		var inputSlice []uuid.UUID

		for i := 0; i < 10; i++ {
			u, err := uuid.NewRandom()
			require.NoError(t, err)
			inputSlice = append(inputSlice, u)
		}

		var outputSlice []uuid.UUID
		err := conn.QueryRow(ctx, `select $1::uuid[]`, inputSlice).Scan(&outputSlice)
		require.NoError(t, err)
		require.Equal(t, inputSlice, outputSlice)
	})
}
