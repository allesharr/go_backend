package routers

import (
	"go_backend/db"
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
	db.GetConection().Create(aukt_router)
}
