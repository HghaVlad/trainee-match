package schemaregistry

import (
	"context"
	"fmt"
	"io/fs"
	"strings"
)

type LocalRegistry struct {
	SchemaIDs map[string]int // subject -> schemaID

	schemas       map[string]string // subject -> raw schema string
	realRegClient *RealRegistryClient
}

// NewLocalRegistry loads, looks up schema ids of local schemas, saves them
func NewLocalRegistry(ctx context.Context, client *RealRegistryClient) (*LocalRegistry, error) {
	schemas, err := loadLocalSchemas()
	if err != nil {
		return nil, fmt.Errorf("new local schema registry: %w", err)
	}

	schemaIDs, err := lookUpSchemaIDs(ctx, client, schemas)
	if err != nil {
		return nil, fmt.Errorf("new local schema registry: %w", err)
	}

	return &LocalRegistry{
		SchemaIDs:     schemaIDs,
		schemas:       schemas,
		realRegClient: client,
	}, nil
}

// looks up schema ids of local schemas in real schema registry
func lookUpSchemaIDs(
	ctx context.Context,
	realRegClient *RealRegistryClient,
	schemas map[string]string,
) (map[string]int, error) {
	schemaIDs := make(map[string]int)

	for subject, schema := range schemas {
		schemaID, err := realRegClient.LookupSchemaID(ctx, subject, schema)
		if err != nil {
			return nil, fmt.Errorf("look up schema id for subject %s: %w", subject, err)
		}

		schemaIDs[subject] = schemaID
	}

	return schemaIDs, nil
}

// returns subject -> raw schema string
func loadLocalSchemas() (map[string]string, error) {
	schemas := make(map[string]string)

	err := fs.WalkDir(schemaFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if !endsWithAvsc(path) {
			return nil
		}

		data, err := fs.ReadFile(schemaFS, path)
		if err != nil {
			return err
		}

		subject := subjectFromPath(path)
		schemas[subject] = string(data)

		return nil
	})

	return schemas, err
}

func subjectFromPath(path string) string {
	// vacancy_published.avsc → vacancy_published-value
	base := path[strings.LastIndex(path, "/")+1:]
	name := strings.TrimSuffix(base, ".avsc")
	return name + "-value"
}

func endsWithAvsc(path string) bool {
	return strings.HasSuffix(path, ".avsc")
}
