package routers

import (
	"go_backend/db"
	"go_backend/prop_manager"
	"testing"
)

func TestInsert(t *testing.T) {
	aukt_router := Aukt_Table_Row{
		Number:     "1",
		Seller:     "arr",
		Object:     "Balaclava",
		WhoGaveMax: "",
		Money:      1000,
		TimeToOut:  800,
	}
	props := prop_manager.AppProperties{}
	err := props.Read()
	if err != nil {
		t.Error("Cannot get props")
	}
	if !db.InitDB(props.Props) || !db.XormEnabled() {
		t.Error("Cannot enable database")
	}
	db.GetConection().Create(&aukt_router)
}
