package repository

import (
	"fmt"
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"

	"github.com/go-semantic-engine-v2/application/tracing"
	"github.com/go-semantic-engine-v2/application/domain/entity"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/codes"
)

type SemanticRepository struct{
	dbConnector connector.IDatabaseConnector
}

type ISemanticRepository interface {
	SearchVector(ctx context.Context, vector entity.VectorSearch) (*entity.VectorSearch, error)
}

func NewSemanticRepository(dbConnector connector.IDatabaseConnector) ISemanticRepository {
	logger.InfoOutCtx("initializing semantic repository SUCCESSFULLY")
	
	return &SemanticRepository{
		dbConnector: dbConnector,
	}
}

func (r *SemanticRepository) SearchVector(ctx context.Context, vector entity.VectorSearch) (res_vector *entity.VectorSearch, err error) {
	logger.Info(ctx, "semantic repository VectorSearch called")

	// Tracing and metrics
	ctx, span := tracing.CustomStartSpanCtx(ctx, "semanticRepository.VectorSearch", trace.SpanKindInternal)
	defer span.End()

	meter := otel.Meter("go-semantic-engine-v2.repository")
    counter, _ := meter.Int64Counter("db_custom_vector_search_requests_total")
    histogram, _ := meter.Float64Histogram("db_custom_vector_search_duration_seconds")
    start := time.Now()

    counter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("operation", "VectorSearch"),
    ))
	
	defer func() {
		if err != nil {
			span.RecordError(err) 
			span.SetStatus(codes.Error, err.Error())
			logger.Error(ctx, "semantic repository VectorSearch failed", zap.Error(err))
		}
        histogram.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
            attribute.String("operation", "VectorSearch"),
        ))
	}()

	connectorReader := r.dbConnector.Reader()

	// Convert vector to string for SQL query
	strVector := "["
	for i, v := range vector.Embedding {
		strVector += fmt.Sprintf("%f", v)
		if i < len(vector.Embedding)-1 {
			strVector += ","
		}
	}
	strVector += "]"

	query := `select  ce."text",
					  cer."type",
					  c.name,
					  ced.endpoint,
					  ced.uri
				from embedding e,
					capability_example_relation cer,
					capability c,
					capability_endpoint ced,
					capability_example ce
				where (e.vector <=> $1) < 0.8
				and e.fk_cap_exp_id = ce.id
				and cer.fk_cap_exp = ce.id
				and c.id = cer.fk_cap_id 
				and c.fk_svc_id = cer.fk_cap_svc_id 
				and ced.fk_cap_id = c.id 
				and ced.fk_cap_svc_id = c.fk_svc_id
				order by (e.vector <=> $1) asc
				limit 10`

    rows, err := connectorReader.Query(ctx, query, strVector)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
	res_vector = &entity.VectorSearch{Text: vector.Text}

	for rows.Next() {
		capability := entity.Capability{}
		err := rows.Scan(
			&capability.Text,
			&capability.Name,
			&capability.Type,
			&capability.Endpoint,
			&capability.Uri,
		)
		if err != nil {
			return nil, err
		}
		res_vector.Capabilities = append(res_vector.Capabilities, capability)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res_vector, nil
}