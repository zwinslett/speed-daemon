package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zwinslett/speed-daemon/calculator"
	"github.com/zwinslett/speed-daemon/format"
	"github.com/zwinslett/speed-daemon/model"
)

func lastActivityCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "last",
		Short: "Display your last activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			detailedActivity, zones, err := fetchLastActivity(cmd.Context())
			if err != nil {
				return err
			}
			if asJSON {
				return format.PrintAsJSON(model.ActivityReport{
					Activity: detailedActivity,
					Zones:    calculator.AggregateZones(zones, calculator.Heartrate),
				})
			}
			fmt.Println(format.ActivityTable(detailedActivity))
			fmt.Println(format.SplitTable(detailedActivity))
			fmt.Println(format.ZonesTable(zones, calculator.Heartrate))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	return cmd
}
