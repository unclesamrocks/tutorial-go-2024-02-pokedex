package json

import "encoding/json"

func Unmarshal(data []byte, project *any) error {
	error := json.Unmarshal(data, *project)
	if error != nil {
		return error
	}
	return nil
}
