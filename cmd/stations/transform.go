package stations

import (
	"slices"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type FormattedStationData struct {
	Id string `yaml:"id"`
	Accessible bool `yaml:"accessible"`
	aliases       mapset.Set[string]
	Aliases []string `yaml:"station_names"`
	lineTags      mapset.Set[string]  
	LineTags []string `yaml:"lines"`
	Transfers []string `yaml:"transfers"`
	Division          string  `yaml:"division"`
	Borough           string  `yaml:"borough"`
    StationLatitude   string  `yaml:"station_latitude"`
    StationLongitude  string  `yaml:"station_longitude"`
    HasPaymentVendingMachine           bool  `yaml:"has_payment_vending"`
    Staffing          string  `yaml:"staffing"`
    StationGeoreference Georeference `yaml:"station_georeference"`
	EntryExits []EntryExit `yaml:"entry_exits"`
}

type EntryExit struct {
    EntryOnly             bool  `yaml:"entry_only"`
    ExitOnly          bool  `yaml:"exit_only"`
    EntranceType      string  `yaml:"entrance_type"`
    NorthSouthStreet  string  `yaml:"north_south_street"`
    EastWestStreet    string  `yaml:"east_west_street"`
    Corner            string  `yaml:"corner"`
    EntranceLatitude  string  `yaml:"entrance_latitude"`
    EntranceLongitude string  `yaml:"entrance_longitude"`
    EntranceGeoreference Georeference `yaml:"entrance_georeference"`
}

func formatStationData(stationData []StationData) ([]FormattedStationData, error) {
	cache := map[string]FormattedStationData{}
	for _, rawStation := range stationData {
		id := formatId(rawStation.Line, rawStation.StationName)
		existingStation, exists := cache[id]
		if !exists {
			cache[id] = FormattedStationData{
				Id: formatId(rawStation.Line, rawStation.StationName),
				Accessible: rawStation.EntranceType == "Elevator",
				aliases: mapset.NewSet[string](rawStation.StationName),
				lineTags: mapset.NewSet[string](normalizeLine(rawStation.Line)),
				Division: normalizeDivision(rawStation.Division, rawStation.Line),
				Borough: normalizeBorough(rawStation.Borough), 
				StationLatitude: rawStation.StationLatitude, 
				StationLongitude: rawStation.StationLongitude, 
				HasPaymentVendingMachine: toBool(rawStation.Vending),
				Staffing: rawStation.Staffing,
				StationGeoreference: rawStation.StationGeoreference,
				EntryExits: []EntryExit{{
					EntryOnly: toBool(rawStation.Entry),
					ExitOnly: toBool(rawStation.ExitOnly),
					EntranceType: rawStation.EntranceType,
					NorthSouthStreet: rawStation.NorthSouthStreet,
					EastWestStreet: rawStation.EastWestStreet,
					Corner: rawStation.Corner,
					EntranceLatitude: rawStation.EntranceLatitude,
					EntranceLongitude: rawStation.EntranceLongitude,
					EntranceGeoreference: rawStation.EntranceGeoreference,
				}},			
			}
		} else {
			existingStation.EntryExits = append(
				existingStation.EntryExits, 
				EntryExit{
					EntryOnly: toBool(rawStation.Entry),
					ExitOnly: toBool(rawStation.ExitOnly),
					EntranceType: rawStation.EntranceType,
					NorthSouthStreet: rawStation.NorthSouthStreet,
					EastWestStreet: rawStation.EastWestStreet,
					Corner: rawStation.Corner,
					EntranceLatitude: rawStation.EntranceLatitude,
					EntranceLongitude: rawStation.EntranceLongitude,
					EntranceGeoreference: rawStation.EntranceGeoreference,
				},			
			)

			existingStation.aliases.Add(rawStation.StationName)
			existingStation.lineTags.Add(normalizeLine(rawStation.Line))
			existingStation.Accessible = existingStation.Accessible || rawStation.EntranceType == "Elevator"
			cache[id] = existingStation
		}		
	}


	results := []FormattedStationData{}
	for _, formattedStation := range cache {
		formattedStation.Aliases = formattedStation.aliases.ToSlice()
		formattedStation.LineTags = formattedStation.lineTags.ToSlice()
		results = append(results, formattedStation)
	}
	return results, nil
}

func toBool(input string) bool {
	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(input)), "y") {
		return true
	}

	return false
}

func newStringSet(input string) (result mapset.Set[string]) {
	result = mapset.NewSet[string]()
	result.Add(input)

	return result
}

func formatId(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		parts = append(
			parts,
			normalize(arg),
		)
	}

	return strings.Join(parts, "_")
}

func normalize(input string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(input), " - ", "_"), "-", "_"), ".", ""), " ", "_")
}

func normalizeDivision(division string, line string) string {
	if slices.Contains([]string{"Archer Av", "63rd St"}, line) {
		return "NYCT"
	}
	parts := strings.Split(division, "/")
	return parts[0]
}

func normalizeLine(line string) string {
	if strings.Contains(line, "/") {
		return "transfer"
	}

	return line
}

func normalizeBorough(borough string) string {
	switch borough {
	
	case "M":
		return "Manhattan"
	case "Q":
		return "Queens"
	case "B":
		return "Brooklyn"
	case "Bx":
		return "Bronx"
	default:
		return "invalid"
	}
}