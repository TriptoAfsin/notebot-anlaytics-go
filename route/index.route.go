package route

import (
	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RouteInit(app *fiber.App, db *gorm.DB) {

	app.Get("/", handler.ApiHandler)

	app.Get("/health", handler.HealthCheckHandler)

	// Daily report routes
	app.Get("/daily_report", handler.GetDailyReport(db))
	app.Post("/daily_report", handler.PostDailyReport(db))

	// NoteBird game routes
	app.Post("/games/notebird", handler.PostNoteBirdScore(db))
	app.Get("/games/notebird", handler.GetNoteBirdHof(db))

	// NoteDino game routes
	app.Post("/games/notedino", handler.PostNoteDinoScore(db))
	app.Get("/games/notedino", handler.GetNoteDinoHof(db))

	// Error logging routes
	app.Post("/logs/err", handler.PostNewError(db, config.GetAppConfig()))
	app.Post("/logs/err/email", handler.GetErrorsByEmail(db, config.GetAppConfig()))
	app.Get("/logs/err", handler.GetErrorLogs(db, config.GetAppConfig()))

	// User routes
	app.Post("/user/new", handler.CreateUser(db))
	app.Get("/users/app", handler.GetAllUsers(db))
	app.Post("/users/app/email", handler.GetUsersByEmail(db))
	app.Post("/users/app/batch_dept", handler.GetUsersByDeptAndBatch(db))

	// Missed words routes
	app.Get("/missed", handler.GetMissedWords(db))
	app.Post("/missed", handler.CreateMissedWord(db))

	// Notes routes
	app.Get("/notes", handler.GetTopNoteSubjects(db))
	app.Get("/notes/top", handler.GetTopNoteSubjects(db))

	// Individual subject routes
	app.Get("/notes/math1", handler.NotesMath1(db))
	app.Get("/notes/math2", handler.NotesMath2(db))
	app.Get("/notes/phy1", handler.NotesPhy1(db))
	app.Get("/notes/phy2", handler.NotesPhy2(db))
	app.Get("/notes/chem1", handler.NotesChem1(db))
	app.Get("/notes/chem2", handler.NotesChem2(db))
	app.Get("/notes/pse", handler.NotesPse(db))
	app.Get("/notes/cp", handler.NotesCp(db))
	app.Get("/notes/ntf", handler.NotesNtf(db))
	app.Get("/notes/bce", handler.NotesBce(db))
	app.Get("/notes/em", handler.NotesEm(db))
	app.Get("/notes/am1", handler.NotesAm1(db))
	app.Get("/notes/am2", handler.NotesAm2(db))
	app.Get("/notes/ym1", handler.NotesYm1(db))
	app.Get("/notes/ym2", handler.NotesYm2(db))
	app.Get("/notes/fm1", handler.NotesFm1(db))
	app.Get("/notes/fm2", handler.NotesFm2(db))
	app.Get("/notes/wp1", handler.NotesWp1(db))
	app.Get("/notes/wp2", handler.NotesWp2(db))
	app.Get("/notes/stat", handler.NotesStat(db))
	app.Get("/notes/market", handler.NotesMarket(db))
	app.Get("/notes/feee", handler.NotesFeee(db))
	app.Get("/notes/ttqc", handler.NotesTtqc(db))
	app.Get("/notes/tp", handler.NotesTp(db))
	app.Get("/notes/mp", handler.NotesMp(db))
	app.Get("/notes/mmtf", handler.NotesMmtf(db))
	app.Get("/notes/acm", handler.NotesAcm(db))
	app.Get("/notes/tqm", handler.NotesTqm(db))
	app.Get("/notes/fsd", handler.NotesFsd(db))
	app.Get("/notes/ace", handler.NotesAce(db))
	app.Get("/notes/mic", handler.NotesMic(db))
	app.Get("/notes/sss1", handler.NotesSss1(db))
	app.Get("/notes/sss2", handler.NotesSss2(db))
	app.Get("/notes/wpp", handler.NotesWpp(db))
	app.Get("/notes/econo", handler.NotesEcono(db))

	// Labs routes
	app.Get("/labs", handler.GetLabSubjects(db))
	app.Get("/labs/top", handler.GetTopLabSubjects(db))

	// We'll need to create these handlers in academic.handler.go
	app.Get("/labs/phy1", handler.LabsPhy1(db))
	app.Get("/labs/phy2", handler.LabsPhy2(db))
	app.Get("/labs/chem1", handler.LabsChem1(db))
	app.Get("/labs/chem2", handler.LabsChem2(db))
	app.Get("/labs/cp", handler.LabsCP(db))
	app.Get("/labs/bce", handler.LabsBce(db))
	app.Get("/labs/msp", handler.LabsMsp(db))
	app.Get("/labs/am1", handler.LabsAm1(db))
	app.Get("/labs/am2", handler.LabsAm2(db))
	app.Get("/labs/ym1", handler.LabsYm1(db))
	app.Get("/labs/ym2", handler.LabsYm2(db))
	app.Get("/labs/wp1", handler.LabsWp1(db))
	app.Get("/labs/wp2", handler.LabsWp2(db))
	app.Get("/labs/fm1", handler.LabsFm1(db))
	app.Get("/labs/fm2", handler.LabsFm2(db))
	app.Get("/labs/feee", handler.LabsFeee(db))
	app.Get("/labs/fme", handler.LabsFme(db))
	app.Get("/labs/ttqc", handler.LabsTtqc(db))
	app.Get("/labs/ap1", handler.LabsAp1(db))
	app.Get("/labs/ap2", handler.LabsAp2(db))
	app.Get("/labs/mp", handler.LabsMp(db))
	app.Get("/labs/fsd", handler.LabsFsd(db))
	app.Get("/labs/lss", handler.LabsLss(db))
}
