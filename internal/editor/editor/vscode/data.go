package vscode

const name string = "vscode"

type VsCode struct{}

func New() *VsCode {
	return &VsCode{}
}

func (e *VsCode) Name() string {
	return name
}
