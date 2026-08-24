package selection

type Provider interface {
	Selection() ([]string, error)
}
