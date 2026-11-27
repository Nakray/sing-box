package ccm

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const sampleCredentials = `{"claudeAiOauth":{"accessToken":"access-test","refreshToken":"refresh-test","expiresAt":1893456000000}}`

func writeCredentialsFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".credentials.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestReadCredentialsFromFile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "normal file",
			content: sampleCredentials,
		},
		{
			name:    "exactly at limit",
			content: sampleCredentials + strings.Repeat(" ", maxCredentialFileSize-len(sampleCredentials)),
		},
		{
			name:    "over limit",
			content: strings.Repeat(" ", maxCredentialFileSize+1),
			wantErr: "credentials file too large",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			path := writeCredentialsFile(t, testCase.content)
			credentials, err := readCredentialsFromFile(context.Background(), path)
			if testCase.wantErr != "" {
				require.ErrorContains(t, err, testCase.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "access-test", credentials.AccessToken)
		})
	}
}
