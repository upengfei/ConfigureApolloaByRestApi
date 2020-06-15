package utils


import (
"github.com/go-ini/ini"
"log"

)

type IniRead struct {
	F *ini.File
}

func NewFile(filePath string) *IniRead {
	var (
		cfg *ini.File
		err error
	)

	ir := new(IniRead)

	cfg, err = ini.Load(filePath)
	if err != nil {
		log.Fatal(err)
	}
	ir.F = cfg
	return ir
}

func (ir *IniRead) GetValue(section, key string) *ini.Key {
	return ir.F.Section(section).Key(key)
}
