package handler

import (
	"log"
	"os"
	"strings"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/lib/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// The NoteBot MCP server (ai.butexnotebot.com) keeps its key bookkeeping in the
// `mcp` schema of the same Postgres: mcp.api_keys (one row per issued key, with a
// running request_count) and mcp.contacts (one row per email that asked for a key).
// Daily request volume comes from bot_daily_report rows with platform "ai".

type mcpDay struct {
	Date     string `json:"date" gorm:"column:date"`
	Emails   int64  `json:"emails" gorm:"column:emails"`
	Keys     int64  `json:"keys" gorm:"column:keys"`
	Requests int64  `json:"requests" gorm:"column:requests"`
}

type mcpSource struct {
	Source string `json:"source" gorm:"column:source"`
	Keys   int64  `json:"keys" gorm:"column:keys"`
}

type mcpTopKey struct {
	Email    string `json:"email" gorm:"column:email"`
	Source   string `json:"source" gorm:"column:source"`
	Requests int64  `json:"requests" gorm:"column:requests"`
	LastUsed string `json:"lastUsed" gorm:"column:last_used"`
}

// maskEmail keeps enough of an address to spot repeats: "afshin@gmail.com" -> "af***@gmail.com".
func maskEmail(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return "***"
	}
	local := email[:at]
	if len(local) > 2 {
		local = local[:2]
	}
	return local + "***" + email[at:]
}

