package objx

import (
	"strings"
)

func (r *Obj) Join(sep string) (string, error) {
	arr, err := r.StringList()
	if err != nil {
		return "", err
	}
	return strings.Join(arr, sep), nil
}
