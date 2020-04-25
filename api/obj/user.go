package obj

type User struct {
	Id              int         `json:"id"`
	FirstName       string      `json:"first_name"`
	Nickname        string      `json:"nickname"`
	LastName        string      `json:"last_name"`
	Deactivated     deactivated `json:"deactivated"`
	IsClosed        bool        `json:"is_closed"`
	CanAccessClosed bool        `json:"can_access_closed"`
	Hidden          int         `json:"hidden"`
	Online          int         `json:"online"`
	LastSeen        LastSeen    `json:"last_seen"`

	BDate     string `json:"bdate"`
	HasMobile int    `json:"has_mobile"`
}

type LastSeen struct {
	//Unix Time
	Time     int64 `json:"time"`
	Platform int   `json:"platform"`
}

type deactivated string

const (
	UserDeleted deactivated = "deleted"
	UserBanned  deactivated = "banned"
)
