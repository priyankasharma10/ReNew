package dbhelperprovider

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/priyankasharma10/ReNew/crypto"
	dbutil "github.com/priyankasharma10/ReNew/dbutils"
	"github.com/priyankasharma10/ReNew/models"
	"github.com/sirupsen/logrus"
	// "github.com/volatiletech/null"
)

func (dh *DBHelper) FetchUserData(userID int) (models.FetchUserData, error) {
	// language=sql
	SQL := `SELECT
		users.id,
		users.name,
		users.email,
		users.phone,
		users.address,
		users.pincode,
		users.city,
		users.country,
		users.aadharcard
	FROM
		users
	WHERE
		users.id = $1 AND users.archived_at IS NULL;`

	var fetchUserData models.FetchUserData
	err := dh.DB.Get(&fetchUserData, SQL, userID)
	fmt.Println(fetchUserData)
	fmt.Println(SQL)
	fmt.Println(userID)
	if err != nil {
		logrus.Errorf("FetchUserData: error getting user data: %v", err)
		return fetchUserData, err
	}
	return fetchUserData, nil
}

func (dh *DBHelper) GetUserInfoByEmail(email string) (models.GetUserDataByEmail, error) {
	//language=sql
	fmt.Println(email)
	SQL := `SELECT  id, name, role,email, phone, address, city, country,pincode ,aadharcard
			FROM users 
	 WHERE email = $1
	 AND archived_at IS NULL`
	var getUserDataByEmail models.GetUserDataByEmail
	err := dh.DB.Get(&getUserDataByEmail, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("GetUserInfoByEmail: error getting user data: %v", err)
		return getUserDataByEmail, err
	}
	if err == sql.ErrNoRows {
		return getUserDataByEmail, errors.New("email does not exist")
	}
	return getUserDataByEmail, nil
}

func (dh *DBHelper) FetchUserSessionData(userID int) ([]models.FetchUserSessionsData, error) {
	//language=sql
	SQL := `SELECT  id, user_id,end_time,  token
			FROM sessions
			WHERE user_id = $1`

	fetchUserSessionData := make([]models.FetchUserSessionsData, 0)
	err := dh.DB.Select(&fetchUserSessionData, SQL, userID)
	if err != nil {
		logrus.Errorf("FetchUserSessionData: error getting user session data from database: %v", err)
		return fetchUserSessionData, err
	}
	return fetchUserSessionData, nil
}

func (dh *DBHelper) UpdateSession(sessionId string) error {
	//language=sql
	SQL := `UPDATE sessions
    		SET end_time = $2
			WHERE token = $1`

	_, err := dh.DB.Exec(SQL, sessionId, time.Now().Add(1*time.Hour))
	if err != nil {
		logrus.Errorf("FetchUserSessionData: error getting user session data from database: %v", err)
		return err
	}
	return nil
}

func (dh *DBHelper) IsUserAlreadyExists(emailID string) (isUserExist bool, user models.UserData, err error) {
	//	language=sql
	SQL := `SELECT id
			FROM users
			WHERE email = lower($1)
			  AND archived_at IS NULL `

	err = dh.DB.Get(&user, SQL, emailID)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("isEmailAlreadyExist: unable to get user from email %v", err)
		return false, user, err
	}

	if err == sql.ErrNoRows {
		return false, user, nil
	}

	return true, user, nil
}
func (dh *DBHelper) IsPhoneNumberAlreadyExist(phone string) (bool, error) {
	// language=sql
	SQL := `SELECT count(*) > 0 
            FROM users
            WHERE archived_at IS NULL
            AND phone  = $1`

	var isPhoneAlreadyExist bool
	err := dh.DB.Get(&isPhoneAlreadyExist, SQL, phone)
	if err != nil {
		logrus.Errorf("IsPhoneNumberAlreadyExist: error getting whether phone exist: %v", err)
		return isPhoneAlreadyExist, err
	}

	return isPhoneAlreadyExist, nil
}

