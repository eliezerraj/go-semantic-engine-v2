package module

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	"go.uber.org/zap"
	"github.com/eliezerraj/go-core/v3/auth"

	"go.opentelemetry.io/otel/trace"

	"github.com/eliezerraj/go-core/v3/httpclient"
	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-semantic-engine-v2/application/config"
	"github.com/go-semantic-engine-v2/application/domain/entity"
	"github.com/go-semantic-engine-v2/application/tracing"
)

const (
	AcceptHeader      	= "Accept"
	ContentTypeHeader 	= "Content-Type"
	ConnectionHeader  	= "Connection"
	KeepAlive         	= "keep-alive"
	XRequestIDHeader	= "X-Request-ID"
	Authorization 		= "Authorization"
	RequestIDHeaderName = "x-request-id"
)

type ModelModule struct {
	cfg *config.Config
	client	httpclient.IHTTPClient
	authClientService *auth.AuthClientService
}

func NewModelModule(cfg *config.Config, 
						client httpclient.IHTTPClient,
						authClientService *auth.AuthClientService) ModelModule {
	logger.InfoOutCtx("NewModelModule called with authClientService", zap.Any("authClientService", authClientService))

	return ModelModule{
		cfg: cfg,
		client: client,
		authClientService: authClientService,
	}
}

func (im *ModelModule) ModelEmbeddingPost(ctx context.Context, vector entity.VectorSearch) (*entity.TeiModel, error) {
	logger.Info(ctx, "model module ModelEmbeddingPost called")

	ctxHttpTimeout, cancel := context.WithTimeout(ctx, im.cfg.Vector.Timeout)
	defer cancel()

	ctx, span := tracing.CustomStartSpanCtx(ctxHttpTimeout, "modelModule.ModelEmbeddingPost", trace.SpanKindInternal)
	defer span.End()

	endpoint := fmt.Sprintf("%s%s", im.cfg.Vector.Endpoint, im.cfg.Vector.UrlPath)

	logger.Info(ctx, "model module ModelEmbeddingPost request", zap.String("method", http.MethodPost), zap.String("endpoint", endpoint))

	payload := &entity.TeiModel{
		Inputs: vector.Text,
		Normalize: true,
		Truncate: false,
		TruncationDirection: "Right",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "Failed to marshal payload", zap.Error(err))
		return nil, err
	}

	logger.Debug(ctx, "model module ModelEmbeddingPost payload", zap.Any("payload", payload))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		logger.Error(ctx, "Failed to create request", zap.Error(err))
		return nil, err
	}

	// Set headers for the request. the const are in model_module.go file
	xrequestid, ok := ctx.Value(RequestIDHeaderName).(string)
	if !ok {
		xrequestid = "not-informed"
	}

	headers := map[string]string{
		ConnectionHeader:  KeepAlive,
		AcceptHeader:      "application/json",
		ContentTypeHeader: "application/json",
		XRequestIDHeader: xrequestid,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := im.client.Do(req.WithContext(ctx))
	if err != nil {
		logger.Error(ctx, "Failed to perform request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(ctx, "Model service returned unexpected status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("model service returned status: %d", resp.StatusCode)
	}

	bytePayload, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, "Failed to read response body", zap.Error(err))
		return nil, err
	}

	var embeddings [][]float64
	err = json.Unmarshal(bytePayload, &embeddings)
	if err != nil {
		logger.Error(ctx, "Failed to unmarshal embeddings", zap.Error(err))
		return nil, err
	}

	if len(embeddings) > 0 {
		payload.Embedding = embeddings[0]
	}

	return payload, nil
}