package schemaregistry

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/hamba/avro/v2"
	"golang.org/x/sync/singleflight"
)

type LocalRegistry struct {
	realRegClient *Client
	subjects      map[string]int
	avroSchemas   map[int]avro.Schema
	mu            sync.RWMutex
	sf            singleflight.Group
}

func NewLocalRegistry(ctx context.Context, realClient *Client) (*LocalRegistry, error) {
	schemas, err := parseSchemasFS()
	if err != nil {
		return nil, err
	}

	subjects := make(map[string]int)
	avroSchemas := make(map[int]avro.Schema, 100)

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
		realRegClient: realClient,
		subjects:      subjects,
		avroSchemas:   avroSchemas,
		mu:            sync.RWMutex{},
		sf:            singleflight.Group{},
	}, nil
}

func (reg *LocalRegistry) GetSchemaIDBySubject(subject string) (int, error) {
	subjectID, ok := reg.subjects[subject]
	if !ok {
		return 0, fmt.Errorf("subject %s not found in subjects %v", subject, reg.subjects)
	}
	return subjectID, nil
}

func (reg *LocalRegistry) GetSchemaByID(id int) (avro.Schema, error) {
	schema, ok := reg.avroSchemas[id]
	if !ok {
		return nil, fmt.Errorf("schema with id %d not found in avroSchemas %v", id, reg.avroSchemas)
	}
	return schema, nil
}

// GetRemoteSchemaByID returns parsed schema by schema id from local cache
// or fetches it via schema registry client and saves it parsed
func (reg *LocalRegistry) GetRemoteSchemaByID(ctx context.Context, schemaID int) (avro.Schema, error) {
	reg.mu.RLock()
	schema, ok := reg.avroSchemas[schemaID]
	if ok {
		reg.mu.RUnlock()
		return schema, nil
	}
	reg.mu.RUnlock()

	val, err, _ := reg.sf.Do(strconv.Itoa(schemaID), func() (any, error) {
		schemaRaw, err := reg.realRegClient.GetSchemaByID(ctx, schemaID)
		if err != nil {
			return nil, err
		}

		avroSchema, err := avro.Parse(schemaRaw)
		if err != nil {
			return nil, fmt.Errorf("local reg get schema by id: %w", err)
		}

		reg.mu.Lock()
		reg.avroSchemas[schemaID] = avroSchema
		reg.mu.Unlock()

		return avroSchema, nil
	})

	if err != nil {
		return nil, err
	}

	c, _ := val.(avro.Schema)
	return c, nil
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
