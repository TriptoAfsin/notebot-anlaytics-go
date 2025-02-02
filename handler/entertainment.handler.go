package handler

import (
	"log"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/lib/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GameScore struct {
	Date     string `json:"date"`
	Score    int    `json:"score"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

// PostNoteBirdScore handles posting scores for NoteBird game
func PostNoteBirdScore(db *gorm.DB) fiber.Handler {
	log.Println("🔵 POST: PostNoteBirdScore handler called")
	return func(c *fiber.Ctx) error {
		// Auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": config.AppMessages.Game.UnauthorizedAccess,
			})
		}

		// Parse body
		var score GameScore
		if err := c.BodyParser(&score); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Game.BadRequest,
			})
		}

		// Validate required fields and score value
		if score.Email == "" || score.Date == "" || score.Score <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Game.InvalidFields,
			})
		}

		// Validate email
		if !utils.ValidateEmail(score.Email) {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Error.InvalidEmail,
			})
		}

		// Insert score using raw SQL
		query := `INSERT INTO game_hof (date, score, email, user_name) VALUES (?, ?, ?, ?)`
		if err := db.Exec(query, score.Date, score.Score, score.Email, score.UserName).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.OperationUnsuccessful,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"gameScoreInfo": score,
			"status":        config.AppMessages.Game.ScoreInsertSuccess,
		})
	}
}

// GetNoteBirdHof handles getting top scores for NoteBird game
func GetNoteBirdHof(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetNoteBirdHof handler called")
	return func(c *fiber.Ctx) error {
		// Get pagination parameters from query
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 100)
		search := c.Query("search", "")
		offset := (page - 1) * limit

		// Build the WHERE clause for search
		whereClause := "1=1"
		params := []interface{}{}

		if search != "" {
			whereClause = "(email LIKE ? OR user_name LIKE ?)"
			searchPattern := "%" + search + "%"
			params = append(params, searchPattern, searchPattern)
		}

		// Get total count with search filter
		var total int64
		countQuery := "SELECT COUNT(*) FROM game_hof WHERE " + whereClause
		if err := db.Raw(countQuery, params...).Scan(&total).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		// Get paginated and filtered scores
		var results []GameScore
		query := `
			SELECT date, score, email, user_name 
			FROM game_hof 
			WHERE ` + whereClause + `
			ORDER BY score DESC 
			LIMIT ? OFFSET ?`

		// Add pagination parameters
		params = append(params, limit, offset)

		if err := db.Raw(query, params...).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof":          results,
			"total":        total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
			"search":       search,
		})
	}
}

// PostNoteDinoScore handles posting scores for NoteDino game
func PostNoteDinoScore(db *gorm.DB) fiber.Handler {
	log.Println("🔵 POST: PostNoteDinoScore handler called")
	return func(c *fiber.Ctx) error {
		// Auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": config.AppMessages.Game.UnauthorizedAccess,
			})
		}

		// Parse body
		var score GameScore
		if err := c.BodyParser(&score); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Game.BadRequest,
			})
		}

		// Validate required fields and score value
		if score.Email == "" || score.Date == "" || score.Score <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Game.InvalidFields,
			})
		}

		// Validate email
		if !utils.ValidateEmail(score.Email) {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.Error.InvalidEmail,
			})
		}

		// Insert score using raw SQL
		query := `INSERT INTO game_hof_noteDino (date, score, email, user_name) VALUES (?, ?, ?, ?)`
		if err := db.Exec(query, score.Date, score.Score, score.Email, score.UserName).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.OperationUnsuccessful,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"gameScoreInfo": score,
			"status":        config.AppMessages.Game.ScoreInsertSuccess,
		})
	}
}

// GetNoteDinoHof handles getting top scores for NoteDino game
func GetNoteDinoHof(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetNoteDinoHof handler called")
	return func(c *fiber.Ctx) error {
		// Get pagination parameters from query
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 100)
		search := c.Query("search", "")
		offset := (page - 1) * limit

		// Build the WHERE clause for search
		whereClause := "1=1"
		params := []interface{}{}

		if search != "" {
			whereClause = "(email LIKE ? OR user_name LIKE ?)"
			searchPattern := "%" + search + "%"
			params = append(params, searchPattern, searchPattern)
		}

		// Get total count with search filter
		var total int64
		countQuery := "SELECT COUNT(*) FROM game_hof_noteDino WHERE " + whereClause
		if err := db.Raw(countQuery, params...).Scan(&total).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		// Get paginated and filtered scores
		var results []GameScore
		query := `
			SELECT date, score, email, user_name 
			FROM game_hof_noteDino 
			WHERE ` + whereClause + `
			ORDER BY score DESC 
			LIMIT ? OFFSET ?`

		// Add pagination parameters
		params = append(params, limit, offset)

		if err := db.Raw(query, params...).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof":          results,
			"total":        total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
			"search":       search,
		})
	}
}
