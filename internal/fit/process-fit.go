package fit

import (
	"fmt"
	"io"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

type FITMessage struct {
	Number          string
	Fields          []FITField
	DeveloperFields []FITDeveloperField
}

type FITField struct {
	Name   string
	Number uint8
	Value  any
	Units  string
	Scale  float64
	Offset float64
}

type FITDeveloperField struct {
	Number             uint8
	DeveloperDataIndex uint8
	Value              any
}

type ActivityMetrics struct {
	RecordCount int
	Messages    []FITMessage
}

func ProcessFIT(r io.Reader) error {

	// The listener will receive every message from the decoder
	// as soon as it is decoded and transform it into an filedef.File.
	lis := filedef.NewListener()
	defer lis.Close() // release channel used by listener

	dec := decoder.New(r,
		// Add activity listener to the decoder:
		decoder.WithMesgListener(lis),
		// Direct the decoder to only broadcast
		// the messages without retaining them:
		decoder.WithBroadcastOnly(),
	)

	_, err := dec.Decode()
	if err != nil {
		panic(err)
	}

	// The resulting File can be retrieved after decoding process completed.
	// filedef.File is just an interface, we can do type assertion like this.
	file, ok := lis.File().(*filedef.Activity)
	if !ok {
		fmt.Printf("%T is not an Activity File\n", lis.File())
		return fmt.Errorf("failed to assert Activity file")
	}

	fmt.Printf("File Type: %s\n", file.FileId.Type)
	fmt.Printf("Sessions count: %d\n", len(file.Sessions))
	fmt.Printf("Laps count: %d\n", len(file.Laps))
	fmt.Printf("Records count: %d\n", len(file.Records))
	i := 100
	fmt.Printf("\nSample value of record[%d]:\n", i)
	fmt.Printf("  Distance: %g m\n", file.Records[i].DistanceScaled())
	fmt.Printf("  Lat: %g degrees\n", file.Records[i].PositionLatDegrees())
	fmt.Printf("  Long: %g degrees\n", file.Records[i].PositionLongDegrees())
	fmt.Printf("  Speed: %g m/s\n", file.Records[i].SpeedScaled())
	fmt.Printf("  HeartRate: %d bpm\n", file.Records[i].HeartRate)
	fmt.Printf("  Cadence: %d rpm\n", file.Records[i].Cadence)
	// Output:
	// File Type: activity
	// Sessions count: 1
	// Laps count: 1
	// Records count: 3601
	//
	// Sample value of record[100]:
	//   Distance: 100 m
	//   Lat: 0 degrees
	//   Long: 0.0008982885628938675 degrees
	//   Speed: 1 m/s
	//   HeartRate: 126 bpm
	//   Cadence: 100 rpm

	// metrics := &ActivityMetrics{
	// 	RecordCount: len(fitFile.Messages),
	// 	Messages:    make([]FITMessage, 0, len(fitFile.Messages)),
	// }

	// for _, message := range fitFile.Messages {
	// 	parsedMessage := FITMessage{
	// 		Number: message.Num.String(),
	// 		Fields: make([]FITField, 0, len(message.Fields)),
	// 	}

	// 	for _, field := range message.Fields {
	// 		parsedMessage.Fields = append(parsedMessage.Fields, FITField{
	// 			Name:   field.Name,
	// 			Number: field.Num,
	// 			Value:  field.Value.Any(),
	// 			Units:  field.Units,
	// 			Scale:  field.Scale,
	// 			Offset: field.Offset,
	// 		})
	// 	}

	// 	metrics.Messages = append(metrics.Messages, parsedMessage)
	// }

	return nil
}
