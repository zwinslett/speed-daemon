package format

import (
	"fmt"
	"strings"

	"github.com/zwinslett/speed-daemon/calculator"
	"github.com/zwinslett/speed-daemon/model"
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
