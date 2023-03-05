package parser

// Resource is an alias for map[string]interface{}
type Resource map[string]interface{}

const (
	kind = "kind"
)

func (r Resource) Kind() (string, error) {
	k, ok := r[kind]
	if !ok {
		return "", ErrKindKeyNotFound
	}

	return k.(string), nil
}
