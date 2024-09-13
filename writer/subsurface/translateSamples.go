package subsurface

import (
	"fmt"

	"littledev.nl/divelogimporter/models"
)

func translateSamples(samples []models.DiveSample) []SSSample {
	var ssSamples []SSSample
	for _, sample := range samples {
		ssSample := SSSample{
			Time: sample.Time,
		}

		if sample.Temperature != nil {
			ssSample.Temperature = fmt.Sprintf("%.1f C", *sample.Temperature)
		}

		if sample.Depth != nil {
			ssSample.Depth = fmt.Sprintf("%.2f m", *sample.Depth)
		}

		if len(sample.Pressure) > 0 {
			ssSample.Pressure = fmt.Sprintf("%.2f bar", sample.Pressure[0].Pressure)
		}

		ssSamples = append(ssSamples, ssSample)
	}
	return ssSamples
}

func translateEventName(name string) string {
	switch name {
	case "safety stop":
		fallthrough
	case "SafetyStop":
		return "Safety Stop"
	case "deepstop":
		fallthrough
	case "DeepStop":
		return "Deep Stop"
	case "heading":
		return "heading"
	}

	return ""
}

func translateFlag(translatedEventName string, modi map[string]bool) uint {
	switch translatedEventName {
	case "Safety Stop":
		fallthrough
	case "Deep Stop":
		isOn := modi[translatedEventName]
		modi[translatedEventName] = !isOn

		if isOn {
			return 10
		}
		return 9
	}
	return 0
}

func translateType(translatedEventName string) uint {
	switch translatedEventName {
	case "heading":
		return 23
	}
	return 0
}

func translateEvents(samples []models.DiveSample) []SSEvent {
	var ssEvents []SSEvent
	modi := make(map[string]bool)

	for _, sample := range samples {
		for _, event := range sample.Events {
			translatedEvent := translateEventName(event.Type)
			if translatedEvent == "" {
				continue
			}

			ssevent := SSEvent{
				Time: sample.Time,
				Name: translatedEvent,
				Type: translateType(translatedEvent),
				Flag: translateFlag(translatedEvent, modi),
			}

			if event.Value != 0 {
				ssevent.Value = fmt.Sprintf("%d", event.Value)
			}

			ssEvents = append(ssEvents, ssevent)
		}
	}
	return ssEvents
}

func translateTank(tank models.Tank) SSCylinder {
	cylinder := SSCylinder{
		Size:            "",
		WorkingPressure: "",
		Description:     "unknown",
		O2:              "",
		Start:           "",
		End:             "",
	}

	if tank.Volume != nil {
		cylinder.Size = fmt.Sprintf("%.1f l", float32(*tank.Volume))
		cylinder.WorkingPressure = "232.0 bar"
		cylinder.Description = fmt.Sprintf("%dℓ 232 bar", *tank.Volume)
	}

	if tank.Oxygen != nil {
		cylinder.O2 = fmt.Sprintf("%.1f%%", *tank.Oxygen)
	}

	if tank.Begin != nil {
		cylinder.Start = fmt.Sprintf("%.1f bar", *tank.Begin)
	}

	if tank.End != nil {
		cylinder.End = fmt.Sprintf("%.1f bar", *tank.End)
	}

	return cylinder
}
