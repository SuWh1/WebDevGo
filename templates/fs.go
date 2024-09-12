package templates

import "embed"

//go:embed *
var FS embed.FS // embedded filesystem
// it embeds all files in that directory to binary to make it easier serve templates
