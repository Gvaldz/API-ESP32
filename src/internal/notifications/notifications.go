package notifications

type Notification struct {
    Event     string      `json:"event"`
    Data      interface{} `json:"data"`
    CageID    int32       `json:"cage_id"`
    Timestamp int64       `json:"timestamp"`
}

const (
    EventNewMotion      = "new_motion"
    EventNewTemperature = "new_temperature"
    EventNewFoodStatus  = "new_food_status"
    EventNewHumidity    = "new_humidity"
)