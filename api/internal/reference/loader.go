package reference

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	filePath string
}

func NewLoader(filePath string) *Loader {
	return &Loader{
		filePath: filePath,
	}
}

func (l *Loader) LoadIngredientData(logger *slog.Logger) ([]*IngredientModel, error) {
	logger.Info("Loading ingredient data from file", "filePath", l.filePath)
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		logger.Error("Failed to read ingredient data file", "error", err)
		return nil, err
	}

	var fileData FileData
	if err := yaml.Unmarshal(data, &fileData); err != nil {
		logger.Error("Failed to unmarshal ingredient data", "error", err)
		return nil, err
	}

	logger.Info("File ingredient data", "ingredients", fmt.Sprintf("%+v", fileData.Ingredients))

	return fileData.Ingredients, nil
}
