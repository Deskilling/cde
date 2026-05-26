package vscodium

const name string = "vscodium"

type VsCodium struct{}

func New() *VsCodium {
	return &VsCodium{}
}

func (vscodium *VsCodium) Name() string {
	return name
}
