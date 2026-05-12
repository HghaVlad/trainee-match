package schemaregistry

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hamba/avro/v2"
)

type LocalRegistry struct {
	client     *RealRegistryClient
	subjects   map[string]int      // subject -> schemaId
	parsedByID map[int]avro.Schema // shemaId -> schema
}

func NewLocalRegistry(ctx context.Context, client *RealRegistryClient) (*LocalRegistry, error) {
	schemas, err := loadAllSchemas()

	if err != nil {
		return nil, err
	}

	subjects := make(map[string]int)
	parsedByID := make(map[int]avro.Schema)

	for subject, schema := range schemas {
		parsedSchema, err := avro.Parse(schema)
		if err != nil {
			return nil, err
		}
		schemaID, err := client.LookupSchemaID(ctx, subject, schema)
		if err != nil {
			return nil, err
		}

		subjects[subject] = schemaID
		parsedByID[schemaID] = parsedSchema
	}

	return &LocalRegistry{
		subjects:   subjects,
		parsedByID: parsedByID,
	}, nil
}

func (reg *LocalRegistry) GetSchemaByID(id int) (avro.Schema, error) {
	schema, exists := reg.parsedByID[id]
	if !exists {
		return nil, errors.New("schema not found")
	}
	return schema, nil
}

func (reg *LocalRegistry) GetSchemaIDBySubject(subject string) (int, error) {
	if schema, exists := reg.subjects[subject]; exists {
		return schema, nil
	}
	return 0, errors.New("subject not found")
}

func loadAllSchemas() (map[string]string, error) {
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
	// user-created.avsc → user-created-value
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ".avsc")
	return name + "-value"
}

func endsWithAvsc(path string) bool {
	return strings.HasSuffix(path, ".avsc")
}
