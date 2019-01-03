package general

import "context"

type BaseService struct {
}

func (s *BaseService) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{}, nil
}

func (s *BaseService) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{Version: Version}, nil
}
