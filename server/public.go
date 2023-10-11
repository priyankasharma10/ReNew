package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/priyankasharma10/ReNew/models"
	authprovider "github.com/priyankasharma10/ReNew/providers/authProvider"
	"github.com/priyankasharma10/ReNew/scmerrors"
	"github.com/priyankasharma10/ReNew/utils"
	"github.com/sirupsen/logrus"
	"github.com/ttacon/libphonenumber"
	"github.com/volatiletech/null"
)

func (srv *Server) register(resp http.ResponseWriter, req *http.Request) {
	var newUserReq models.CreateNewUserRequest
	uc := srv.MiddlewareProvider.UserFromContext(req.Context())

	err := json.NewDecoder(req.Body).Decode(&newUserReq)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "error creating user", "Error parsing request")
		return
	}

	if newUserReq.Email.String == "" {
		scmerrors.RespondClientErr(resp, errors.New("email cannot be empty"), http.StatusBadRequest, "Email  cannot be empty", "Email cannot be empty")
		return
	}

	name := strings.TrimSpace(newUserReq.Name)
	if name == "" {
		scmerrors.RespondClientErr(resp, errors.New("name cannot be empty"), http.StatusBadRequest, "Name cannot be empty", "Name cannot be empty")
		return
	}
	// checking if the user is already exist
	isUserExist, _, err := srv.DBHelper.IsUserAlreadyExists(newUserReq.Email.String)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "error in processing request")
		return
	}

	if isUserExist {
		scmerrors.RespondClientErr(resp, errors.New("error creating user"), http.StatusBadRequest, "this email is already linked with one of our account please use a different email address", "unable to create a user with duplicate email address")
		return
	}

	newUserReq.Email.String = strings.ToLower(newUserReq.Email.String)
	if newUserReq.Name == "" {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "name cannot be empty", "name cannot be empty")
		return
	}

	if newUserReq.Phone.String == "" {
		scmerrors.RespondClientErr(resp, errors.New("phone number cannot be empty"), http.StatusBadRequest, "Phone number cannot be empty", "Phone number cannot be empty")
		return
	}

	if newUserReq.Address == "" {
		scmerrors.RespondClientErr(resp, errors.New("address cannot be empty"), http.StatusBadRequest, "password cannot be empty", "address cannot be empty")
		return
	}

	if newUserReq.Pincode == 0 {
		scmerrors.RespondClientErr(resp, errors.New("pincode cannot be empty"), http.StatusBadRequest, "pincode cannot be empty", "pincode cannot be empty")
		return
	}

	if newUserReq.City == "" {
		scmerrors.RespondClientErr(resp, errors.New("city  cannot be empty"), http.StatusBadRequest, "city cannot be empty", "city cannot be empty")
		return
	}

	if newUserReq.Country == "" {
		scmerrors.RespondClientErr(resp, errors.New("country  cannot be empty"), http.StatusBadRequest, "country cannot be empty", "country cannot be empty")
		return
	}

	uncleanPhoneNumber := newUserReq.Phone.String

	if strings.Count(uncleanPhoneNumber, "+") == 2 {
		uncleanPhoneNumber = uncleanPhoneNumber[strings.LastIndex(uncleanPhoneNumber, "+"):]
	}

	phone := strings.ReplaceAll(uncleanPhoneNumber, " ", "")

	num, err := libphonenumber.Parse(phone, "IN")
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Phone Number not in a correct format", "invalid format for phone number")
		return
	}

	isValidNumber := libphonenumber.IsValidNumber(num)

	if !isValidNumber {
		scmerrors.RespondClientErr(resp, errors.New("invalid phone number"), http.StatusBadRequest, "invalid phone number", "invalid phone number")
		return
	}

	phoneNumber := libphonenumber.Format(num, libphonenumber.E164)
	newUserReq.Phone = null.StringFrom(phoneNumber)

	// if !crypto.IsGoodPassword(newUserReq.Password) {
	// 	scmerrors.RespondClientErr(resp, errors.New("password length should be at least 6"), http.StatusBadRequest, "password length should be at least 6", "password length should be at least 6")
	// 	return
	// }

	isMobileAlreadyExist, err := srv.DBHelper.IsPhoneNumberAlreadyExist(phoneNumber)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "unable to create user")
		return
	}

	if isMobileAlreadyExist {
		scmerrors.RespondClientErr(resp, errors.New("phone number already exist"), http.StatusBadRequest, "this phone number is already linked with one of our account please use a different phone number", "unable to create a user")
		return
	}

	if uc.Role == string(models.Admin) && newUserReq.Role == models.Admin {
		scmerrors.RespondClientErr(resp, errors.New("admins cannot create admins"), http.StatusBadRequest, "admins cannot create admins. Only super admins can create admins", "country cannot be empty")
		return
	}

	// Creating user in the database
	userID, err := srv.DBHelper.CreateNewUser(&newUserReq, uc.UserId)
	if err != nil {
		scmerrors.RespondGenericServerErr(resp, err, "Error registering new user")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"message": "success",
		"userId":  userID,
	})
}

