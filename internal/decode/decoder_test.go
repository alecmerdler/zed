package decode

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authzed/spicedb/pkg/validationfile"
)

func TestRewriteURL(t *testing.T) {
	tests := []struct {
		name string
		in   url.URL
		out  url.URL
	}{
		{
			name: "gist",
			in: url.URL{
				Scheme: "https",
				Host:   "gist.github.com",
				Path:   "/ecordell/9e2110ac4a1292b899784ed809d44b8f",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "gist.githubusercontent.com",
				Path:   "/ecordell/9e2110ac4a1292b899784ed809d44b8f/raw",
			},
		},
		{
			name: "playground schema",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/schema",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground relationships",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/relationships",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground assertions",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/assertions",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground expected",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/expected",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "pastebin",
			in: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/LuCwwBwU",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
		},
		{
			name: "pastebin raw",
			in: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
		},
		{
			name: "direct",
			in: url.URL{
				Scheme: "https",
				Host:   "somethingelse.com",
				Path:   "/any/path",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "somethingelse.com",
				Path:   "/any/path",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rewriteURL(&tt.in)
			require.EqualValues(t, tt.out, tt.in)
		})
	}
}

func TestUnmarshalAsYAMLOrSchema(t *testing.T) {
	tests := []struct {
		name         string
		in           []byte
		isOnlySchema bool
		outSchema    string
		wantErr      bool
	}{
		{
			name: "valid yaml",
			in: []byte(`
schema:
  definition user {}
`),
			outSchema:    `definition user {}`,
			isOnlySchema: false,
			wantErr:      false,
		},
		{
			name:         "valid schema",
			in:           []byte(`definition user {}`),
			isOnlySchema: true,
			outSchema:    `definition user {}`,
			wantErr:      false,
		},
		{
			name:         "invalid yaml",
			in:           []byte(`invalid yaml`),
			isOnlySchema: false,
			outSchema:    "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := validationfile.ValidationFile{}
			isOnlySchema, err := unmarshalAsYAMLOrSchema(tt.in, &block)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.isOnlySchema, isOnlySchema)
			if !tt.wantErr {
				require.Equal(t, tt.outSchema, block.Schema.Schema)
			}
		})
	}
}
