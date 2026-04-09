package utils

import "encoding/json"

func StructData2Map(structObj interface{}) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(structObj)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(inrec, &inInterface); err != nil {
		return nil, err
	}
	return inInterface, nil
}
