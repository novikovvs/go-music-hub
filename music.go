package main

import (
	"os"
)

type Track struct {
	Label string
}

type TrackLibrary struct {
	Children []Track
	Label    string
}

func getLibrary() []TrackLibrary {
	var result []Track
	var libraries []TrackLibrary
	entries, err := os.ReadDir("downloads/")
	if err != nil {
		AppLogger.Fatal(err.Error())
	}

	for _, e := range entries {
		AppLogger.Info(e.Name())
		result = append(result, Track{Label: e.Name()})
	}
	libraries = append(libraries, TrackLibrary{
		Children: result,
		Label:    "Домашняя библиотека",
	})

	return libraries
}
