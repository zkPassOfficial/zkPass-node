package typings

type Settings struct {
	Session *Session
}

type Session struct {
	MAX     int32
	LIFE    int32 //ms
	TIMEOUT int32 //ms
}
