package labelselector

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	cases := []struct {
		name             string
		input            string
		expectedSelector LabelSelector
		expectedError    error
	}{
		{
			name:  "one equals test",
			input: `foo=bar`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationEquals,
						Value:     "bar",
					},
				},
			},
		},
		{
			name:  "one lower than test",
			input: `foo < 5`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationLowerThan,
						Value:     "5",
					},
				},
			},
		},
		{
			name:  "one lower than equal test",
			input: `foo <= 5`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationLowerThanEqual,
						Value:     "5",
					},
				},
			},
		},
		{
			name:  "one greater than test",
			input: `foo > 5`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationGreaterThan,
						Value:     "5",
					},
				},
			},
		},
		{
			name:  "one greater than equal test",
			input: `foo >= 5`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationGreaterThanEqual,
						Value:     "5",
					},
				},
			},
		},
		{
			name:  "one equals test with '=='",
			input: `foo == bar`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationEquals,
						Value:     "bar",
					},
				},
			},
		},
		{
			name:  "one equals test with extra whitespaces",
			input: ` foo  = bar   `,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationEquals,
						Value:     "bar",
					},
				},
			},
		},
		{
			name:  "one not equals test",
			input: `foo != bar`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Key:       "foo",
						Operation: OperationNotEquals,
						Value:     "bar",
					},
				},
			},
		},
		{
			name:  "one existance test",
			input: `foo`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationExists,
						Key:       "foo",
					},
				},
			},
		},
		{
			name:  "one existance test with quoted name",
			input: `"foo bar"`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationExists,
						Key:       "foo bar",
					},
				},
			},
		},
		{
			name:  "one non-existance test",
			input: `!foo`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationNotExists,
						Key:       "foo",
					},
				},
			},
		},
		{
			name:  "one in test",
			input: `foo in (a, b, c)`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationIn,
						Key:       "foo",
						Values:    []string{"a", "b", "c"},
					},
				},
			},
		},
		{
			name:  "one notin test",
			input: `foo notin (a, b, c)`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationNotIn,
						Key:       "foo",
						Values:    []string{"a", "b", "c"},
					},
				},
			},
		},
		{
			name:  "multiple checks",
			input: `foo, bar, !baz, bla=blub`,
			expectedSelector: LabelSelector{
				Requirements: []Requirement{
					{
						Operation: OperationExists,
						Key:       "foo",
					},
					{
						Operation: OperationExists,
						Key:       "bar",
					},
					{
						Operation: OperationNotExists,
						Key:       "baz",
					},
					{
						Operation: OperationEquals,
						Key:       "bla",
						Value:     "blub",
					},
				},
			},
		},
		{
			name:          "illegal characters are rejected",
			input:         `❤`,
			expectedError: errors.New("illegal token"),
		},
		{
			name:          "not-existance test needs argument",
			input:         `!`,
			expectedError: errors.New("expect identifier after exclamation mark"),
		},
		{
			name:          "identifier needs operator if any",
			input:         `foo bar`,
			expectedError: errors.New("unexpected token 'bar'"),
		},
		{
			name:          "equal operator needs argument",
			input:         `foo=`,
			expectedError: errors.New("expect identifier after equal operator"),
		},
		{
			name:          "lower than operator needs argument",
			input:         `foo<`,
			expectedError: errors.New("expect identifier after < operator"),
		},
		{
			name:          "lower than equal operator needs argument",
			input:         `foo<=`,
			expectedError: errors.New("expect identifier after <= operator"),
		},
		{
			name:          "greater than operator needs argument",
			input:         `foo>`,
			expectedError: errors.New("expect identifier after > operator"),
		},
		{
			name:          "greater than equal operator needs argument",
			input:         `foo>=`,
			expectedError: errors.New("expect identifier after >= operator"),
		},
		{
			name:          "not equal operator needs argument",
			input:         `foo!=`,
			expectedError: errors.New("expect identifier after not-equal operator"),
		},
		{
			name:          "in operator needs argument",
			input:         `foo in`,
			expectedError: errors.New("expect opening bracket after in operator"),
		},
		{
			name:          "notin operator needs argument",
			input:         `foo notin`,
			expectedError: errors.New("expect opening bracket after in operator"),
		},
		{
			name:          "in operator needs properly formatted argument list",
			input:         `foo in (]`,
			expectedError: errors.New("unexpected token in value list (])"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			selector, err := ParseString(tc.input)
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedSelector, selector)
		})
	}
}
