package components

type MeshComponent struct {
	Path string
}

func NewMeshComponent(path string) *MeshComponent {
	return &MeshComponent{
		Path: path,
	}
}

func (*MeshComponent) Type() string { return "mesh" }
