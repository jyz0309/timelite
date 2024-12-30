package query

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/stretchr/testify/require"
)

func mockTSDB(t time.Time) error {
	db, _ := tsdb.Open("././storage", nil, nil, tsdb.DefaultOptions(), nil)
	appender := db.Appender(context.Background())
	for i := 0; i < 100000; i++ {
		appender.Append(0, labels.Labels{
			{Name: "__name__", Value: "my_counter"},
		}, t.UnixMilli()+int64(i), float64(rand.Float64()))
	}
	return appender.Commit()
}

func TestQuery(t *testing.T) {
	ts := time.Now()
	err := mockTSDB(ts)
	require.Nil(t, err)
	querier, err := NewQuerier("././storage")
	require.Nil(t, err)
	rs, timestamps, err := querier.NewRangeQuery(querier.ctx, "my_counter{}", ts.Truncate(5*time.Minute), ts.Add(5*time.Minute), 15*time.Second)
	require.Nil(t, err)
	fmt.Println(rs)
	fmt.Println(timestamps)
}
