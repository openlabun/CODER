package cache_publisher
import "encoding/json"

type SyncMessage struct {
	Model 	string          `json:"model"`
	Request string        	`json:"request"`
	Data  	json.RawMessage `json:"data"`
}