package utils

import (
	"fmt"
	"strings"
)

type URLResolver struct {
	XForwardedProto string `header:"x-forwarded-proto" doc:"Protocol from CloudFront/ALB (http or https)"`
	XForwardedHost  string `header:"x-forwarded-host" doc:"Host from CloudFront/ALB"`
	Host            string `header:"host" doc:"Fallback host header"`
}

func (r *URLResolver) ResolveBaseURL(configuredBaseURL string) string {
	// 1. Explicit override from environment config (if provided)
	if configuredBaseURL != "" {
		return strings.TrimRight(configuredBaseURL, "/")
	}

	// 2. Resolve protocol (defaults to https in serverless/CloudFront environments)
	proto := r.XForwardedProto
	if proto == "" {
		proto = "https"
	}

	// 3. Resolve host (prefers X-Forwarded-Host from CloudFront, falls back to Host)
	host := r.XForwardedHost
	if host == "" {
		host = r.Host
	}
	if host == "" {
		host = "localhost:8080"
	}

	return fmt.Sprintf("%s://%s", proto, host)
}
