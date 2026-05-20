package schemaregistry

import "embed"

//go:embed avroschemas/*.avsc
var schemasFS embed.FS
