package usecase

import (
	"context"
		
	"go.uber.org/zap"
	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-semantic-engine-v2/application/tracing"
    "github.com/go-semantic-engine-v2/application/domain/entity"
    "github.com/go-semantic-engine-v2/application/infrastructure/repository"

	"go.opentelemetry.io/otel/trace"
)

type SemanticUseCase struct {
    semanticRepository 	repository.ISemanticRepository
}

type ISemanticUseCase interface {
	SearchVector(ctx context.Context, vector entity.VectorSearch) (*entity.VectorSearch, error)
	IntentDecompose(ctx context.Context, intent entity.IntentDecompose) (*entity.IntentDecompose, error)
}

func NewSemanticUseCase(semanticRepository repository.ISemanticRepository) ISemanticUseCase {
	logger.InfoOutCtx("initializing semantic usecase SUCCESSFULLY")
	
	return &SemanticUseCase{
		semanticRepository: 	semanticRepository,
	}
}

// Login handles the login use case, generating an OAuth token for the given login credentials.
func (uc *SemanticUseCase) SearchVector(ctx context.Context, vector entity.VectorSearch) (*entity.VectorSearch, error) {
	logger.Info(ctx, "semantic usecase SearchVector called with vector ", zap.Any("vector", vector))

	// Tracer for OpenTelemetry
	ctx, span := tracing.CustomStartSpanCtx(ctx, "semanticUsecase.SearchVector", trace.SpanKindInternal)
	defer span.End()

	result, err := uc.semanticRepository.SearchVector(ctx, vector)
	if err != nil {
		logger.Error(ctx, "error searching vector", zap.Any("error", err))
		return nil, err
	}

	return result, nil	
}

// VerifyJWT handles the verification of a JWT, returning the claims if the token is valid.
func (uc *SemanticUseCase) IntentDecompose(ctx context.Context, intent entity.IntentDecompose) (*entity.IntentDecompose, error) {
	logger.Info(ctx, "semantic usecase IntentDecompose called with intent ", zap.Any("intent", intent))

	// Tracer for OpenTelemetry
	ctx, span := tracing.CustomStartSpanCtx(ctx, "semanticUsecase.IntentDecompose", trace.SpanKindInternal)
	defer span.End()

	// TODO: Implement the actual intent decomposition logic here
	return &intent, nil
}
