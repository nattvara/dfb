package groups

const (
	AllRepositories = "*"
)

type Repository struct {
	Name       string
	ResticPath string
}
