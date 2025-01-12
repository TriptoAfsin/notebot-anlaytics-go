package handler

import (
	"log"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	Email  string `json:"email"`
	UniID  string `json:"uni_id"`
	Batch  string `json:"batch"`
	Dept   string `json:"dept"`
	Role   string `json:"role"`
	ImgUrl string `json:"imgUrl"`
}

// CreateUser handles creating new users
func CreateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": "🔴 Unauthorized Access !",
			})
		}

		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		// Validate required fields
		if user.Email == "" || user.UniID == "" || user.Batch == "" || user.Dept == "" || user.Role == "" {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		// Set default image URL if not provided
		if user.ImgUrl == "" {
			user.ImgUrl = "not given"
		}

		// Raw SQL query
		query := `INSERT INTO app_users (email, uni_id, batch, dept, role, img_url) VALUES (?, ?, ?, ?, ?, ?)`
		if err := db.Exec(query, user.Email, user.UniID, user.Batch, user.Dept, user.Role, user.ImgUrl).Error; err != nil {
			log.Printf("🔴 Error while inserting new user info: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Operation was unsuccessful!"})
		}

		return c.Status(200).JSON(fiber.Map{
			"user":   user,
			"status": "🟢 New user info insertion was successful",
		})
	}
}

// GetAllUsers handles fetching all users
func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get pagination parameters from query
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 500)
		offset := (page - 1) * limit

		// Get total count
		var total int64
		if err := db.Raw("SELECT COUNT(*) FROM app_users").Scan(&total).Error; err != nil {
			log.Printf("🔴 Error while counting users: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while counting users"})
		}

		// Get paginated users
		var users []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_users LIMIT ? OFFSET ?", limit, offset).Scan(&users).Error; err != nil {
			log.Printf("🔴 Error while fetching all users: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while fetching all users"})
		}

		return c.Status(200).JSON(fiber.Map{
			"users":        users,
			"total":        total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
		})
	}
}

// GetUserCount handles fetching user count
func GetUserCount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM app_users").Scan(&count).Error; err != nil {
			log.Printf("🔴 Error while fetching app user count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while fetching app user count"})
		}

		return c.Status(200).JSON(fiber.Map{"app_users_count": count})
	}
}

// GetUsersByEmail handles fetching users by email
func GetUsersByEmail(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{"error": "🔴 Unauthorized Access !"})
		}

		email := struct {
			Email string `json:"email"`
		}{}

		if err := c.BodyParser(&email); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		var users []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_users WHERE email LIKE ?", email.Email).Scan(&users).Error; err != nil {
			log.Printf("🔴 Error while fetching users by email: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while fetching app users"})
		}

		return c.Status(200).JSON(fiber.Map{"searched_users": users})
	}
}

// GetUsersByDeptAndBatch handles fetching users by department and batch
func GetUsersByDeptAndBatch(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{"error": "🔴 Unauthorized Access !"})
		}

		filter := struct {
			Dept  string `json:"dept"`
			Batch string `json:"batch"`
		}{}

		if err := c.BodyParser(&filter); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		var users []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_users WHERE batch = ? AND dept LIKE ? ORDER BY batch DESC",
			filter.Batch, filter.Dept).Scan(&users).Error; err != nil {
			log.Printf("🔴 Error while fetching filtered users: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while fetching app users"})
		}

		return c.Status(200).JSON(fiber.Map{"searched_users": users})
	}
}

// IncrementUserCount handles incrementing user count
func IncrementUserCount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Exec("UPDATE user_count SET count = count + 1").Error; err != nil {
			log.Printf("🔴 Error while incrementing user count: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while incrementing user count"})
		}

		return c.Status(200).JSON(fiber.Map{"status": "🟢 Incrementing user count was successful"})
	}
}
