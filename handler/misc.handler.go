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
		search := c.Query("search", "")

		// Prevent negative values
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 500 {
			limit = 500
		}

		offset := (page - 1) * limit

		// Build the WHERE clause for search
		whereClause := "1=1"
		params := []interface{}{}

		if search != "" {
			whereClause = "missed_words LIKE ?"
			searchPattern := "%" + search + "%"
			params = append(params, searchPattern)
		}

		var missedWords []map[string]interface{}
		var total int64

		// Get total count with search filter
		countQuery := "SELECT COUNT(*) FROM missed_words_table WHERE " + whereClause
		if err := db.Raw(countQuery, params...).Scan(&total).Error; err != nil {
			log.Printf("🔴 Error while counting missed words: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.MissedWord.FetchError,
			})
		}

		// Get paginated and filtered results
		query := `
			SELECT * FROM missed_words_table 
			WHERE ` + whereClause + `
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		// Add pagination parameters
		params = append(params, limit, offset)

		if err := db.Raw(query, params...).Scan(&missedWords).Error; err != nil {
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
				"search":       search,
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