func (srv *Server) loginWithEmailPassword(resp http.ResponseWriter, req *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"function": "loginWithEmailPassword",
	})

	var token string
	var authLoginRequest models.AuthLoginRequest
	err := json.NewDecoder(req.Body).Decode(&authLoginRequest)
	if err != nil {
		log.WithError(err).Error("Unable to decode request body")
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Error decoding request body", "error decoding request body")
		return
	}

	log.Infof("Received login request for email: %s", authLoginRequest.Email)

	if authLoginRequest.Password == "" {
		log.Warn("Empty password received")
		scmerrors.RespondClientErr(resp, errors.New("password can not be empty"), http.StatusBadRequest, "Empty password!", "password field can not be empty")
		return
	}

	if authLoginRequest.Email == "" {
		log.Warn("Empty email received")
		scmerrors.RespondClientErr(resp, errors.New("email can not be empty"), http.StatusBadRequest, "Please enter email to login", "email can not be empty")
		return
	}

	UserDataByEmail, err := srv.DBHelper.GetUserInfoByEmail(authLoginRequest.Email)
	if err != nil {
		log.WithError(err).Error("Error getting user info by email")
		log.Errorf("Error details: %v", err) // Log the entire error for troubleshooting
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest, "Error getting user info", "error getting user info")
		return
	}

	log.Infof("User info retrieved for email: %s", UserDataByEmail.Email)

	// if authLoginRequest.Role != string(models.Admin) {
	// 	log.Warn("Invalid user role")
	// 	scmerrors.RespondClientErr(resp, errors.New("invalid user role"), http.StatusBadRequest, "Error role does not match", "error role does not match")
	// 	return
	// }

	loginReq := models.EmailAndPassword{
		Role:     authLoginRequest.Role,
		Email:    authLoginRequest.Email,
		Password: authLoginRequest.Password,
	}
	loginReq.Email = strings.ToLower(loginReq.Email)

	userID, errorMessage, err := srv.DBHelper.LogInUserUsingEmailAndRole(loginReq, UserDataByEmail.Role)
	fmt.Println("this is userid", userID)
	if err != nil {
		log.WithError(err).Error("Error logging in user")
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, errorMessage, errorMessage)
		return
	}

	createSessionRequest := models.CreateSessionRequest{
		Platform:  authLoginRequest.Platform,
		ModelName: authLoginRequest.ModelName.String,
		OSVersion: authLoginRequest.OSVersion.String,
		DeviceID:  authLoginRequest.DeviceID.String,
	}

	UUIDToken, err := srv.DBHelper.StartNewSession(userID, &createSessionRequest)
	if err != nil {
		log.WithError(err).Error("Error creating session")
		scmerrors.RespondGenericServerErr(resp, err, "Error in creating session")
		return
	}

	userInfo, err := srv.DBHelper.FetchUserData(userID)
	if err != nil {
		log.WithError(err).Error("Error getting user info")
		scmerrors.RespondGenericServerErr(resp, err, "Error in getting user info")
		return
	}

	devClaims := make(map[string]interface{})
	devClaims["UUIDToken"] = UUIDToken
	devClaims["userInfo"] = UserDataByEmail
	devClaims["UserSession"] = createSessionRequest

	token, err = authprovider.GenerateJWT(devClaims)
	if err != nil {
		log.WithError(err).Error("Error generating JWT")
		scmerrors.RespondClientErr(resp, err, http.StatusInternalServerError, "Error while login", "Error while login")
		return
	}

	log.Info("Login successful")

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"userInfo": userInfo,
		"token":    token,
	})
}
