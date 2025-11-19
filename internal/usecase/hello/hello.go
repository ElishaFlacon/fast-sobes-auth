package hello

func (u *usecase) Hello() string {
	u.log.Infof("Hello")
	return "hello"
}
