package services

import "testing"

func TestHealthAPI(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "API Health Check", want: "API is healthy."},
	}

	healthAPI := HealthAPI{BaseURL: "http://localhost:8090"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiHealth := healthAPI.CheckAPIHealth().Message
			if apiHealth != tt.want {
				t.Fatalf("API Health check failed - got: %s, want: %s", apiHealth, tt.want)
			}
			t.Log("API Health check passed")
		})
	}
}
