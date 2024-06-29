package fbgraph

import "context"

type ContextVar string

const (
	ContextGraphAPIVersion ContextVar = "graph_api_version"
)

type APIVersion string

const (
	APIVersion15 APIVersion = "v15.0"
	APIVersion16 APIVersion = "v16.0"
	APIVersion17 APIVersion = "v17.0"
	APIVersion20 APIVersion = "v20.0"
)

func WithAPIVersion(ctx context.Context, version APIVersion) context.Context {
	return context.WithValue(ctx, ContextGraphAPIVersion, version)
}

func GetAPIVersion(ctx context.Context, defaultv APIVersion) APIVersion {
	if ctx == nil || ctx.Value(ContextGraphAPIVersion) == nil {
		if defaultv != "" {
			return defaultv
		}

		return APIVersion15
	}
	return ctx.Value(ContextGraphAPIVersion).(APIVersion)
}
