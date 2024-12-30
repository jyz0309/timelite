package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListDashboard(t *testing.T) {
	dashboard := ListDashboards()
	require.Equal(t, 1, len(dashboard))
}
