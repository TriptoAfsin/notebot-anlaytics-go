package utils

import (
	"crypto/subtle"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/gofiber/fiber/v2"
)

// ValidateAdminKey checks if the provided admin key matches the configured key
func ValidateAdminKey(c *fiber.Ctx, appConfig config.AppConfig) error {
	if adminKey := c.Query("adminKey"); adminKey != appConfig.ADMIN_AUTH_KEY {
		return c.Status(401).JSON(fiber.Map{
			"Error": "🔴 Unauthorized Access !",
		})
	}
	return nil
}

// ValidateEmail checks if the provided email matches a valid email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// Failed admin attempts per client IP. The dashboard password is short, so
// guessing is capped: after maxAdminFailures wrong keys in adminFailWindow an
// IP is refused until the window passes, even with the right key.
const (
	maxAdminFailures = 10
	adminFailWindow  = 15 * time.Minute
)

type adminFailures struct {
	count int
	since time.Time
}

var (
	adminFailMu sync.Mutex
	adminFails  = map[string]*adminFailures{}
)

// clientIP prefers the first X-Forwarded-For hop: on Railway c.IP() is the proxy.
func clientIP(c *fiber.Ctx) string {
	if fwd := c.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	return c.IP()
}

func secretMatches(key, secret string) bool {
	return secret != "" && subtle.ConstantTimeCompare([]byte(key), []byte(secret)) == 1
}

// IsAdmin reports whether the request carries the admin key or the dashboard
// password, either in the X-Admin-Key header (preferred: stays out of URLs and
// logs) or the legacy adminKey query parameter. Constant-time compare.
func IsAdmin(c *fiber.Ctx, appConfig config.AppConfig) bool {
	key := c.Get("X-Admin-Key")
	if key == "" {
		key = c.Query("adminKey")
	}
	if key == "" {
		return false
	}

	ip := clientIP(c)
	adminFailMu.Lock()
	defer adminFailMu.Unlock()
	f := adminFails[ip]
	if f != nil && time.Since(f.since) > adminFailWindow {
		delete(adminFails, ip)
		f = nil
	}
	if f != nil && f.count >= maxAdminFailures {
		return false
	}

	if secretMatches(key, appConfig.ADMIN_AUTH_KEY) || secretMatches(key, appConfig.DASHBOARD_PASS) {
		return true
	}
	if f == nil {
		if len(adminFails) > 10_000 { // bound memory under a spray of IPs
			adminFails = map[string]*adminFailures{}
		}
		f = &adminFailures{since: time.Now()}
		adminFails[ip] = f
	}
	f.count++
	return false
}