func (dh *DBHelper) CreateNewUser(newUserRequest *models.CreateNewUserRequest, userID int) (*int, error) {
	var newUserID int

	txErr := dbutil.WithTransaction(dh.DB, func(tx *sqlx.Tx) error {

		SQL := `INSERT INTO users
		(name, email, phone, password, address, country, pincode, city, role,aadharcard,created_at, id)
		VALUES (lower(trim($1)), $2, $3, $4, $5, $6, $7, $8, $9,$10, $11, $12)
		RETURNING id;`


    args := []interface{}{
		newUserRequest.Name,
		newUserRequest.Email.String,
		newUserRequest.Phone.String,
		newUserRequest.Password,
		newUserRequest.Address,
		newUserRequest.Country,
		newUserRequest.Pincode,
		newUserRequest.City,
		newUserRequest.Role,
		newUserRequest.Aadharcard,
		time.Now().UTC(),
		3,
    }


		fmt.Println("here i am printing argument -----------------------")

		for i, arg := range args {
			fmt.Printf("Arg %d: %v\n", i+1, arg)
		}

		fmt.Println("hey ------------------ i am here ")

		err := tx.Get(&newUserID, SQL, args...)
		if err != nil {
			logrus.Errorf("CreateNewUser: error creating user %v", err)

			// Additional print statements for debugging
			fmt.Println("SQL Query:", SQL)
			fmt.Println("Arguments:", args)

			return err
		}

		fmt.Println("xbsabdsavdh-------------------")

		// SQL = `INSERT INTO user_roles
		// 		 (user_id, role)
		// 		 VALUES ($1,$2)`

		// _, err = tx.Exec(SQL, newUserID, string(newUserRequest.Role))
		// if err != nil {
		// 	logrus.Errorf("CreateNewUser: error creating user roles err %v", err)
		// 	return err
		// }
		return nil
	})

	fmt.Println("fbhdgfsgdfhgf-----------------")

	if txErr != nil {
		logrus.Errorf("CreateNewUser: error in creating user: %v", txErr)
		return nil, txErr
	}

	fmt.Println("hello ----------------------------------")

	return &newUserID, nil
}

func (dh *DBHelper) LogInUserUsingEmailAndRole(loginReq models.EmailAndPassword, role models.UserRoles) (UserId int, message string, err error) {
	// language=SQL
	SQL := `SELECT 	id,   
					password
			FROM users
		WHERE email = $1
		AND archived_at IS NULL`

	var user = struct {
		UserId         int    `db:"id"`
		HashedPassword string `db:"password"`
	}{}
	// password := crypto.HashAndSalt(loginReq.Password)
	// _ = password
	passwrod := crypto.HashAndSalt(loginReq.Password)
	fmt.Println("hash password is ", loginReq.Password)
	fmt.Println("Bcrypt Hash:", passwrod)
	if err = dh.DB.Get(&user, SQL, loginReq.Email); err != nil && err != sql.ErrNoRows {
		logrus.Errorf("LogInUserUsingEmailAndRole: error while getting user %v", err)
		return UserId, "error getting user", err
	}

	isPasswordMatched := crypto.ComparePasswords(user.HashedPassword, loginReq.Password)

	if !isPasswordMatched {
		return UserId, "Password Not Correct", models.ErrorPasswordNotMatched
	}

	// var userRole models.UserRoles
	// SQL = `
	// 	SELECT
	// 		role
	// 	FROM user_roles
	// 	WHERE id = $1
	// 	  	  AND role = $2
	// 		  AND archived_at IS NULL
	// `

	// err = dh.DB.Get(&userRole, SQL, user.UserId, role)
	// if err != nil && err != sql.ErrNoRows {
	// 	logrus.Errorf("LogInUserUsingEmailAndRole: error while getting user role:  %v", err)
	// 	return UserId, "error getting user role", err
	// }
	// if err == sql.ErrNoRows {
	// 	return UserId, "user role not matched", errors.New("user does not have required access")
	// }

	return user.UserId, "", nil
}
func (dh *DBHelper) StartNewSession(userID int, request *models.CreateSessionRequest) (string, error) {

	// language=sql
	fmt.Println(userID)
	SQL := `INSERT INTO sessions 
    (user_id, start_time, end_time, platform, model_name, os_version, device_id, token) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING token, id`

	fmt.Println("------------------- i am here ")

	args := []interface{}{
		userID,
		time.Now(),
		time.Now().Add(1 * time.Hour),
		request.Platform,
		request.ModelName,
		request.OSVersion,
		request.DeviceID,
		uuid.New(),
	}

	fmt.Println("hey ------------------")

	type sessionDetails struct {
		Token     string     `db:"token"`
		SessionID int64      `db:"id"`
		UserID    int        `db:"user_id"`
		StartTime time.Time  `db:"start_time" json:"start_time" sql:"not null"`
		EndTime   *time.Time `db:"end_time" json:"end_time"`
		Platform  string     `db:"platform" json:"platform"`
		ModelName string     `db:"model_name" json:"model_name"`
		OSVersion string     `db:"os_version" json:"os_version"`
		DeviceID  string     `db:"device_id" json:"device_id"`
	}

	fmt.Println("hiiiiiiiiiiiiiiiiiiiiiiiiiii")

	var session sessionDetails

	err := dh.DB.Get(&session, SQL, args...)
	if err != nil {
		logrus.Errorf("StartNewSession: error while starting new session: %v\n", err)
		return session.Token, err
	}

	return session.Token, nil
}
