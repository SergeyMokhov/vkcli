package auth

type UserPermissions struct{}

func AllUserPermissions() map[string]int32 {
	allUserPerm := map[string]int32{
		"notify":        1,
		"friends":       2,
		"photos":        4,
		"audio":         8,
		"video":         16,
		"stories":       64,
		"pages":         128,
		"+256":          256, //Addition of link to the application in the left menu.
		"status":        1024,
		"notes":         2048,
		"messages":      4096,
		"wall":          8192,
		"ads":           32768,
		"offline":       65536,
		"docs":          131072,
		"groups":        262144,
		"notifications": 524288,
		"stats":         1048576,
		"email":         4194304,
		"market":        134217728,
	}
	return allUserPerm
}

func FullUserScope() int32 {
	sum := int32(0)
	for _, i := range AllUserPermissions() {
		sum += i
	}
	return sum
}

func Friends() int32 {
	return AllUserPermissions()["friends"]
}
