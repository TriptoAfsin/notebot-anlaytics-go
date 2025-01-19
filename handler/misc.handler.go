package handler

import (
	"log"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MissedWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// GetMissedWords handles fetching all missed words with pagination
func GetMissedWords(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GET: GetMissedWords handler called")
	return func(c *fiber.Ctx) error {
		// Get pagination parameters with defaults
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 500)

		// Prevent negative values
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 500 {
			limit = 500
		}

		offset := (page - 1) * limit

		var missedWords []map[string]interface{}
		var total int64

		// Get total count
		if err := db.Raw("SELECT COUNT(*) FROM missed_words_table").Scan(&total).Error; err != nil {
			log.Printf("🔴 Error while counting missed words: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.MissedWord.FetchError,
			})
		}

		// Get paginated results
		query := `
			SELECT * FROM missed_words_table 
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		if err := db.Raw(query, limit, offset).Scan(&missedWords).Error; err != nil {
			log.Printf("🔴 Error while fetching missed words: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.MissedWord.FetchError,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"missed_words": missedWords,
			"pagination": fiber.Map{
				"current_page": page,
				"limit":        limit,
				"total":        total,
				"total_pages":  (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}

// CreateMissedWord handles creating or updating missed word entries
func CreateMissedWord(db *gorm.DB) fiber.Handler {
	log.Println("🔵 POST: CreateMissedWord handler called")
	return func(c *fiber.Ctx) error {
		word := struct {
			Word string `json:"word"`
		}{}

		if err := c.BodyParser(&word); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.MissedWord.BadRequest,
			})
		}

		// Simple insert query without count column
		query := `
			INSERT INTO missed_words_table (missed_words) 
			VALUES (?)
		`

		if err := db.Exec(query, word.Word).Error; err != nil {
			log.Printf("🔴 Error while inserting missed word: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.MissedWord.OperationUnsuccessful,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"word":   word.Word,
			"status": config.AppMessages.MissedWord.InsertSuccess,
		})
	}
}
