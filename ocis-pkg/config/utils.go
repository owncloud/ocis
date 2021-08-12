package config

func TracingService(enabled bool) string {
	if enabled {
		return "opentelemetry"
	}
	return ""
}
