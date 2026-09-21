package validator

import (
	"context"
	"errors"
	"github.com/go-semantic-engine-v2/application/domain/external"
)

// Schema struct defines a validation schema for product requests.
type Schema struct {
    Validate func(context.Context, any) error
}

// Use in ProductAdd and ProductPut.
func (s *Schema)VectorSearchSchema() Schema {
    return Schema{
        Validate: func(ctx context.Context, data any) error {
            
			req, ok := data.(external.VectorSearchRequest)
			if !ok {
                return errors.New("schema validation failed ! Please check the request body and try again.")
            }

            if req.Text == "" {
                return errors.New("schema validation failed ! field text is mandatory")
            }

            return nil
        },
    }
}
