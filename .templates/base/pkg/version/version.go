// Package version contains global variables for
// nolint: gochecknoglobals, gochecknoinits
package version

var (
	// RELEASE returns the release version.
	RELEASE = "UNKNOWN"
	{{[- if .API.Enabled ]}}
	// API returns the API version.
	API = "UNKNOWN"
	{{[- end ]}}
	// DATE returns the release date.
	DATE = "UNKNOWN"
	// REPO returns the git repository URL.
	REPO = "UNKNOWN"
	// COMMIT returns the short sha from git.
	COMMIT = "UNKNOWN"
	// BRANCH returns deployed git brunch.
	BRANCH = "UNKNOWN"
	// DESC returns the service description.
	DESC = "{{[ .Description ]}}"
)
