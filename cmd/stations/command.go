package stations

import (
	"fmt"

	"github.com/spf13/cobra"
)

const NY_STATE_JSON_API_ENDPOINT string = "https://data.ny.gov/resource/i9wp-a4ja.json?$limit=10000"

var Command = &cobra.Command{
	Use: "stations",
	Short: "Download station data",
	Long: "Download station data from NY State (https://data.ny.gov/Transportation/MTA-NYCT-Subway-Entrances-and-Exits-2015/i9wp-a4ja/about_data)",
	RunE: func(cmd *cobra.Command, args []string) error {
        endpointUrl, err := cmd.Flags().GetString("station_entrances_url")
        if err != nil {
            return fmt.Errorf("No value for 'station_entrances_url': %w", err)
        }

        baseDir, err := cmd.Flags().GetString("base_dir")
        if err != nil {
            return fmt.Errorf("No value for 'base_dir': %w", err)
        }        

        stationData, err := downloadStationData(endpointUrl)
        if err != nil {
            return fmt.Errorf("Unable to download station data: %w", err)
        }

        fmt.Printf("Found %d entries and exits\n", len(stationData))

        formattedStationData, err := formatStationData(stationData)
        if err != nil {
            return fmt.Errorf("Unable to parse station data: %w", err)
        }

        fmt.Printf("Parsed %d stations\n", len(formattedStationData))

        err = upsertFormattedStationData(baseDir, formattedStationData)
        if err != nil {
            return fmt.Errorf("Unable to upsert formatted station data: %w", err)
        }

        return nil
	},
}

func init() {
    Command.Flags().String("base_dir", "", "Where to output the finished data")
    Command.Flags().String("station_entrances_url", NY_STATE_JSON_API_ENDPOINT, "URI to JSON containing NY state subway entrances and exit data.")
}


