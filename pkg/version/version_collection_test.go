package version

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     []string
	}{
		{
			name: "happy path",
			versions: []string{
				"1.1.1",
				"1.0.0",
				"1.2.0",
				"2.0.0",
				"0.7.1",
			},
			want: []string{
				"0.7.1",
				"1.0.0",
				"1.1.1",
				"1.2.0",
				"2.0.0",
			},
		},
		{
			name: "revisions",
			versions: []string{
				"1.0.0-1.1",
				"1.0.0-3.1",
				"1.0.0",
				"1.0.0-2.2",
				"1.0.0-2.11",
				"1.0.0-1",
				"1.0.0-1.2",
				"1.0.0-2",
			},
			want: []string{
				"1.0.0",
				"1.0.0-1",
				"1.0.0-1.1",
				"1.0.0-1.2",
				"1.0.0-2",
				"1.0.0-2.2",
				"1.0.0-2.11",
				"1.0.0-3.1",
			},
		},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			versions := make(Collection, len(test.versions))
			for i, raw := range test.versions {
				v, err := Parse(raw)
				require.NoError(tt, err)
				versions[i] = v
			}

			sort.Sort(Collection(versions))

			got := make([]string, len(versions))
			for i, v := range versions {
				got[i] = v.String()
			}

			assert.Equal(tt, test.want, got)
		})
	}
}
