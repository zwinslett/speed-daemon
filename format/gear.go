package format

import (
	"fmt"
	"strings"

	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/model"
)

func GearMessage(activities []model.DetailedActivity) string {
	var message strings.Builder
	aggregatedGear := calculator.AggregateGear(activities)
	for _, gear := range aggregatedGear {
		fmt.Fprintf(
			&message,
			"👟 %s 📍 %.2f\n",
			gear.Name,
			gear.Distance,
		)
	}
	return "<u><b>Shoe Rotation</b></u>\n" + message.String()
}
