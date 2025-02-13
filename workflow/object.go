package workflow

import (
	"encoding/json"
	"strings"
)

type Object map[string]any

func NewObject() Object {
	return make(Object)
}

func NewJSONObject(v []byte) (Object, error) {
	result := NewObject()

	err := json.Unmarshal(v, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c Object) With(key string, value any) Object {
	c[key] = value
	return c
}

func (c Object) Get(key string) any {
	return c[key]
}

func (c Object) GetInt(key string) (int, bool) {
	return GetObjectValue[int](c, key)
}

func (c Object) Int(key string) int {
	return first(c.GetInt(key))
}

func (c Object) GetBool(key string) (bool, bool) {
	return GetObjectValue[bool](c, key)
}

func (c Object) Bool(key string) bool {
	return first(c.GetBool(key))
}

func (c Object) GetString(key string) (string, bool) {
	return GetObjectValue[string](c, key)
}

func (c Object) String(key string) string {
	return first(c.GetString(key))
}

func (c Object) GetObject(key string) (Object, bool) {
	return GetObjectValue[Object](c, key)
}

func (c Object) Object(key string) Object {
	return first(c.GetObject(key))
}

func (c Object) GetObjectArray(key string) ([]Object, bool) {
	return GetObjectValue[[]Object](c, key)
}

func (c Object) ObjectArray(key string) []Object {
	return first(c.GetObjectArray(key))
}

func (c Object) GetStringArray(key string) ([]string, bool) {
	return GetObjectValue[[]string](c, key)
}

func (c Object) StringArray(key string) []string {
	return first(c.GetStringArray(key))
}

func (c Object) GetByPaths(paths []string) any {
	var (
		obj    = c
		result any
	)

	for i, path := range paths {
		if i == len(paths)-1 {
			result = obj.Get(path)
		} else {
			obj = obj.Object(path)
		}
	}

	return result
}

func (c Object) BatchOutput() []Object {
	return ObjectValue[[]Object](c, batchNodeOutputName)
}

func (c Object) Copy() Object {
	result := make(Object, len(c))
	for k, v := range c {
		result[k] = v
	}
	return result
}

func (c Object) JSON() ([]byte, error) {
	result, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Object) Pretty() string {
	return pretty(c)
}

func (c Object) PrettyIndent() string {
	return prettyIndent(c)
}

type ObjectField struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Node          string `json:"node"`
	Path          string `json:"path"`
	IsConstant    bool   `json:"is_constant"`
	Constant      any    `json:"constant"`
	IsBatchOutput bool   `json:"is_batch_output"`
}

func (c ObjectField) Paths() []string {
	return strings.Split(c.Path, ".")
}

type ObjectMapper []ObjectField

func GetObjectValue[T any](note Object, key string) (T, bool) {
	value, ok := note[key].(T)
	return value, ok
}

func ObjectValue[T any](note Object, key string) T {
	value, _ := note[key].(T)
	return value
}
