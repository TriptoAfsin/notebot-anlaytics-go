package handler

import (
	"log"
	"time"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApiHandler(c *fiber.Ctx) error {
	db := db.DB
	dbStatus := true

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = false
	}

	log.Println("🟢 GET: ApiHandler handler called")

	apiStatus := fiber.Map{
		"endPoints": []string{
			"/",
			"/users",
			"/missed",
			"/notes",
			"/games/notebird",
			"/games/notedino",
		},
		"db_connection": dbStatus,
		"mode":          config.GetAppConfig().ENVIRONMENT,
	}

	return c.JSON(apiStatus)
}

func HealthCheckHandler(c *fiber.Ctx) error {
	log.Println("💉 GET: HealthCheckHandler handler called")
	db := db.DB
	dbStatus := true

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = false
	}
	return c.JSON(fiber.Map{
		"status": "OK",
		"db":     dbStatus,
	})
}

// GetDailyReport handles fetching daily reports with optional platform filter
func GetDailyReport(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetDailyReport handler called")
	return func(c *fiber.Ctx) error {
		// Get query parameters with defaults
		platform := c.Query("platform")
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 500)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 500
		}

		offset := (page - 1) * limit

		// Get date filter parameters
		startDate := c.Query("startDate") // Format: YYYY-MM-DD
		endDate := c.Query("endDate")     // Format: YYYY-MM-DD

		var reports []map[string]interface{}
		var totalCount int64
		var query, countQuery string
		var queryParams []interface{}
		var err error

		// Build base queries
		baseQuery := "SELECT * FROM bot_daily_report WHERE 1=1"
		baseCountQuery := "SELECT COUNT(*) FROM bot_daily_report WHERE 1=1"

		// Add platform filter if specified
		if platform != "" {
			baseQuery += " AND platform = ?"
			baseCountQuery += " AND platform = ?"
			queryParams = append(queryParams, platform)
		}

		// Add date filters if specified
		if startDate != "" {
			baseQuery += " AND date >= ?"
			baseCountQuery += " AND date >= ?"
			queryParams = append(queryParams, startDate)
		}
		if endDate != "" {
			baseQuery += " AND date <= ?"
			baseCountQuery += " AND date <= ?"
			queryParams = append(queryParams, endDate)
		}

		// Add ordering and pagination
		query = baseQuery + " ORDER BY date DESC LIMIT ? OFFSET ?"
		countQuery = baseCountQuery

		// Add pagination params
		queryParams = append(queryParams, limit, offset)

		// Get total count
		err = db.Raw(countQuery, queryParams[:len(queryParams)-2]...).Scan(&totalCount).Error
		if err != nil {
			log.Printf("🔴 Error while counting records: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Get paginated data
		err = db.Raw(query, queryParams...).Scan(&reports).Error
		if err != nil {
			log.Printf("🔴 Error while fetching daily reports: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Calculate total pages
		totalPages := (totalCount + int64(limit) - 1) / int64(limit)

		return c.Status(200).JSON(fiber.Map{
			"status": config.AppMessages.API.OperationSuccessful,
			"data":   reports,
			"meta": fiber.Map{
				"current_page": page,
				"per_page":     limit,
				"total_items":  totalCount,
				"total_pages":  totalPages,
			},
		})
	}
}

func GetDailyReportSummary(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetDailyReportSummary handler called")
	return func(c *fiber.Ctx) error {
		// Get total counts for each platform
		//Each query uses COALESCE to handle null values
		var totalAppCount, totalBotCount int64
		if err := db.Raw("SELECT COALESCE(SUM(count), 0) FROM bot_daily_report WHERE platform = 'app'").Scan(&totalAppCount).Error; err != nil {
			log.Printf("🔴 Error fetching app count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}
		if err := db.Raw("SELECT COALESCE(SUM(count), 0) FROM bot_daily_report WHERE platform = 'bot'").Scan(&totalBotCount).Error; err != nil {
			log.Printf("🔴 Error fetching bot count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Calculate percentages
		total := float64(totalAppCount + totalBotCount)
		var appPercentage, botPercentage float64
		if total > 0 {
			appPercentage = float64(totalAppCount) / total * 100
			botPercentage = float64(totalBotCount) / total * 100
		}

		// Get highest counts and dates
		type MaxCount struct {
			Count int64
			Date  string
		}
		var highestApp, highestBot MaxCount
		var lowestDate, highestDate string
		var lowestApiCount int64

		// Get highest app count
		if err := db.Raw(`
			SELECT count, date 
			FROM bot_daily_report 
			WHERE platform = 'app' 
			ORDER BY count DESC 
			LIMIT 1
		`).Scan(&highestApp).Error; err != nil {
			log.Printf("🔴 Error fetching highest app count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Get highest bot count
		if err := db.Raw(`
			SELECT count, date 
			FROM bot_daily_report 
			WHERE platform = 'bot' 
			ORDER BY count DESC 
			LIMIT 1
		`).Scan(&highestBot).Error; err != nil {
			log.Printf("🔴 Error fetching highest bot count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Get dates with highest and lowest total API counts
		if err := db.Raw(`
			SELECT date 
			FROM (
				SELECT date, SUM(count) as total_count 
				FROM bot_daily_report 
				GROUP BY date
			) t 
			ORDER BY total_count DESC 
			LIMIT 1
		`).Scan(&highestDate).Error; err != nil {
			log.Printf("🔴 Error fetching highest date: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		if err := db.Raw(`
			SELECT date 
			FROM (
				SELECT date, SUM(count) as total_count 
				FROM bot_daily_report 
				GROUP BY date
			) t 
			ORDER BY total_count ASC 
			LIMIT 1
		`).Scan(&lowestDate).Error; err != nil {
			log.Printf("🔴 Error fetching lowest date: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Get lowest total API count
		if err := db.Raw(`
			SELECT total_count 
			FROM (
				SELECT date, SUM(count) as total_count 
				FROM bot_daily_report 
				GROUP BY date
			) t 
			ORDER BY total_count ASC 
			LIMIT 1
		`).Scan(&lowestApiCount).Error; err != nil {
			log.Printf("🔴 Error fetching lowest API count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		// Determine highest platform
		highestPlatform := "app"
		if totalBotCount > totalAppCount {
			highestPlatform = "bot"
		}

		return c.Status(200).JSON(fiber.Map{
			"status": config.AppMessages.API.OperationSuccessful,
			"kpi": fiber.Map{
				"totalAppPlatformCount":   totalAppCount,
				"totalBotPlatformCount":   totalBotCount,
				"appPlatformPercentage":   float64(int(appPercentage*100)) / 100, // Round to 2 decimal places
				"botPlatformPercentage":   float64(int(botPercentage*100)) / 100, // Round to 2 decimal places
				"highestPlatform":         highestPlatform,
				"lowestPlatform":          map[string]string{"app": "bot", "bot": "app"}[highestPlatform],
				"highestAppPlatformCount": highestApp.Count,
				"highestBotPlatformCount": highestBot.Count,
				"highestApiCountDate":     highestDate,
				"lowestApiCountDate":      lowestDate,
				"lowestApiCount":          lowestApiCount,
			},
		})
	}
}

// PostDailyReport handles creating or updating daily report entries
func PostDailyReport(db *gorm.DB) fiber.Handler {
	log.Println("🔵 POST: PostDailyReport handler called")
	return func(c *fiber.Ctx) error {
		// Admin auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": config.AppMessages.API.UnauthorizedAccess,
			})
		}

		// Parse request body
		report := struct {
			Platform string `json:"platform"`
		}{}

		if err := c.BodyParser(&report); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": config.AppMessages.API.BadRequest})
		}

		// Validate platform
		validPlatforms := []string{"bot", "app"}
		isValidPlatform := false
		for _, p := range validPlatforms {
			if p == report.Platform {
				isValidPlatform = true
				break
			}
		}

		if !isValidPlatform {
			return c.Status(400).JSON(fiber.Map{"status": config.AppMessages.API.InvalidPlatform})
		}

		currentDate := time.Now().Format("2006-01-02")

		// Check if entry exists
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM bot_daily_report WHERE date = ? AND platform = ?",
			currentDate, report.Platform).Scan(&count).Error

		if err != nil {
			log.Printf("🔴 Error while checking existing log: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.LogCheckError})
		}

		if count > 0 {
			// Update existing entry
			err = db.Exec("UPDATE bot_daily_report SET count = count + 1 WHERE date = ? AND platform = ?",
				currentDate, report.Platform).Error

			if err != nil {
				log.Printf("🔴 Error while updating daily api count: %v", err)
				return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.UpdateCountError})
			}

			return c.Status(200).JSON(fiber.Map{
				"status": config.AppMessages.API.IncrementSuccess,
			})
		}

		// Create new entry
		err = db.Exec("INSERT INTO bot_daily_report (date, count, platform) VALUES (?, 1, ?)",
			currentDate, report.Platform).Error

		if err != nil {
			log.Printf("🔴 Error while inserting new log: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.API.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": config.AppMessages.API.NewLogSuccess,
		})
	}
}
