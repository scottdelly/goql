package runner

import (
	"database/sql"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/mgutz/logxi/v1"
	"gopkg.in/mgutz/dat.v2/dat"
	"gopkg.in/mgutz/dat.v2/kvs"
	"gopkg.in/mgutz/dat.v2/postgres"
)

var logger log.Logger

// LogQueriesThreshold is the threshold for logging "slow" queries
var LogQueriesThreshold time.Duration

// PendingTransactionsTimeout is the timeout for pending transactions
// (used when the strict mode is enabled)
var PendingTransactionsTimeout = 1 * time.Minute

func init() {
	dat.Dialect = postgres.New()
	logger = log.New("dat:sqlx")
}

// Cache caches query results.
var Cache kvs.KeyValueStore

// SetCache sets this runner's cache. The default cache is in-memory
// based. See cache.MemoryKeyValueStore.
func SetCache(store kvs.KeyValueStore) {
	Cache = store
}

// MustPing pings a database with an exponential backoff. The
// function panics if the database cannot be pinged after 15 minutes
func MustPing(db *sql.DB) {
	var err error
	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)

	// Ticks will continue to arrive when the previous operation is still running,
	// so operations that take a while to fail could run in quick succession.
	for range ticker.C {
		if err = db.Ping(); err != nil {
			logger.Info("pinging database...", err.Error())
			continue
		}

		ticker.Stop()
		return
	}

	panic("Could not ping database!")
}
