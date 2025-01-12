package handler

import (
	"regexp"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GameScore struct {
	Date     string `json:"date"`
	Score    int    `json:"score"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

// validateEmail checks if the email is valid
// validateEmail uses regex to validate email format
func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// PostNoteBirdScore handles posting scores for NoteBird game
func PostNoteBirdScore(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": "🔴 Unauthorized Access !",
			})
		}

		// Parse body
		var score GameScore
		if err := c.BodyParser(&score); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request",
			})
		}

		// Validate required fields
		if score.Email == "" || score.Score == 0 || score.Date == "" {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request",
			})
		}

		// Validate email
		if !validateEmail(score.Email) {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request, Invalid Email",
			})
		}

		// Insert score using raw SQL
		query := `INSERT INTO game_hof (date, score, email, user_name) VALUES (?, ?, ?, ?)`
		if err := db.Exec(query, score.Date, score.Score, score.Email, score.UserName).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "🔴 Operation was unsuccessful!",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"gameScoreInfo": score,
			"status":        "🟢 Game score insertion was successful",
		})
	}
}

// GetNoteBirdHof handles getting top scores for NoteBird game
func GetNoteBirdHof(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var results []GameScore
		query := `SELECT date, score, email, user_name FROM game_hof ORDER BY score DESC LIMIT 10`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "🔴 Error while fetching hof",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof": results,
		})
	}
}

// PostNoteDinoScore handles posting scores for NoteDino game
func PostNoteDinoScore(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": "🔴 Unauthorized Access !",
			})
		}

		// Parse body
		var score GameScore
		if err := c.BodyParser(&score); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request",
			})
		}

		// Validate required fields
		if score.Email == "" || score.Score == 0 || score.Date == "" {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request",
			})
		}

		// Validate email
		if !validateEmail(score.Email) {
			return c.Status(400).JSON(fiber.Map{
				"status": "🔴 Bad Request, Invalid Email",
			})
		}

		// Insert score using raw SQL
		query := `INSERT INTO game_hof_noteDino (date, score, email, user_name) VALUES (?, ?, ?, ?)`
		if err := db.Exec(query, score.Date, score.Score, score.Email, score.UserName).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "🔴 Operation was unsuccessful!",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"gameScoreInfo": score,
			"status":        "🟢 Game score insertion was successful",
		})
	}
}

// GetNoteDinoHof handles getting top scores for NoteDino game
func GetNoteDinoHof(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var results []GameScore
		query := `SELECT date, score, email, user_name FROM game_hof_noteDino ORDER BY score DESC LIMIT 10`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "🔴 Error while fetching hof",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof": results,
		})
	}
}
