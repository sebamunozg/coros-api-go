package coros

type QueryActivitiesData struct {
	Count      int        `json:"count"`
	DataList   []DataList `json:"datalist"`
	PageNumber int        `json:"pageNumber"`
	TotalPage  int        `json:"totalPage"`
}

type QueryActivitiesResponse struct {
	CorosResponseBase
	Data QueryActivitiesData `json:"data"`
}

type DataList struct {
	AdjustedPace     int     `json:"adjustedPace"`
	Ascent           int     `json:"ascent"`
	Avg5x10s         int     `json:"avg5x10s"`
	AvgCadence       int     `json:"avgCadence"`
	AvgHr            int     `json:"avgHr"`
	AvgPower         int     `json:"avgPower"`
	AvgSpeed         float64 `json:"avgSpeed"`
	AvgStrkRate      int     `json:"avgStrkRate"`
	Best             int     `json:"best"`
	Best500m         int     `json:"best500m"`
	BestKm           int     `json:"bestKm"`
	BestLen          int     `json:"bestLen"`
	BodyTemperature  int     `json:"bodyTemperature"`
	Cadence          int     `json:"cadence"`
	Calorie          int     `json:"calorie"`
	Date             int     `json:"date"`
	Descent          int     `json:"descent"`
	Device           string  `json:"device"`
	DeviceID         string  `json:"deviceId"`
	DeviceSportMode  int     `json:"deviceSportMode"`
	Distance         float64 `json:"distance"`
	DownhillDesc     int     `json:"downhillDesc"`
	DownhillDist     float64 `json:"downhillDist"`
	DownhillTime     int     `json:"downhillTime"`
	EndTime          int64   `json:"endTime"`
	EndTimezone      int     `json:"endTimezone"`
	HasMessage       int     `json:"hasMessage"`
	HasPodDistance   int     `json:"hasPodDistance"`
	ImageURL         string  `json:"imageUrl"`
	ImageURLType     int     `json:"imageUrlType"`
	IsRunTest        int     `json:"isRunTest"`
	IsShowMs         int     `json:"isShowMs"`
	LabelID          string  `json:"labelId"`
	Lengths          int     `json:"lengths"`
	Max2s            int     `json:"max2s"`
	MaxSlope         int     `json:"maxSlope"`
	MaxSpeed         int     `json:"maxSpeed"`
	Mode             int     `json:"mode"`
	Name             string  `json:"name"`
	Np               int     `json:"np"`
	Pitch            int     `json:"pitch"`
	Sets             int     `json:"sets"`
	SpeedType        int     `json:"speedType"`
	SportType        int     `json:"sportType"`
	StartTime        int64   `json:"startTime"`
	StartTimezone    int     `json:"startTimezone"`
	Step             int     `json:"step"`
	SubMode          int     `json:"subMode"`
	Swolf            int     `json:"swolf"`
	Total            float64 `json:"total"`
	TotalDescent     int     `json:"totalDescent"`
	TotalFishingTime int     `json:"totalFishingTime"`
	TotalReps        int     `json:"totalReps"`
	TotalTime        int     `json:"totalTime"`
	TrainingLoad     int     `json:"trainingLoad"`
	UnitType         int     `json:"unitType"`
	WaterTemperature int     `json:"waterTemperature"`
	WorkoutTime      int     `json:"workoutTime"`
}
