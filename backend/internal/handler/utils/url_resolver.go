package utils

import (
	"net/url"
	"strings"
)

type URLResolver struct {
	XForwardedProto          string `header:"x-forwarded-proto"`
	CloudFrontForwardedProto string `header:"cloudfront-forwarded-proto"`
	XForwardedHost           string `header:"x-forwarded-host"`
	Host                     string `header:"host"`
	Origin                   string `header:"origin"`
	Referer                  string `header:"referer"`
}

func (r *URLResolver) ResolveBaseURL(configuredBaseURL string) string {
	if configuredBaseURL != "" {
		return strings.TrimRight(configuredBaseURL, "/")
	}

	// 1. Extract Host (Fall back across headers, Origin, and Referer)
	host := r.XForwardedHost
	if host == "" {
		host = r.Host
	}
	if host == "" && r.Origin != "" {
		if u, err := url.Parse(r.Origin); err == nil && u.Host != "" {
			host = u.Host
		}
	}
	if host == "" && r.Referer != "" {
		if u, err := url.Parse(r.Referer); err == nil && u.Host != "" {
			host = u.Host
		}
	}
	if host == "" {
		host = "localhost:8080"
	}

	// 2. Extract Protocol / Scheme
	scheme := r.XForwardedProto
	if scheme == "" {
		scheme = r.CloudFrontForwardedProto
	}
	if scheme == "" && r.Origin != "" {
		if u, err := url.Parse(r.Origin); err == nil && u.Scheme != "" {
			scheme = u.Scheme
		}
	}
	if scheme == "" && r.Referer != "" {
		if u, err := url.Parse(r.Referer); err == nil && u.Scheme != "" {
			scheme = u.Scheme
		}
	}
	if scheme == "" {
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}

	return scheme + "://" + host
}
