package http

// Request data
type RequestSaveVideoProgress struct {
	Timebar int `form:"timebar,default=0" json:"timebar" xml:"timebar"  binding:"number"`
}
