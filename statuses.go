package tpb

// UserStatus represents the status of a suser
type UserStatus string

// List of all the available categories
const (
	Member    UserStatus = "member"
	VIP       UserStatus = "vip"
	Trusted   UserStatus = "trusted"
	Helper    UserStatus = "helper"
	Moderator UserStatus = "moderator"
	SuperMod  UserStatus = "supermod"
	Admin     UserStatus = "admin"
)
