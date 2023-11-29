package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		version string
		wantErr bool
	}{
		{"1.2.3", false},
		{"1.2.3-alpha.01", true},
		{"1.2.3+test.01", false},
		{"1.2.3-alpha.-1", false},
		{"1.0", true},
		{"1", true},
		{"1.2.beta", true},
		{"foo", true},
		{"1.2-5", true},
		{"1.2-beta.5", true},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hypen", false},
		{"1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", true},
		{"1.2.2147483648", false},
		{"1.2147483648.3", false},
		{"2147483648.3.0", false},
		{"1.0.0-alpha", false},
		{"1.0.0-alpha.1", false},
		{"1.0.0-0.3.7", false},
		{"1.0.0-x.7.z.92", false},
		{"1.0.0-x-y-z.-", false},
		{"1.2.3.4", true},
		{"foo1.2.3", true},
		{"1.7rc2", true},
		{"1.0-", true},
		{"v1.2.3", true},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(test.version, func(tt *testing.T) {
			tt.Parallel()
			_, err := Parse(test.version)
			if test.wantErr {
				assert.NotNil(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.4.5", -1},
		{"2.2.3", "1.5.1", 1},
		{"2.2.3", "2.2.2", 1},
		{"1.0.0-1", "1.0.0-2", -1},
		{"1.0.0-2", "1.0.0-1", 1},
		{"1.2.3-1", "1.2.3-1", 0},
		{"1.2.3", "1.2.3-1", -1},
		{"1.2.3-1", "1.2.3", 1},
		{"1.2.3+foo", "1.2.3+bar", 0},
		{"1.2.3+foo", "1.2.3+bar", 0},
		{"1.2.0", "1.2.0-1+metadata", -1},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(fmt.Sprintf("%s vs %s", test.v1, test.v2), func(tt *testing.T) {
			tt.Parallel()
			v1, err := Parse(test.v1)
			require.NoError(tt, err, test.v1)

			v2, err := Parse(test.v2)
			require.NoError(tt, err, test.v2)

			assert.Equal(tt, test.expected, v1.Compare(v2))
		})
	}
}
