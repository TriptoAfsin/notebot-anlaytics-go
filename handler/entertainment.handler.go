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
		var results []GameScore
		query := `SELECT date, score, email, user_name FROM game_hof ORDER BY score DESC LIMIT 10`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof": results,
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
		var results []GameScore
		query := `SELECT date, score, email, user_name FROM game_hof_noteDino ORDER BY score DESC LIMIT 10`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.Game.FetchError,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"hof": results,
		})
	}
}
