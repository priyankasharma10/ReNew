package models

type CreateSessionRequest struct {
	// ID        int64      `db:"id" json:"id" sql:"primary_key"`
	// UserId    int        `db:"user_id" json:"user_id" sql:"references:users(id)"`
	Platform  string `db:"platform" json:"platform"`
	ModelName string `db:"model_name" json:"model_name"`
	OSVersion string `db:"os_version" json:"os_version"`
	DeviceID  string `db:"device_id" json:"device_id"`
	// StartTime time.Time  `db:"start_time" json:"start_time" sql:"not null"`
	// EndTime   *time.Time `db:"end_time" json:"end_time"`
}
