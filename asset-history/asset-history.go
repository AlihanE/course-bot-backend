package asset_history

import "time"

type AssetHistory struct {
	Id string `json:"id" db:"id"`
	Data string `json:"data" db:"data"`
	CreateDate time.Time `json:"createDate" db:"create_date"`
}
