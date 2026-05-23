package schemaregistry

import "embed"

//go:embed avroschemas/*
var schemaFS embed.FS
