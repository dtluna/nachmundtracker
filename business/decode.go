package business

import (
	"fmt"
	"io"
	"os"

	"github.com/dtluna/nachmundtracker/model"
	"github.com/goccy/go-yaml"
)

func DecodeData() ([]model.GameRecord, error) {
	file, err := os.Open("campaign.yaml")
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file data: %w", err)
	}

	var games []model.GameRecord

	if err := yaml.Unmarshal(data, &games); err != nil {
		return nil, fmt.Errorf("decoding yaml: %w", err)
	}

	return games, nil
}
