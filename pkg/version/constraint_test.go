package version

import (
	"testing"

	"github.com/aquasecurity/go-version/pkg/part"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConstraints(t *testing.T) {
	tests := []struct {
		input   string
		opts    []ConstraintOption
		want    Constraints
		wantErr bool
	}{
		{
			">= 1.1",
			nil,
			Constraints{
				constraints: [][]constraint{{{
					version: Version{
						major: part.NewPart("1"),
						minor: part.NewPart("1"),
						patch: part.NewEmpty(true),
						revision: part.Parts{
							part.Any(true),
						},
						original: ">= 1.1",
					},
					operator: constraintGreaterThanEqual,
					original: ">= 1.1",
				}}},
			},
			false,
		},
		{
			">= 1.1",
			[]ConstraintOption{
				WithZeroPadding(true),
			},
			Constraints{
				constraints: [][]constraint{{{
					version: Version{
						major:    part.NewPart("1"),
						minor:    part.NewPart("1"),
						patch:    part.NewEmpty(false),
						revision: nil,
						original: ">= 1.1",
					},
					operator: constraintGreaterThanEqual,
					original: ">= 1.1",
				}}},
			},
			false,
		},
		{
			">40.50.60, < 50.70",
			nil,
			Constraints{
				constraints: [][]constraint{{{
					version: Version{
						major:    part.NewPart("40"),
						minor:    part.NewPart("50"),
						patch:    part.NewPart("60"),
						revision: nil,
						original: ">40.50.60",
					},
					operator: constraintGreaterThan,
					original: ">40.50.60",
				}, {
					version: Version{
						major: part.NewPart("50"),
						minor: part.NewPart("70"),
						patch: part.NewEmpty(true),
						revision: part.Parts{
							part.Any(true),
						},
						original: "< 50.70",
					},
					operator: constraintLessThan,
					original: "< 50.70",
				}}},
			},
			false,
		},
		{
			">= bar",
			nil,
			Constraints{},
			true,
		},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(test.input, func(tt *testing.T) {
			tt.Parallel()
			got, err := NewConstraints(test.input, test.opts...)
			if test.wantErr {
				assert.NotNil(tt, err)
			} else {
				assert.NoError(tt, err)
				assert.Len(tt, got.constraints, len(test.want.constraints))
				for i, c := range got.constraints {
					for j, cc := range c {
						assert.Equal(tt, test.want.constraints[i][j].version, cc.version)
						assert.Equal(tt, test.want.constraints[i][j].original, cc.original)
					}
				}
			}
		})
	}
}

func TestConstraints_Check(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		// Equal: =
		{"=2.0.0", "1.2.3", false},
		{"=2.0.0", "2.0.0", true},
		{"=2.0", "1.2.3", false},
		{"=2.0", "2.0.0", true},
		{"=2.0", "2.0.1", true},
		{"=0", "1.0.0", false},
		// Equal: ==
		{"== 2.0.0", "1.2.3", false},
		{"==2.0.0", "2.0.0", true},
		{"== 2.0", "1.2.3", false},
		{"==2.0", "2.0.0", true},
		{"== 2.0", "2.0.1", true},
		{"==0", "1.0.0", false},
		// Equal without "="
		{"4.1", "4.1.0", true},
		{"2", "1.0.0", false},
		{"2", "3.4.5", false},
		{"2", "2.1.1", true},
		{"2.1", "2.1.1", true},
		{"2.1", "2.2.1", false},
		// Not equal
		{"!=4.1.0", "4.1.0", false},
		{"!=4.1.0", "4.1.1", true},
		{"!=4.1", "4.1.0", false},
		{"!=4.1", "4.1.1", false},
		{"!=4.1", "5.1.0", true},
		// Less than
		{"<11", "0.1.0", true},
		{"<11", "11.1.0", false},
		{"<1.1", "0.1.0", true},
		{"<1.1", "1.1.0", false},
		{"<1.1", "1.1.1", false},
		// Less than or equal
		{"<=11", "1.2.3", true},
		{"<=11", "12.2.3", false},
		{"<=11", "11.2.3", true},
		{"<=1.1", "1.2.3", false},
		{"<=1.1", "0.1.0", true},
		{"<=1.1", "1.1.0", true},
		{"<=1.1", "1.1.1", true},
		// Greater than
		{">1.1", "4.1.0", true},
		{">1.1", "1.1.0", false},
		{">0", "0.0.0", false},
		{">0", "1.0.0", true},
		{">11", "11.1.0", false},
		{">11.1", "11.1.0", false},
		{">11.1", "11.1.1", false},
		{">11.1", "11.2.1", true},
		// Greater than or equal
		{">=11", "11.1.2", true},
		{">=11.1", "11.1.2", true},
		{">=11.1", "11.0.2", false},
		{">=1.1", "4.1.0", true},
		{">=1.1", "1.1.0", true},
		{">=1.1", "0.0.9", false},
		{">=0", "0.0.0", true},
		{"=0", "1.0.0", false},
		// Asterisk
		{"*", "1.0.0", true},
		{"*", "4.5.6", true},
		{"2.*", "1.0.0", false},
		{"2.*", "3.4.5", false},
		{"2.*", "2.1.1", true},
		{"2.1.*", "2.1.1", true},
		{"2.1.*", "2.2.1", false},
		// Empty
		{"", "1.0.0", true},
		{"", "4.5.6", true},
		{"2", "1.0.0", false},
		{"2", "3.4.5", false},
		{"2", "2.1.1", true},
		{"2.1", "2.1.1", true},
		{"2.1", "2.2.1", false},
		// Tilde
		{"~1.2.3", "1.2.4", true},
		{"~1.2.3", "1.3.4", false},
		{"~1.2", "1.2.4", true},
		{"~1.2", "1.3.4", false},
		{"~1", "1.2.4", true},
		{"~1", "2.3.4", false},
		{"~0.2.3", "0.2.5", true},
		{"~0.2.3", "0.3.5", false},
		{"^1.2.3", "1.8.9", true},
		{"^1.2.3", "2.8.9", false},
		{"^1.2.3", "1.2.1", false},
		{"^1.1.0", "2.1.0", false},
		{"^1.2.0", "2.2.1", false},
		{"^1.2", "1.8.9", true},
		{"^1.2", "2.8.9", false},
		{"^1", "1.8.9", true},
		{"^1", "2.8.9", false},
		{"^0.2.3", "0.2.5", true},
		{"^0.2.3", "0.5.6", false},
		{"^0.2", "0.2.5", true},
		{"^0.2", "0.5.6", false},
		{"^0.0.3", "0.0.3", true},
		{"^0.0.3", "0.0.4", false},
		{"^0.0", "0.0.3", true},
		{"^0.0", "0.1.4", false},
		{"^0.0", "1.0.4", false},
		{"^0", "0.2.3", true},
		{"^0", "1.1.4", false},
		// revision: Not equal
		{"!=4.1", "5.1.0-1", true},
		{"!=4.1-1", "4.1.0", false},
		// revision: Greater than
		{">0", "0.0.1-1", false},
		{">0.0", "0.0.1-1", false},
		{">0-0", "0.0.1-1", false},
		{">0.0-0", "0.0.1-1", false},
		{">0", "0.0.0-1", false},
		{">0-0", "0.0.0-1", false},
		{">0.0.0-0", "0.0.0-1", true},
		{">1.2.3-1", "1.2.3-2", true},
		{">1.2.3-1", "1.3.3-2", true},
		// revision: Less than
		{"<0", "0.0.0-1", false},
		{"<0-2", "0.0.0-1", true},
		// revision: Greater than or equal
		{">=0", "0.0.1-1", true},
		{">=0.0", "0.0.1-1", true},
		{">=0-0", "0.0.1-1", true},
		{">=0.0-0", "0.0.1", true},
		{">=0", "0.0.0-1", true},
		{">=0-0", "0.0.0-1", true},
		{">=0.0.0", "0.0.0-1", true},
		{">=0.0.0-0", "0.0.0-1", true},
		{">=0.0.0-2", "0.0.0-1", false},
		{">=0.0.0-0", "1.2.3", true},
		{">=0.0.0-0", "3.4.5-1", true},
		// revision: Asterisk
		{"*", "1.2.3-1", true},
		// revision: Empty
		{"", "1.2.3-1", true},
		// revision: Tilde
		{"~1.2.3-1", "1.2.3-2", true},
		{"~1.2.3-1", "1.2.4-1", true},
		{"~1.2.3-1", "1.3.4-1", false},
		// revision: Caret
		{"^1.2.0", "1.2.1-1", true},
		{"^1.2.0-0", "1.2.1-1", true},
		{"^1.2.0-0", "1.2.1-0", true},
		{"^1.2.0-2", "1.2.0-1", false},
		{"^0.2.3-3", "0.2.3-4", true},
		{"^0.2.3-3", "0.2.4-3", true},
		{"^0.2.3-3", "0.3.4-3", false},
		{"^0.2.3-3", "0.2.3-3", true},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(test.constraint, func(tt *testing.T) {
			tt.Parallel()
			c, err := NewConstraints(test.constraint)
			require.NoError(tt, err)

			v, err := Parse(test.version)
			require.NoError(tt, err)

			got := c.Check(v)
			assert.Equal(tt, test.want, got)
		})
	}
}

