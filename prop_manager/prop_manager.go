package prop_manager

import (
	"encoding/json"
	"os"

	"go_backend/api"
	"go_backend/logger"
)

type AppProperties struct {
	Props api.PropertiesStruct
}

func (p *AppProperties) Read() error {
	dat, err := os.ReadFile("properties.json")

	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot read properties.json: "+err.Error())
		return err
	}

	var fileProps api.PropertiesStruct
	err = json.Unmarshal([]byte(dat), &fileProps)
	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot unmarshal properties: "+err.Error())
		return err
	}

	p.Props = fileProps
	return nil
}
