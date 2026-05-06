package schemaregistry

import (
	"context"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"sync"

	"github.com/hamba/avro/v2"
	"golang.org/x/sync/singleflight"
)

type LocalRegistry struct {
	schemaIDs map[string]int // subject -> schemaID

	schemas    map[string]string   // subject -> raw schema string
	parsedByID map[int]avro.Schema // schemaID -> parsed schema

	realRegClient *RealRegistryClient

	mu sync.RWMutex
	sf singleflight.Group
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
		schemaIDs:     schemaIDs,
		schemas:       schemas,
		parsedByID:    make(map[int]avro.Schema),
		realRegClient: client,
		mu:            sync.RWMutex{},
		sf:            singleflight.Group{},
	}, nil
}

// GetSchemaByID returns parsed schema by schema id from local cache
// or fetches it via schema registry client and saves it parsed
func (r *LocalRegistry) GetSchemaByID(ctx context.Context, schemaID int) (avro.Schema, error) {
	r.mu.RLock()
	schema, ok := r.parsedByID[schemaID]
	if ok {
		r.mu.RUnlock()
		return schema, nil
	}
	r.mu.RUnlock()

	val, err, _ := r.sf.Do(strconv.Itoa(schemaID), func() (any, error) {
		schemaRaw, err := r.realRegClient.GetSchemaByID(ctx, schemaID)
		if err != nil {
			return nil, err
		}

		avroSchema, err := avro.Parse(schemaRaw)
		if err != nil {
			return nil, fmt.Errorf("local reg get schema by id: %w", err)
		}

		r.mu.Lock()
		r.parsedByID[schemaID] = avroSchema
		r.mu.Unlock()

		return avroSchema, nil
	})

	if err != nil {
		return nil, err
	}

	c, _ := val.(avro.Schema)
	return c, nil
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
	// vacancy-published.avsc → vacancy-published-value
	base := path[strings.LastIndex(path, "/")+1:]
	name := strings.TrimSuffix(base, ".avsc")
	return name + "-value"
}

func endsWithAvsc(path string) bool {
	return strings.HasSuffix(path, ".avsc")
}
