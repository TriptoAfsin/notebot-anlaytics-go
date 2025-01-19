package handler

import (
	"log"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetTopNoteSubjects handles fetching top 5 note subjects
func GetTopNoteSubjects(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GetTopNoteSubjects handler called")
	return func(c *fiber.Ctx) error {
		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM subnamedb ORDER BY count DESC LIMIT 5").Scan(&results).Error; err != nil {
			log.Printf("🔴 Error while retrieving top note subjects: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.Academic.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"topNoteSubjects": results,
		})
	}
}

// GetLabSubjects handles fetching all lab subjects
func GetLabSubjects(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GetLabSubjects handler called")
	return func(c *fiber.Ctx) error {
		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM labsdb").Scan(&results).Error; err != nil {
			log.Printf("🔴 Error while retrieving lab subjects: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.Academic.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"LabSubjects": results,
		})
	}
}

// GetTopLabSubjects handles fetching top 5 lab subjects
func GetTopLabSubjects(db *gorm.DB) fiber.Handler {
	log.Println("🟢 GetTopLabSubjects handler called")
	return func(c *fiber.Ctx) error {
		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM labsdb ORDER BY count DESC LIMIT 5").Scan(&results).Error; err != nil {
			log.Printf("🔴 Error while retrieving top lab subjects: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.Academic.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"topLabSubjects": results,
		})
	}
}

// IncrementSubjectCount is a generic handler for incrementing subject counts
func IncrementSubjectCount(db *gorm.DB, subject string) fiber.Handler {
	log.Println("🟢 IncrementSubjectCount handler called with subject: ", subject)
	return func(c *fiber.Ctx) error {
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": config.AppMessages.Academic.UnauthorizedAccess,
			})
		}

		if err := db.Exec("UPDATE subnamedb SET count = count + 1 WHERE sub_name = ?", subject).Error; err != nil {
			log.Printf("🔴 Error while updating count for %s: %v", subject, err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.Academic.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": config.AppMessages.Academic.OperationSuccessful,
		})
	}
}

// Subject-specific handlers using the generic IncrementSubjectCount
func NotesMath1(db *gorm.DB) fiber.Handler  { return IncrementSubjectCount(db, "math1") }
func NotesMath2(db *gorm.DB) fiber.Handler  { return IncrementSubjectCount(db, "math2") }
func NotesPhy1(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "phy1") }
func NotesPhy2(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "phy2") }
func NotesChem1(db *gorm.DB) fiber.Handler  { return IncrementSubjectCount(db, "chem1") }
func NotesChem2(db *gorm.DB) fiber.Handler  { return IncrementSubjectCount(db, "chem2") }
func NotesPse(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "pse") }
func NotesCp(db *gorm.DB) fiber.Handler     { return IncrementSubjectCount(db, "cp") }
func NotesNtf(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "ntf") }
func NotesEm(db *gorm.DB) fiber.Handler     { return IncrementSubjectCount(db, "em") }
func NotesBce(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "bce") }
func NotesAm1(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "am1") }
func NotesAm2(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "am2") }
func NotesYm1(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "ym1") }
func NotesYm2(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "ym2") }
func NotesFm1(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "fm1") }
func NotesFm2(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "fm2") }
func NotesWp1(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "wp1") }
func NotesWp2(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "wp2") }
func NotesStat(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "stat") }
func NotesFeee(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "feee") }
func NotesMarket(db *gorm.DB) fiber.Handler { return IncrementSubjectCount(db, "marketing") }
func NotesTtqc(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "ttqc") }
func NotesTp(db *gorm.DB) fiber.Handler     { return IncrementSubjectCount(db, "tp") }
func NotesMp(db *gorm.DB) fiber.Handler     { return IncrementSubjectCount(db, "mp") }
func NotesMmtf(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "mmtf") }
func NotesAcm(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "acm") }
func NotesTqm(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "tqm") }
func NotesFsd(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "fsd") }
func NotesAce(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "ace") }
func NotesMic(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "mic") }
func NotesSss1(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "sss1") }
func NotesSss2(db *gorm.DB) fiber.Handler   { return IncrementSubjectCount(db, "sss2") }
func NotesWpp(db *gorm.DB) fiber.Handler    { return IncrementSubjectCount(db, "wpp") }
func NotesEcono(db *gorm.DB) fiber.Handler  { return IncrementSubjectCount(db, "econo") }

// IncrementLabCount is a generic handler for incrementing lab counts
func IncrementLabCount(db *gorm.DB, subject string) fiber.Handler {
	log.Println("🟢 IncrementLabCount handler called with subject: ", subject)
	return func(c *fiber.Ctx) error {
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": config.AppMessages.Academic.UnauthorizedAccess,
			})
		}

		if err := db.Exec("UPDATE labsdb SET count = count + 1 WHERE lab_name = ?", subject).Error; err != nil {
			log.Printf("🔴 Error while updating count for lab %s: %v", subject, err)
			return c.Status(500).JSON(fiber.Map{"status": config.AppMessages.Academic.OperationUnsuccessful})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": config.AppMessages.Academic.OperationSuccessful,
		})
	}
}

// Lab-specific handlers using the generic IncrementLabCount
func LabsPhy1(db *gorm.DB) fiber.Handler  { return IncrementLabCount(db, "phy1") }
func LabsPhy2(db *gorm.DB) fiber.Handler  { return IncrementLabCount(db, "phy2") }
func LabsChem1(db *gorm.DB) fiber.Handler { return IncrementLabCount(db, "chem1") }
func LabsChem2(db *gorm.DB) fiber.Handler { return IncrementLabCount(db, "chem2") }
func LabsCP(db *gorm.DB) fiber.Handler    { return IncrementLabCount(db, "cp") }
func LabsBce(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "bce") }
func LabsMsp(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "msp") }
func LabsAm1(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "am1") }
func LabsAm2(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "am2") }
func LabsYm1(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "ym1") }
func LabsYm2(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "ym2") }
func LabsWp1(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "wp1") }
func LabsWp2(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "wp2") }
func LabsFm1(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "fm1") }
func LabsFm2(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "fm2") }
func LabsFeee(db *gorm.DB) fiber.Handler  { return IncrementLabCount(db, "feee") }
func LabsFme(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "fme") }
func LabsTtqc(db *gorm.DB) fiber.Handler  { return IncrementLabCount(db, "ttqc") }
func LabsAp1(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "ap1") }
func LabsAp2(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "ap2") }
func LabsMp(db *gorm.DB) fiber.Handler    { return IncrementLabCount(db, "mp") }
func LabsFsd(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "fsd") }
func LabsLss(db *gorm.DB) fiber.Handler   { return IncrementLabCount(db, "lss") }
