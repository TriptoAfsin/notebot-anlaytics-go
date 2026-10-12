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

// mcpAvailable reports whether the mcp schema can be there at all: it lives in
// the notebot v2 Postgres, never in the legacy MySQL.
func mcpAvailable() bool {
	switch strings.ToLower(os.Getenv("DB_DRIVER")) {
	case "postgres", "postgresql", "pg":
		return true
	}
	return false
}

// mcpRequestTotal is the all-time request counter across MCP keys, or 0 when
// the mcp schema isn't reachable.
func mcpRequestTotal(db *gorm.DB) int64 {
	if !mcpAvailable() {
		return 0
	}
	var total int64
	if err := db.Raw(`SELECT COALESCE(SUM(request_count), 0) FROM mcp.api_keys`).Scan(&total).Error; err != nil {
		return 0
	}
	return total
}

// pageParams reads page/limit/search with the same shape as /users/app.
func pageParams(c *fiber.Ctx) (page, limit, offset int, search string) {
	page = c.QueryInt("page", 1)
	limit = c.QueryInt("limit", 25)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 500 {
		limit = 25
	}
	return page, limit, (page - 1) * limit, strings.TrimSpace(c.Query("search", ""))
}

func mcpGuard(c *fiber.Ctx) error {
	if !utils.IsAdmin(c, config.GetAppConfig()) {
		return c.Status(401).JSON(fiber.Map{"error": config.AppMessages.User.UnauthorizedAccess})
	}
	if !mcpAvailable() {
		return c.Status(503).JSON(fiber.Map{"status": "MCP stats need DB_DRIVER=postgres"})
	}
	return nil
}

type mcpEmailRow struct {
	Email          string `json:"email" gorm:"column:email"`
	FirstSeen      string `json:"firstSeen" gorm:"column:first_seen"`
	LastSeen       string `json:"lastSeen" gorm:"column:last_seen"`
	Sources        string `json:"sources" gorm:"column:sources"`
	KeyRequests    int64  `json:"keyRequests" gorm:"column:key_requests"`
	RecoveryEmails int64  `json:"recoveryEmails" gorm:"column:recovery_emails"`
	Keys           int64  `json:"keys" gorm:"column:keys"`
	ActiveKeys     int64  `json:"activeKeys" gorm:"column:active_keys"`
	Requests       int64  `json:"requests" gorm:"column:requests"`
	LastUsed       string `json:"lastUsed" gorm:"column:last_used"`
	Unsubscribed   bool   `json:"unsubscribed" gorm:"column:unsubscribed"`
}

// GetMcpEmails lists every email registered for an MCP key, newest first, with
// when it signed up and what its keys have done. Admin only.
func GetMcpEmails(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetMcpEmails handler called")
	return func(c *fiber.Ctx) error {
		if err := mcpGuard(c); err != nil || c.Response().StatusCode() >= 400 {
			return err
		}
		page, limit, offset, search := pageParams(c)
		pattern := "%" + strings.ToLower(search) + "%"

		var total int64
		if err := db.Raw(`SELECT count(*) FROM mcp.contacts WHERE email LIKE ?`, pattern).Scan(&total).Error; err != nil {
			log.Printf("🔴 Error counting MCP emails: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		rows := []mcpEmailRow{}
		if err := db.Raw(`
			SELECT c.email,
			       to_char(c.first_seen AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI') AS first_seen,
			       to_char(c.last_seen AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI') AS last_seen,
			       array_to_string(c.sources, ', ') AS sources,
			       c.key_requests, c.recovery_emails,
			       COALESCE(k.keys, 0) AS keys, COALESCE(k.active_keys, 0) AS active_keys,
			       COALESCE(k.requests, 0) AS requests,
			       COALESCE(to_char(k.last_used AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI'), '') AS last_used,
			       c.unsubscribed_at IS NOT NULL AS unsubscribed
			FROM mcp.contacts c
			LEFT JOIN (
			  SELECT lower(email) AS email, count(*) AS keys,
			         count(*) FILTER (WHERE revoked_at IS NULL) AS active_keys,
			         SUM(request_count) AS requests, max(last_used_at) AS last_used
			  FROM mcp.api_keys WHERE tier <> 'admin' GROUP BY lower(email)
			) k ON k.email = c.email
			WHERE c.email LIKE ?
			ORDER BY c.first_seen DESC
			LIMIT ? OFFSET ?`, pattern, limit, offset).Scan(&rows).Error; err != nil {
			log.Printf("🔴 Error fetching MCP emails: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"emails":       rows,
			"total":        total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
			"search":       search,
		})
	}
}

type mcpKeyRow struct {
	ID        string `json:"id" gorm:"column:id"`
	Email     string `json:"email" gorm:"column:email"`
	Source    string `json:"source" gorm:"column:source"`
	Label     string `json:"label" gorm:"column:label"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at"`
	LastUsed  string `json:"lastUsed" gorm:"column:last_used"`
	Requests  int64  `json:"requests" gorm:"column:requests"`
	Revoked   bool   `json:"revoked" gorm:"column:revoked"`
}

// GetMcpKeys lists issued MCP keys, newest first. Only the key id is stored, so
// no usable secret ever leaves the database. Admin only.
func GetMcpKeys(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetMcpKeys handler called")
	return func(c *fiber.Ctx) error {
		if err := mcpGuard(c); err != nil || c.Response().StatusCode() >= 400 {
			return err
		}
		page, limit, offset, search := pageParams(c)
		pattern := "%" + strings.ToLower(search) + "%"

		var total int64
		if err := db.Raw(`SELECT count(*) FROM mcp.api_keys WHERE tier <> 'admin' AND lower(email) LIKE ?`, pattern).Scan(&total).Error; err != nil {
			log.Printf("🔴 Error counting MCP keys: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		rows := []mcpKeyRow{}
		if err := db.Raw(`
			SELECT id, lower(email) AS email, source, COALESCE(label, '') AS label,
			       to_char(created_at AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI') AS created_at,
			       COALESCE(to_char(last_used_at AT TIME ZONE 'Asia/Dhaka', 'YYYY-MM-DD HH24:MI'), '') AS last_used,
			       request_count AS requests, revoked_at IS NOT NULL AS revoked
			FROM mcp.api_keys
			WHERE tier <> 'admin' AND lower(email) LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?`, pattern, limit, offset).Scan(&rows).Error; err != nil {
			log.Printf("🔴 Error fetching MCP keys: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"keys":         rows,
			"total":        total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
			"search":       search,
		})
	}
}
