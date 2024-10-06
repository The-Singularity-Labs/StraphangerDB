package stations

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type StationData struct {
    Division          string  `json:"division"`
    Line              string  `json:"line"`
    Borough           string  `json:"borough"`
    StationName       string  `json:"station_name"`
    StationLatitude   string  `json:"station_latitude"`
    StationLongitude  string  `json:"station_longitude"`
    DaytimeRoutes     string  `json:"daytime_routes"`
    EntranceType      string  `json:"entrance_type"`
    Entry             string  `json:"entry"`
    ExitOnly          string  `json:"exit_only"`
    Vending           string  `json:"vending"`
    Staffing          string  `json:"staffing"`
    NorthSouthStreet  string  `json:"north_south_street"`
    EastWestStreet    string  `json:"east_west_street"`
    Corner            string  `json:"corner"`
    EntranceLatitude  string  `json:"entrance_latitude"`
    EntranceLongitude string  `json:"entrance_longitude"`
    EntranceGeoreference Georeference `json:"entrance_georeference"`
    StationGeoreference Georeference `json:"station_georeference"`
}

func downloadStationData(url string) ([]StationData, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("Unable to execute get request on URL: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Unable to read respone body: %w", err)
    }

    var stations []StationData
    err = json.Unmarshal(body, &stations)
    if err != nil {
        return nil, fmt.Errorf("Unable to parse respone body: %w", err)
    }

    return stations, nil
}