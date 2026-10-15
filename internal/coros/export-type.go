package coros

import (
	"fmt"
	"strings"
)

// ExportType is the file format accepted by the activity download endpoint.
type ExportType int

const (
	ExportTypeGPX ExportType = 1
	ExportTypeTCX ExportType = 3
	ExportTypeFIT ExportType = 4
)

func (e ExportType) String() string {
	switch e {
	case ExportTypeGPX:
		return "gpx"
	case ExportTypeTCX:
		return "tcx"
	case ExportTypeFIT:
		return "fit"
	default:
		return ""
	}
}

func ParseExportType(value string) (ExportType, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "gpx":
		return ExportTypeGPX, nil
	case "tcx":
		return ExportTypeTCX, nil
	case "fit":
		return ExportTypeFIT, nil
	default:
		return 0, fmt.Errorf("unsupported export type %q; expected fit, tcx or gpx", value)
	}
}
