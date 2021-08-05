package joexzherror

type BizError struct {
	Err error
}

func (err BizError) Error() string {
	return err.Err.Error()
}
