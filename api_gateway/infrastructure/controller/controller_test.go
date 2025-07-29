package controller

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/metrics"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/dns"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/port"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/prefix"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"testing"
)

func createMockReverseProxy() *httputil.ReverseProxy {
	target, _ := url.Parse(prefix.HttpPrefix + dns.Localhost + ":" + strconv.Itoa(port.Http))
	return httputil.NewSingleHostReverseProxy(target)
}

func TestRerouteHandler(t *testing.T) {
	proxy := createMockReverseProxy()

	// create the handler
	metricsInstance := metrics.New()
	ctrl := NewController(metricsInstance)
	handler := ctrl.RerouteHandler(endpoint.Service, proxy)

	tests := []struct {
		name         string
		requestPath  string
		expectedPath string
		description  string
	}{
		{
			name:         "remove service prefix",
			requestPath:  "/service/api/v1/users",
			expectedPath: "/api/v1/users",
			description:  "should remove /service prefix",
		},
		{
			name:         "exact service path",
			requestPath:  "/service",
			expectedPath: "/",
			description:  "should convert to root when only service path",
		},
		{
			name:         "service with trailing slash",
			requestPath:  "/service/",
			expectedPath: "/",
			description:  "should handle trailing slash correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			// we can't easily test the actual proxy call, but we can test URL modification
			originalPath := req.URL.Path
			removePrefix(req, endpoint.Service)

			if req.URL.Path != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, req.URL.Path)
			}

			// restore original path for the handler test
			req.URL.Path = originalPath

			// test the handler doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RerouteHandler panicked: %v", r)
				}
			}()

			handler(w, req)
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	tests := []struct {
		name         string
		originalPath string
		service      string
		expectedPath string
	}{
		{
			name:         "remove simple prefix",
			originalPath: "/service/users",
			service:      "/service",
			expectedPath: "/users",
		},
		{
			name:         "remove complex prefix",
			originalPath: "/api/v1/service/data",
			service:      "/api/v1/service",
			expectedPath: "/data",
		},
		{
			name:         "exact match becomes root",
			originalPath: "/service",
			service:      "/service",
			expectedPath: "/",
		},
		{
			name:         "no prefix match",
			originalPath: "/other/service",
			service:      "/service",
			expectedPath: "/other/service",
		},
		{
			name:         "empty path after prefix removal",
			originalPath: "/service",
			service:      "/service",
			expectedPath: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.originalPath, nil)

			removePrefix(req, tt.service)

			if req.URL.Path != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, req.URL.Path)
			}
		})
	}
}
