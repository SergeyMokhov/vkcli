package obj

type User struct {
	Id          int         `json:"id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Nickname    string      `json:"nickname"`
	Deactivated deactivated `json:"deactivated"`
	BDate       string      `json:"bdate"`
	HasMobile   int         `json:"has_mobile"`
}

type deactivated string
