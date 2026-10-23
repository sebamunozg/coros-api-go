package coros

import "fmt"

type SportType struct {
	Key   string
	Value string
}

var AllSportTypes = map[string]string{
	"all":            "0",
	"run":            "100",
	"indoorRun":      "101",
	"trailRun":       "102",
	"trackRun":       "103",
	"hike":           "104",
	"mtnClimb":       "105",
	"bike":           "200",
	"indoorBike":     "201",
	"roadEbike":      "202",
	"gravelRoadBike": "203",
	"mountainRiding": "204",
	"mountainEbike":  "205",
	"helmetBike":     "299",
	"poolSwim":       "300",
	"openWater":      "301",
	"triathlon":      "10000",
	"strength":       "402",
	"gymCardio":      "400",
	"gpsCardio":      "401",
	"ski":            "500",
	"snowboard":      "501",
	"xcSki":          "502",
	"skiTouring":     "503",
	"skiTouringOld":  "10002",
	"multiSport":     "10001",
	"speedsurfing":   "706",
	"windsurfing":    "705",
	"row":            "700",
	"indoorRow":      "701",
	"whitewater":     "702",
	"flatwater":      "704",
	"multiPitch":     "10003",
	"climb":          "106",
	"indoorClimb":    "800",
	"bouldering":     "801",
	"walk":           "900",
	"jumpRope":       "901",
	"climbStairs":    "902",
	"customSport":    "98",
}

func GetSportTypeCode(key string) (string, error) {
	code, ok := AllSportTypes[key]
	if !ok {
		return "", fmt.Errorf("invalid sport type: %s", key)
	}

	return code, nil
}