func TestConstraints_CheckWithZeroPadding(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		// Equal
		{"=2.0.0", "1.2.3", false},
		{"=2.0.0", "2.0.0", true},
		// Not equal
		{"!=4.1.0", "4.1.0", false},
		{"!=4.1.0", "4.1.1", true},
		// Less than
		{"<0.0.5", "0.1.0", false},
		{"<1.0.0", "0.1.0", true},
		// Less than or equal
		{"<=0.2.3", "1.2.3", false},
		{"<=1.2.3", "1.2.3", true},
		// Greater than
		{">5.0.0", "4.1.0", false},
		{">4.0.0", "4.1.0", true},
		// Greater than or equal
		{">=11.1.3", "11.1.2", false},
		{">=11.1.2", "11.1.2", true},
		// Asterisk
		{"*", "1.0.0", true},
		{"*", "4.5.6", true},
		{"2.*", "1.0.0", false},
		{"2.*", "3.4.5", false},
		{"2.*", "2.1.1", true},
		{"2.1.*", "2.1.1", true},
		{"2.1.*", "2.2.1", false},
		// Empty
		{"", "1.0.0", true},
		{"", "4.5.6", true},
		{"2", "1.0.0", false},
		{"2", "3.4.5", false},
		{"2", "2.1.1", false},
		{"2.1", "2.1.1", false},
		{"2.1", "2.2.1", false},
		// Tilde
		{"~1.2.3", "1.2.4", true},
		{"~1.2.3", "1.3.4", false},
		{"~1.2", "1.2.4", true},
		{"~1.2", "1.3.4", false},
		{"~1", "1.2.4", true},
		{"~1", "2.3.4", false},
		{"~0.2.3", "0.2.5", true},
		{"~0.2.3", "0.3.5", false},
		{"~1.2.3-beta.2", "1.2.3-beta.4", true},
		// Caret
		{"^1.2.3", "1.8.9", true},
		{"^1.2.3", "2.8.9", false},
		{"^1.2.3", "1.2.1", false},
		{"^1.1.0", "2.1.0", false},
		{"^1.2.0", "2.2.1", false},
		{"^1.2", "1.8.9", true},
		{"^1.2", "2.8.9", false},
		{"^1", "1.8.9", true},
		{"^1", "2.8.9", false},
		{"^0.2.3", "0.2.5", true},
		{"^0.2.3", "0.5.6", false},
		{"^0.2", "0.2.5", true},
		{"^0.2", "0.5.6", false},
		{"^0.0.3", "0.0.3", true},
		{"^0.0.3", "0.0.4", false},
		{"^0.0", "0.0.3", true},
		{"^0.0", "0.1.4", false},
		{"^0.0", "1.0.4", false},
		{"^0", "0.2.3", true},
		{"^0", "1.1.4", false},
		// revision: Equal
		{"=4.1", "4.1.0-1", false},
		{"=4.1-1", "4.1.0-1", true},
		{"== 4.1", "4.1.0-1", false},
		{"==4.1-1", "4.1.0-1", true},
		// revision: Not equal
		{"!=4.1", "5.1.0-1", true},
		{"!=4.1-1", "4.1.0", true},
		// revision: Greater than
		{">0", "0.0.1-1", true},
		{">0.0", "0.0.1-1", true},
		{">0-0", "0.0.1-1", true},
		{">0.0-0", "0.0.1-1", true},
		{">0", "0.0.0-1", true},
		{">0-0", "0.0.0-1", true},
		{">0.0.0-0", "0.0.0-1", true},
		{">1.2.3-1", "1.2.3-1", false},
		{">1.2.3-1", "1.2.3-2", true},
		{">1.2.3-1", "1.3.3-2", true},
		// revision: Less than
		{"<0", "0.0.0-1", false},
		{"<0-2", "0.0.0-1", true},
		{"<0", "1.0.0-1", false},
		{"<1", "1.0.0-1", false},
		{"<2", "1.0.0-1", true},
		// revision: Greater than or equal
		{">=0", "0.0.1-1", true},
		{">=0.0", "0.0.1-1", true},
		{">=0-0", "0.0.1-1", true},
		{">=0.0-0", "0.0.1-1", true},
		{">=0", "0.0.0-1", true},
		{">=0-0", "0.0.0-1", true},
		{">=0.0.0-1", "0.0.0-1", true},
		{">=0.0.0-2", "0.0.0-1", false},
		{">=0.0.0-0", "1.2.3", true},
		{">=0.0.0-0", "3.4.5-1", true},
		// revision: Asterisk
		{"*", "1.2.3-1", true},
		// revision: Empty
		{"", "1.2.3-1", true},
		// revision: Tilde
		{"~1.2.3-1", "1.2.3-2", true},
		{"~1.2.3-1", "1.2.4-1", true},
		{"~1.2.3-1", "1.3.4-1", false},
		// revision: Caret
		{"^1.2.0", "1.2.1-alpha.1", true},
		{"^1.2.0-0", "1.2.1-1", true},
		{"^1.2.0-0", "1.2.1-0", true},
		{"^1.2.0-2", "1.2.0-1", false},
		{"^0.2.3-2", "0.2.3-4", true},
		{"^0.2.3-2", "0.2.4-2", true},
		{"^0.2.3-2", "0.3.4-2", false},
		{"^0.2.3-2", "0.2.3-2", true},
	}
	t.Parallel()
	for _, testToRun := range tests {
		test := testToRun
		t.Run(test.constraint, func(tt *testing.T) {
			tt.Parallel()
			c, err := NewConstraints(test.constraint, WithZeroPadding(true))
			require.NoError(tt, err)

			v, err := Parse(test.version)
			require.NoError(tt, err)

			got := c.Check(v)
			assert.Equal(tt, test.want, got)
		})
	}
}