// GetMcpSummary returns how many emails and keys the MCP server has registered
// and how many requests it has served. Admin only: it describes who signed up.
func GetMcpSummary(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetMcpSummary handler called")
	return func(c *fiber.Ctx) error {
		if !utils.IsAdmin(c, config.GetAppConfig()) {
			return c.Status(401).JSON(fiber.Map{"error": config.AppMessages.User.UnauthorizedAccess})
		}
		if driver := strings.ToLower(os.Getenv("DB_DRIVER")); driver != "postgres" && driver != "postgresql" && driver != "pg" {
			return c.Status(503).JSON(fiber.Map{"status": "MCP stats need DB_DRIVER=postgres"})
		}

		days := c.QueryInt("days", 30)
		if days < 7 {
			days = 7
		}
		if days > 365 {
			days = 365
		}

		fail := func(what string, err error) error {
			log.Printf("🔴 Error fetching MCP %s: %v", what, err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		var emails struct {
			Total        int64 `json:"total" gorm:"column:total"`
			New7d        int64 `json:"new7d" gorm:"column:new7d"`
			New30d       int64 `json:"new30d" gorm:"column:new30d"`
			Unsubscribed int64 `json:"unsubscribed" gorm:"column:unsubscribed"`
		}
		if err := db.Raw(`
			SELECT count(*) AS total,
			       count(*) FILTER (WHERE first_seen > now() - interval '7 days') AS new7d,
			       count(*) FILTER (WHERE first_seen > now() - interval '30 days') AS new30d,
			       count(*) FILTER (WHERE unsubscribed_at IS NOT NULL) AS unsubscribed
			FROM mcp.contacts`).Scan(&emails).Error; err != nil {
			return fail("emails", err)
		}

		// The owner's own admin-tier tokens are not registrations; leave them out.
		var keys struct {
			Total     int64 `json:"total" gorm:"column:total"`
			Active    int64 `json:"active" gorm:"column:active"`
			Revoked   int64 `json:"revoked" gorm:"column:revoked"`
			Used      int64 `json:"used" gorm:"column:used"`
			Active7d  int64 `json:"active7d" gorm:"column:active7d"`
			Active24h int64 `json:"active24h" gorm:"column:active24h"`
			New7d     int64 `json:"new7d" gorm:"column:new7d"`
			New30d    int64 `json:"new30d" gorm:"column:new30d"`
		}
		if err := db.Raw(`
			SELECT count(*) AS total,
			       count(*) FILTER (WHERE revoked_at IS NULL) AS active,
			       count(*) FILTER (WHERE revoked_at IS NOT NULL) AS revoked,
			       count(*) FILTER (WHERE request_count > 0) AS used,
			       count(*) FILTER (WHERE last_used_at > now() - interval '7 days') AS active7d,
			       count(*) FILTER (WHERE last_used_at > now() - interval '24 hours') AS active24h,
			       count(*) FILTER (WHERE created_at > now() - interval '7 days') AS new7d,
			       count(*) FILTER (WHERE created_at > now() - interval '30 days') AS new30d
			FROM mcp.api_keys WHERE tier <> 'admin'`).Scan(&keys).Error; err != nil {
			return fail("keys", err)
		}

		var requests struct {
			Total    int64 `json:"total" gorm:"column:total"`
			Today    int64 `json:"today" gorm:"column:today"`
			Last7d   int64 `json:"last7d" gorm:"column:last7d"`
			Last30d  int64 `json:"last30d" gorm:"column:last30d"`
			Reported int64 `json:"reported" gorm:"column:reported"`
		}
		if err := db.Raw(`
			SELECT COALESCE(SUM(count), 0) AS reported,
			       COALESCE(SUM(count) FILTER (WHERE date = (now() AT TIME ZONE 'Asia/Dhaka')::date), 0) AS today,
			       COALESCE(SUM(count) FILTER (WHERE date > (now() AT TIME ZONE 'Asia/Dhaka')::date - 7), 0) AS last7d,
			       COALESCE(SUM(count) FILTER (WHERE date > (now() AT TIME ZONE 'Asia/Dhaka')::date - 30), 0) AS last30d
			FROM bot_daily_report WHERE platform = 'ai'`).Scan(&requests).Error; err != nil {
			return fail("daily requests", err)
		}
		// Scan above resets the struct, so the all-time key counter goes in after it.
		if err := db.Raw(`SELECT COALESCE(SUM(request_count), 0) FROM mcp.api_keys`).Scan(&requests.Total).Error; err != nil {
			return fail("request total", err)
		}

		daily := []mcpDay{}
		if err := db.Raw(`
			WITH d AS (
			  SELECT generate_series((now() AT TIME ZONE 'Asia/Dhaka')::date - (? - 1), (now() AT TIME ZONE 'Asia/Dhaka')::date, interval '1 day')::date AS day
			)
			SELECT to_char(d.day, 'YYYY-MM-DD') AS date,
			       (SELECT count(*) FROM mcp.contacts WHERE (first_seen AT TIME ZONE 'Asia/Dhaka')::date = d.day) AS emails,
			       (SELECT count(*) FROM mcp.api_keys WHERE tier <> 'admin' AND (created_at AT TIME ZONE 'Asia/Dhaka')::date = d.day) AS keys,
			       (SELECT COALESCE(SUM(count), 0) FROM bot_daily_report WHERE platform = 'ai' AND date = d.day) AS requests
			FROM d ORDER BY d.day`, days).Scan(&daily).Error; err != nil {
			return fail("daily series", err)
		}

		sources := []mcpSource{}
		if err := db.Raw(`
			SELECT source, count(*) AS keys FROM mcp.api_keys
			WHERE tier <> 'admin' GROUP BY source ORDER BY keys DESC`).Scan(&sources).Error; err != nil {
			return fail("sources", err)
		}

		top := []mcpTopKey{}
		if err := db.Raw(`
			SELECT email, source, request_count AS requests,
			       COALESCE(to_char(last_used_at AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI'), '') AS last_used
			FROM mcp.api_keys
			WHERE tier <> 'admin' AND request_count > 0
			ORDER BY request_count DESC LIMIT 10`).Scan(&top).Error; err != nil {
			return fail("top keys", err)
		}
		for i := range top {
			top[i].Email = maskEmail(top[i].Email)
		}

		return c.Status(200).JSON(fiber.Map{
			"status":   config.AppMessages.API.OperationSuccessful,
			"emails":   emails,
			"keys":     keys,
			"requests": requests,
			"daily":    daily,
			"sources":  sources,
			"topKeys":  top,
		})
	}
}
