package dao

import (
	"context"
	"time"

	analysispb "proto/analysis"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalysisResult struct {
	Intent         string
	RawCategory    string
	Category       string
	Confidence     float32
	TranslatedText string
}

type AnalysisAccessor struct {
	conn   *grpc.ClientConn
	client analysispb.AnalysisServiceClient
}

func NewAnalysisAccessor(addr string) (*AnalysisAccessor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &AnalysisAccessor{
		conn:   conn,
		client: analysispb.NewAnalysisServiceClient(conn),
	}, nil
}

func (a *AnalysisAccessor) Close() error {
	if a.conn == nil {
		return nil
	}

	return a.conn.Close()
}

func (a *AnalysisAccessor) ClassifyMessage(
	ctx context.Context,
	message string,
) (*AnalysisResult, error) {
	resp, err := a.client.ClassifyMessage(ctx, &analysispb.ClassifyMessageRequest{
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	return &AnalysisResult{
		Intent:         resp.Intent,
		RawCategory:    resp.RawCategory,
		Category:       resp.Category,
		Confidence:     resp.Confidence,
		TranslatedText: resp.TranslatedText,
	}, nil
}
