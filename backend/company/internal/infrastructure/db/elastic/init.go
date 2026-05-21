package elastic

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

const vacancyIndex = "vacancies"

const vacancyIndexMapping = `{
  "mappings": {
    "dynamic": "strict",

    "properties": {
      "id": {
        "type": "keyword"
      },

      "company_id": {
        "type": "keyword"
      },

      "company_name": {
        "type": "text"
      },

      "title": {
        "type": "text"
      },

      "description": {
        "type": "text"
      },

      "work_format": {
        "type": "keyword"
      },

      "city": {
        "type": "keyword"
      },

      "employment_type": {
        "type": "keyword"
      },

      "status": {
        "type": "keyword"
      },

      "moderation_status": {
        "type": "keyword"
      },

      "flexible_schedule": {
        "type": "boolean"
      },

      "is_paid": {
        "type": "boolean"
      },

      "internship_to_offer": {
        "type": "boolean"
      },

      "salary_from": {
        "type": "integer"
      },

      "salary_to": {
        "type": "integer"
      },

      "hours_per_week_from": {
        "type": "integer"
      },

      "hours_per_week_to": {
        "type": "integer"
      },

      "duration_from_days": {
        "type": "integer"
      },

      "duration_to_days": {
        "type": "integer"
      },

      "published_at": {
        "type": "date"
      },

      "created_at": {
        "type": "date"
      }
    }
  }
}`

func Init(ctx context.Context, cl *elasticsearch.Client) error {
	exists, err := cl.Indices.Exists([]string{vacancyIndex})
	if err != nil {
		return err
	}

	if exists.StatusCode == http.StatusOK {
		return nil
	}

	res, err := cl.Indices.Create(
		vacancyIndex,
		cl.Indices.Create.WithContext(ctx),
		cl.Indices.Create.WithBody(strings.NewReader(vacancyIndexMapping)),
	)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return errors.New("elastic create index error")
	}

	return nil
}
