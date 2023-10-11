package providers

import "github.com/priyankasharma10/ReNew/models"

type DBHelperProvider interface {
	CreateNewUser(newUserRequest *models.CreateNewUserRequest, userID int) (*int, error)
	IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error)
	UpdateSession(sessionId string) error
	FetchUserData(userID int) (models.FetchUserData, error)
	FetchUserSessionData(userID int) ([]models.FetchUserSessionsData, error)
	IsPhoneNumberAlreadyExist(phone string) (bool, error)
	GetUserInfoByEmail(email string) (models.GetUserDataByEmail, error)
	LogInUserUsingEmailAndRole(loginReq models.EmailAndPassword, role models.UserRoles) (userID int, message string, err error)
	StartNewSession(userID int, request *models.CreateSessionRequest) (string, error)
}
