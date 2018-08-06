package services

type ServiceQuery struct {
	Term      string `form:"term"`
	Location  string `form:"location"`
	Latitude  string `form:"latitude" binding:"omitempty,latitude"`
	Longitude string `form:"longitude" binding:"omitempty,longitude"`
	Limit     int    `form:"limit" binding:"omitempty,gt=0"`
}
