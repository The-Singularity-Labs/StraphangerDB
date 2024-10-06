package stations

import (
	"fmt"
	"os"
	"path/filepath"
 
	"gopkg.in/yaml.v2"
	mapset "github.com/deckarep/golang-set/v2"
 )


 func (s FormattedStationData) Filepath(baseDir string) string {
	name, _ := s.aliases.Pop()
	line, _ := s.lineTags.Pop()
    return filepath.Join(
		baseDir, 
		fmt.Sprintf(
			"%s/stations/%s/%s.yaml", 
			normalize(s.Division), 
			normalize(line), 
			normalize(name),
		),
	)
}

func (s FormattedStationData) Merge(other FormattedStationData) FormattedStationData {
    s.aliases = mapset.NewSet[string](s.Aliases...).Union(mapset.NewSet[string](other.Aliases...))
	s.lineTags = mapset.NewSet[string](s.LineTags...).Union(mapset.NewSet[string](other.LineTags...))

	s.Aliases = s.aliases.ToSlice()
	s.LineTags = s.lineTags.ToSlice()
	s.Transfers = mapset.NewSet[string](s.Transfers...).Union(mapset.NewSet[string](other.Transfers...)).ToSlice()
    return s
}

func upsertFormattedStationData(baseDir string, data []FormattedStationData) error {
    for _, station := range data {
        fpath := station.Filepath(baseDir)
		
		fmt.Println(fpath)
        if _, err := os.Stat(fpath); err == nil {
            // File exists, load existing data
            var existingData FormattedStationData
            yamlFile, err := os.ReadFile(fpath)
            if err != nil {
                return fmt.Errorf("Unable to read existing YAML file: %w", err)
            }

            if err := yaml.Unmarshal(yamlFile, &existingData); err != nil {
                return fmt.Errorf("Unable to unmarshal existing YAML data: %w", err)
            }

            // Merge with current struct
            mergedData := existingData.Merge(station)

            // Serialize and write merged data
            yamlData, err := yaml.Marshal(&mergedData)
            if err != nil {
                return fmt.Errorf("Unable to marshal merged data to YAML: %w", err)
            }

            err = os.WriteFile(fpath, yamlData, 0644)
            if err != nil {
                return fmt.Errorf("Unable to write merged YAML data to file: %w", err)
            }
        } else {
			err = os.MkdirAll(filepath.Dir(fpath), 0777)
            if err != nil {
                return fmt.Errorf("Unable to create directory: %w", err)
            }
            // File doesn't exist, serialize and write directly
            yamlData, err := yaml.Marshal(&station)
            if err != nil {
                return fmt.Errorf("Unable to marshal station data to YAML: %w", err)
            }

            err = os.WriteFile(fpath, yamlData, 0777)
            if err != nil {
                return fmt.Errorf("Unable to write YAML data to file: %w", err)
            }
        }
    }

    return nil
}