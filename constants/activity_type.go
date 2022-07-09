package constants

const (
	LOGIN    = 1
	LOGOUT   = 2
	CHECKIN  = 3
	CHECKOUT = 4
)

var ActivityType map[int]string

func init() {
	ActivityType = map[int]string{
		1: "LOGIN",
		2: "LOGOUT",
		3: "CHECKIN",
		4: "CHECKOUT",
	}
}
