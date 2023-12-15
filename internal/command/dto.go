package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"requester/internal/model"
)

func requestDataFromCmd(cmd *cobra.Command) (*model.RequestData, error) {
	res := &model.RequestData{
		URL:       cmd.Flag(flagURL).Value.String(),
		Amount:    0,
		PerSecond: 0,
	}

	var err error

	res.Amount, err = cmd.Flags().GetInt(flagAmount)
	if err != nil {
		return nil, fmt.Errorf("parse flag %q value: %w", flagAmount, err)
	}

	res.PerSecond, err = cmd.Flags().GetInt(flagPerSecond)
	if err != nil {
		return nil, fmt.Errorf("parse flag %q value: %w", flagPerSecond, err)
	}

	return res, nil
}

func requestDataFromFile(file string) (*model.RequestData, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	res := &requestDataByConfig{}

	if err := yaml.NewDecoder(f).Decode(res); err != nil {
		return nil, fmt.Errorf("decode config from file: %w", err)
	}

	return res.toModel(), nil
}

type requestDataByConfig struct {
	URL      string `yaml:"url"`
	Requests struct {
		Amount    int `yaml:"amount"`
		PerSecond int `yaml:"per_second"`
	} `yaml:"requests"`
}

func (r *requestDataByConfig) toModel() *model.RequestData {
	return &model.RequestData{
		URL:       r.URL,
		Amount:    r.Requests.Amount,
		PerSecond: r.Requests.PerSecond,
	}
}
