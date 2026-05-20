package schemaregistry

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/hamba/avro/v2"
)

type LocalRegistry struct {
	realClient  *Client
	subjects    map[string]int
	avroSchemas map[int]avro.Schema
}

func NewLocalRegistry(ctx context.Context, realClient *Client) (*LocalRegistry, error) {
	schemas, err := parseSchemasFS()
	if err != nil {
		return nil, err
	}

	subjects := make(map[string]int)
	avroSchemas := make(map[int]avro.Schema)

	for subject, schema := range schemas {
		id, err := realClient.LookUpSchemaID(ctx, subject, schema)
		if err != nil {
			return nil, err
		}
		subjects[subject] = id

		avroSchema, err := avro.Parse(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to parse avro schema: %w", err)
		}
		avroSchemas[id] = avroSchema
	}
	return &LocalRegistry{
		realClient:  realClient,
		subjects:    subjects,
		avroSchemas: avroSchemas,
	}, nil
}

func (reg *LocalRegistry) GetSchemaIDBySubject(subject string) (int, error) {
	subjectId, ok := reg.subjects[subject]
	if !ok {
		return 0, fmt.Errorf("subject %s not found in subjects %v", subject, reg.subjects)
	}
	return subjectId, nil
}

func (reg *LocalRegistry) GetSchemaByID(id int) (avro.Schema, error) {
	schema, ok := reg.avroSchemas[id]
	if !ok {
		return nil, fmt.Errorf("schema with id %d not found in avroSchemas %v", id, reg.avroSchemas)
	}
	return schema, nil
}

func parseSchemasFS() (map[string]string, error) {
	schemas := make(map[string]string)

	files, err := schemasFS.ReadDir("avroschemas")
	if err != nil {
		return nil, fmt.Errorf("failed to read schemas dir: %w", err)
	}
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), ".avsc") + "-value"
		content, err := schemasFS.ReadFile(path.Join("avroschemas", f.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read schemas file: %w", err)
		}
		schemas[name] = string(content)
	}
	return schemas, nil
}
