package schemaregistry

import "embed"

//go:embed avroschemas/*.avsc
var schemaFS embed.FS
