package fit

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

// ActivitySummary is the standardized summary of a single activity.
// Optional metrics are nil when the FIT file does not contain them
// (for example power for a swim), so they serialize as null instead of 0.
type ActivitySummary struct {
	LabelID string `json:"label_id,omitempty"`

	Sport     string     `json:"sport"`
	StartTime *time.Time `json:"start_time"`

	DistanceM     *float64 `json:"distance_m"`
	TimerTimeS    *float64 `json:"timer_time_s"`
	ElapsedTimeS  *float64 `json:"elapsed_time_s"`
	Calories      *uint16  `json:"calories_kcal"`
	AscentM       *uint16  `json:"ascent_m"`
	DescentM      *uint16  `json:"descent_m"`
	AvgHeartRate  *uint8   `json:"avg_heart_rate_bpm"`
	MaxHeartRate  *uint8   `json:"max_heart_rate_bpm"`
	AvgSpeedMS    *float64 `json:"avg_speed_ms"`
	MaxSpeedMS    *float64 `json:"max_speed_ms"`
	AvgPowerW     *uint16  `json:"avg_power_w"`
	AvgCadenceRPM *uint8   `json:"avg_cadence_rpm"`

	Laps    int `json:"laps"`
	Records int `json:"records"`
}

// ParseActivity decodes a FIT activity file and returns its summary.
// Metrics come from the first session of the file.
func ParseActivity(r io.Reader) (*ActivitySummary, error) {
	// The listener builds an filedef.File from the decoded messages.
	lis := filedef.NewListener()
	defer lis.Close()

	dec := decoder.New(r,
		decoder.WithMesgListener(lis),
		decoder.WithBroadcastOnly(),
	)

	if _, err := dec.Decode(); err != nil {
		return nil, fmt.Errorf("decoding FIT file: %w", err)
	}

	file, ok := lis.File().(*filedef.Activity)
	if !ok {
		return nil, fmt.Errorf("FIT file is not an activity file (%T)", lis.File())
	}
	if len(file.Sessions) == 0 {
		return nil, fmt.Errorf("FIT activity has no session")
	}

	s := file.Sessions[0]

	summary := &ActivitySummary{
		Sport:   s.Sport.String(),
		Laps:    len(file.Laps),
		Records: len(file.Records),
	}

	if !s.StartTime.IsZero() && s.StartTime.Year() > 1990 {
		t := s.StartTime.UTC()
		summary.StartTime = &t
	}

	summary.DistanceM = scaled32(s.TotalDistance, 100)
	summary.TimerTimeS = scaled32(s.TotalTimerTime, 1000)
	summary.ElapsedTimeS = scaled32(s.TotalElapsedTime, 1000)
	summary.Calories = valid16(s.TotalCalories)
	summary.AscentM = valid16(s.TotalAscent)
	summary.DescentM = valid16(s.TotalDescent)
	summary.AvgHeartRate = valid8(s.AvgHeartRate)
	summary.MaxHeartRate = valid8(s.MaxHeartRate)
	summary.AvgPowerW = valid16(s.AvgPower)
	summary.AvgCadenceRPM = valid8(s.AvgCadence)

	// Prefer the enhanced (32-bit) speed fields, fall back to the 16-bit ones.
	if summary.AvgSpeedMS = scaled32(s.EnhancedAvgSpeed, 1000); summary.AvgSpeedMS == nil {
		summary.AvgSpeedMS = scaled16(s.AvgSpeed, 1000)
	}
	if summary.MaxSpeedMS = scaled32(s.EnhancedMaxSpeed, 1000); summary.MaxSpeedMS == nil {
		summary.MaxSpeedMS = scaled16(s.MaxSpeed, 1000)
	}

	return summary, nil
}

// FIT marks "no value" with the maximum value of the field's base type.

func valid8(v uint8) *uint8 {
	if v == math.MaxUint8 {
		return nil
	}
	return &v
}

func valid16(v uint16) *uint16 {
	if v == math.MaxUint16 {
		return nil
	}
	return &v
}

func scaled16(v uint16, scale float64) *float64 {
	if v == math.MaxUint16 {
		return nil
	}
	f := float64(v) / scale
	return &f
}

func scaled32(v uint32, scale float64) *float64 {
	if v == math.MaxUint32 {
		return nil
	}
	f := float64(v) / scale
	return &f
}
