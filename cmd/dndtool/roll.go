package dndtool

import (
	"fmt"
	"github.com/spf13/cobra"
	"go_rpg_sims_dnd_5e/roll"
	"math"
)

var sides int
var rolls int
var advantage bool
var disadvantage bool
var adjustment int
var dropLowest bool
var dropHighest bool

var rollCmd = &cobra.Command{
	Use:     "roll",
	Version: version,
	Short:   "roll - roll some dice",
	Long:    `roll is a simple CLI for rolling dice.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sides:        %d\n", sides)
		fmt.Printf("Rolls:        %d\n", rolls)
		fmt.Printf("Adjustment:   %d\n", adjustment)
		fmt.Printf("Advantage:    %t\n", advantage)
		fmt.Printf("Disadvantage: %t\n", disadvantage)
		fmt.Printf("Drop_lowest:  %t\n", dropLowest)
		fmt.Printf("Drop_highest: %t\n", dropHighest)
		var opts []string
		if advantage == true {
			opts = []string{"drop lowest 1"}
			rolls = 2
		} else if disadvantage == true {
			opts = []string{"drop highest 1"}
			rolls = 2
		} else if dropLowest {
			opts = []string{"drop lowest 1"}
		} else if dropHighest {
			opts = []string{"drop highest 1"}
		}
		if adjustment != 0 {
			var tDir string
			if adjustment > 0 {
				tDir = "add"
			} else {
				tDir = "subtract"
			}
			tStr := fmt.Sprintf("%s %d", tDir, int(math.Abs(float64(adjustment))))
			opts = append(opts, tStr)
		}
		rollObj, err := roll.Perform(sides, rolls, "dndtool roll", opts...)
		// rollObj, err := roll.Perform(sides, rolls, "dndtool roll")

		if err != nil {
			panic(err)
		}
		fmt.Println(rollObj.ToPrettyString())
	},
}

func init() {
	rollCmd.Flags().IntVarP(&sides, "sides", "s", 20,
		"Sides of the dice to roll")
	rollCmd.Flags().IntVarP(&rolls, "rolls", "r", 1,
		"Nbr of rolls to make")
	rollCmd.Flags().IntVarP(&adjustment, "adjustment", "j", 0,
		"Value to add to result")
	rollCmd.Flags().BoolVarP(&advantage, "advantage", "a",
		false, "Roll twice, keep the highest (overrides nbr of rolls)")
	rollCmd.Flags().BoolVarP(&disadvantage, "disadvantage", "d",
		false, "Roll twice, keep the lowest (overrides nbr of rolls)")
	rollCmd.Flags().BoolVarP(&dropLowest, "drop_lowest", "w",
		false, "Drop the lowest of the results (evals nbr of rolls -1)")
	rollCmd.Flags().BoolVarP(&dropHighest, "drop_highest", "g",
		false, "Drop the highest of the results (evals nbr of rolls -1)")
	rootCmd.AddCommand(rollCmd)
}
