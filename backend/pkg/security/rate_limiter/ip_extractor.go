// Package ratelimiter provides IP address extraction utilities for rate limiting
// implementations. It defines interfaces and default implementations to reliably
// obtain client IP addresses from network requests, supporting trusted proxy
// detection and various address formats.
package ratelimiter

import (
	"net"
	"net/http"
	"strings"
)

// IPExtractor defines the interface for parsing client IP addresses from HTTP
// requests. Implementations should handle proxy headers like X-Forwarded-For
// and X-Real-IP when the connection originates from a trusted proxy.
type IPExtractor interface {
	// Extract processes the HTTP request to isolate the client IP address.
	// Returns the IP as a string suitable for rate limiting lookups.
	Extract(r *http.Request) string
}

// DefaultIPExtractor provides a standard implementation of IPExtractor with
// trusted proxy support. Extracts the client IP following this priority:
//  1. X-Real-IP header (only if RemoteAddr is a trusted proxy)
//  2. X-Forwarded-For header, first IP (only if RemoteAddr is a trusted proxy)
//  3. RemoteAddr (fallback when no trusted proxy headers are present)
//
// Without trusted proxies configured, the extractor always falls back to
// RemoteAddr, ignoring proxy headers to prevent IP spoofing.
type DefaultIPExtractor struct {
	trustedProxies []*net.IPNet
}

// NewDefaultIPExtractor creates a new DefaultIPExtractor with the specified
// trusted proxy CIDR ranges. Invalid CIDR notations are silently skipped.
// Parameters:
//
//	trustedCIDRs - CIDR notation ranges (e.g. "10.0.0.0/8", "192.168.0.0/16")
//
// Example:
//
//	NewDefaultIPExtractor("10.0.0.0/8", "172.16.0.0/12")
func NewDefaultIPExtractor(trustedCIDRs ...string) *DefaultIPExtractor {
	e := &DefaultIPExtractor{}
	for _, cidr := range trustedCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil {
			e.trustedProxies = append(e.trustedProxies, ipNet)
		}
	}
	return e
}

// Extract retrieves the real client IP from an HTTP request.
// When RemoteAddr belongs to a trusted proxy, it checks X-Real-IP and
// X-Forwarded-For headers. Otherwise, it returns RemoteAddr directly to
// prevent header spoofing by untrusted clients.
func (e *DefaultIPExtractor) Extract(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}

	if e.isTrusted(remoteIP) {
		if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
			if ip := cleanIP(xRealIP); ip != "" {
				return ip
			}
		}

		if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
			if firstIP := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0]); firstIP != "" {
				if ip := cleanIP(firstIP); ip != "" {
					return ip
				}
			}
		}
	}

	return remoteIP
}

// isTrusted checks whether the given IP address belongs to a configured trusted
// proxy range. Returns false if the IP cannot be parsed.
func (e *DefaultIPExtractor) isTrusted(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	for _, trusted := range e.trustedProxies {
		if trusted.Contains(parsedIP) {
			return true
		}
	}
	return false
}

// cleanIP attempts to separate an IP address from an optional port suffix.
// Handles IPv4 ("192.168.1.1:8080" → "192.168.1.1"), IPv6 with brackets
// ("[::1]:8080" → "::1"), and plain IPs without ports. Returns the original
// string if parsing fails.
func cleanIP(addr string) string {
	if addr == "" {
		return ""
	}
	if ip, _, err := net.SplitHostPort(addr); err == nil {
		return ip
	}
	if net.ParseIP(addr) != nil {
		return addr
	}
	return addr
}
