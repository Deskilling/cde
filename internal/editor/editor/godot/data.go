package godot

const name string = "godot"

type Godot struct{}

func New() *Godot {
	return &Godot{}
}

func (e *Godot) Name() string {
	return name
}
