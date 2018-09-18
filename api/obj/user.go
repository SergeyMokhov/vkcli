package obj

type User struct {
	Id          int         `json:"id"`
	FirstName   string      `json:"first_name"`
	Nickname    string      `json:"nickname"`
	LastName    string      `json:"last_name"`
	Deactivated deactivated `json:"deactivated"`
	Hidden      int         `json:"hidden"`

	BDate     string `json:"bdate"`
	HasMobile int    `json:"has_mobile"`
}

type deactivated string
