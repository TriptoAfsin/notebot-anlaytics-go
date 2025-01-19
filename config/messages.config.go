package config

// Messages holds all application messages
type Messages struct {
	Success    SuccessMessages
	Error      ErrorMessages
	Validation ValidationMessages
	Game       GameMessages
	ErrorLog   ErrorLogMessages
	MissedWord MissedWordMessages
	User       UserMessages
	API        APIMessages
	Academic   AcademicMessages
}

// SuccessMessages contains all success related messages
type SuccessMessages struct {
	Created   string
	Updated   string
	Deleted   string
	Retrieved string
}

// ErrorMessages contains all error related messages
type ErrorMessages struct {
	Internal      string
	NotFound      string
	Unauthorized  string
	BadRequest    string
	AlreadyExists string
	InvalidEmail  string
	FetchError    string
}

// ValidationMessages contains all validation related messages
type ValidationMessages struct {
	Required      string
	InvalidFormat string
	TooLong       string
	TooShort      string
}

// GameMessages contains all game related messages
type GameMessages struct {
	ScoreInsertSuccess    string
	UnauthorizedAccess    string
	BadRequest            string
	InvalidFields         string
	OperationUnsuccessful string
	FetchError            string
}

// ErrorLogMessages contains all error logging related messages
type ErrorLogMessages struct {
	BadRequest            string
	InvalidEmail          string
	OperationUnsuccessful string
	FetchError            string
	LogInsertSuccess      string
	LogsFetchSuccess      string
	EmailFetchError       string
}

// MissedWordMessages contains all missed word related messages
type MissedWordMessages struct {
	FetchError            string
	OperationUnsuccessful string
	BadRequest            string
	InsertSuccess         string
}

// UserMessages contains all user related messages
type UserMessages struct {
	UnauthorizedAccess    string
	BadRequest            string
	InsertSuccess         string
	FetchError            string
	CountError            string
	IncrementSuccess      string
	OperationUnsuccessful string
}

// APIMessages contains all API related messages
type APIMessages struct {
	UnauthorizedAccess    string
	BadRequest            string
	InvalidPlatform       string
	OperationUnsuccessful string
	OperationSuccessful   string
	LogCheckError         string
	UpdateCountError      string
	IncrementSuccess      string
	NewLogSuccess         string
}

// AcademicMessages contains all academic related messages
type AcademicMessages struct {
	UnauthorizedAccess    string
	OperationUnsuccessful string
	OperationSuccessful   string
	TopSubjectsError      string
	LabSubjectsError      string
	TopLabsError          string
	SubjectUpdateError    string
	LabUpdateError        string
}

// AppMessages is the global messages instance
var AppMessages = Messages{
	Success: SuccessMessages{
		Created:   "Resource created successfully",
		Updated:   "Resource updated successfully",
		Deleted:   "Resource deleted successfully",
		Retrieved: "Resource retrieved successfully",
	},
	Error: ErrorMessages{
		Internal:      "Internal server error occurred",
		NotFound:      "Resource not found",
		Unauthorized:  "Unauthorized access",
		BadRequest:    "Invalid request",
		AlreadyExists: "Resource already exists",
		InvalidEmail:  "🔴 Bad Request, Invalid Email",
		FetchError:    "🔴 Error while fetching hof",
	},
	Validation: ValidationMessages{
		Required:      "This field is required",
		InvalidFormat: "Invalid format",
		TooLong:       "Value is too long",
		TooShort:      "Value is too short",
	},
	Game: GameMessages{
		ScoreInsertSuccess:    "🟢 Game score insertion was successful",
		UnauthorizedAccess:    "🔴 Unauthorized Access !",
		BadRequest:            "🔴 Bad Request",
		InvalidFields:         "🔴 Bad Request - Invalid or missing fields",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
		FetchError:            "🔴 Error while fetching hof",
	},
	ErrorLog: ErrorLogMessages{
		BadRequest:            "🔴 Bad Request",
		InvalidEmail:          "🔴 Bad Request, Invalid Email",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
		FetchError:            "🔴 Error while fetching logs by email",
		LogInsertSuccess:      "🟢 New Error log insertion was successful",
		LogsFetchSuccess:      "🟢 Logs Data fetching was successful",
		EmailFetchError:       "🔴 Error while fetching logs by email",
	},
	MissedWord: MissedWordMessages{
		FetchError:            "🔴 Error while fetching missed words",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
		BadRequest:            "🔴 Bad Request",
		InsertSuccess:         "🟢 Word insertion was successful",
	},
	User: UserMessages{
		UnauthorizedAccess:    "🔴 Unauthorized Access !",
		BadRequest:            "🔴 Bad Request",
		InsertSuccess:         "🟢 New user info insertion was successful",
		FetchError:            "🔴 Error while fetching app users",
		CountError:            "🔴 Error while fetching app user count",
		IncrementSuccess:      "🟢 Incrementing user count was successful",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
	},
	API: APIMessages{
		UnauthorizedAccess:    "🔴 Unauthorized Access !",
		BadRequest:            "🔴 Bad Request",
		InvalidPlatform:       "🔴 Bad Request - Invalid Platform",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
		OperationSuccessful:   "🟢 Operation was successful",
		LogCheckError:         "🔴 Error while checking existing log!",
		UpdateCountError:      "🔴 Error while updating daily api count!",
		IncrementSuccess:      "🟢 Incrementing api call count was successful",
		NewLogSuccess:         "🟢 Creating new log entry was successful",
	},
	Academic: AcademicMessages{
		UnauthorizedAccess:    "🔴 Unauthorized Access !",
		OperationUnsuccessful: "🔴 Operation was unsuccessful!",
		OperationSuccessful:   "🟢 Operation was successful",
		TopSubjectsError:      "🔴 Error while retrieving top note subjects",
		LabSubjectsError:      "🔴 Error while retrieving lab subjects",
		TopLabsError:          "🔴 Error while retrieving top lab subjects",
		SubjectUpdateError:    "🔴 Error while updating count for subject",
		LabUpdateError:        "🔴 Error while updating count for lab",
	},
}
