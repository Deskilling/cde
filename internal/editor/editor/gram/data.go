package gram

const name string = "gram"

type Gram struct{}

func New() *Gram {
	return &Gram{}
}

func (e *Gram) Name() string {
	return name
}
