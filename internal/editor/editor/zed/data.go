package zed

const name string = "zed"

type Zed struct{}

func New() *Zed {
	return &Zed{}
}

func (e *Zed) Name() string {
	return name
}
