package redisbloomfilter

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
	"github.com/stretchr/testify/require"
)

func TestRedisBackend(t *testing.T) {
	var (
		filterName = "unit-testing"
		key1       = "test-key1"
		key2       = "test-key2"
	)

	conn := redigomock.NewConn()

	reserve := conn.Command(
		"BF.RESERVE",
		filterName,
		fmt.Sprintf("%.1f", defaultErrorRate),
		defaultCapacity,
	)

	addKey1 := conn.Command(
		"BF.ADD",
		filterName,
		key1,
	).Expect(int64(1))

	existsKey1 := conn.Command(
		"BF.EXISTS",
		filterName,
		key1,
	).Expect(int64(1))

	existsKey2 := conn.Command(
		"BF.EXISTS",
		filterName,
		key2,
	).Expect(int64(0))

	addKey2 := conn.Command(
		"BF.ADD",
		filterName,
		key2,
	).Expect(int64(1))

	pool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}

	rbf, err := New(FilterName(filterName), pool)
	require.NoError(t, err)

	require.Equal(t, 1, conn.Stats(reserve), "BF.RESERVE was not called exactly once")

	t.Run("key1 can be set and tested for existence", func(t *testing.T) {
		require.NoError(t, rbf.Set(key1))
		require.Equal(t, 1, conn.Stats(addKey1))

		exists, err := rbf.Exists(key1)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, 1, conn.Stats(existsKey1))
	})

	t.Run("key2 should not exist yet", func(t *testing.T) {
		exists, err := rbf.Exists(key2)
		require.NoError(t, err)
		require.False(t, exists)
		require.Equal(t, 1, conn.Stats(existsKey2))
	})

	t.Run("key2 can be set with CheckAndSet", func(t *testing.T) {
		exists, err := rbf.CheckAndSet(key2)
		require.NoError(t, err)
		require.False(t, exists)
		require.Equal(t, 1, conn.Stats(addKey2))

		addKey2 = conn.Command(
			"BF.ADD",
			filterName,
			key2,
		).Expect(int64(0))

		exists, err = rbf.CheckAndSet(key2)
		require.NoError(t, err)
		require.True(t, exists)

		require.Equal(t, 2, conn.Stats(addKey2))
	})
}

func TestReserveIdempotency(t *testing.T) {
	filterName := "unit-testing"

	conn := redigomock.NewConn()

	reserve := conn.Command(
		"BF.RESERVE",
		filterName,
		fmt.Sprintf("%.1f", defaultErrorRate),
		defaultCapacity,
	).ExpectError(redis.Error("ERR item exists"))

	pool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}

	_, err := New(FilterName(filterName), pool)
	require.NoError(t, err, "should not have recieved an error from New")

	require.Equal(t, 1, conn.Stats(reserve), "BF.RESERVE was not called exactly once")

	reserve = conn.Command(
		"BF.RESERVE",
		filterName,
		fmt.Sprintf("%.1f", defaultErrorRate),
		defaultCapacity,
	).ExpectError(redis.Error("unexpected error"))

	_, err = New(FilterName(filterName), pool)
	require.Error(t, err, "should have recieved an error from New")
}
