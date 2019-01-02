package general

import "context"

type BaseService struct {
}

func (s *BaseService) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, nil
}

func (s *BaseService) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	return nil, nil
}
