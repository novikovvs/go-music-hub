package main

import "os"

type Track struct {
	Label string
	Path  string
}

type TrackLibrary struct {
	Children []Track
	Label    string
}

func getLibrary() []TrackLibrary {
	var libraries []TrackLibrary

	libraries = append(libraries, TrackLibrary{
		Children: getTracks(),
		Label:    "Домашняя библиотека",
	})

	return libraries
}

func getTracks() []Track {
	var result []Track

	entries, err := os.ReadDir("downloads/")
	if err != nil {
		AppLogger.Fatal(err.Error())
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		result = append(result, Track{
			Label: e.Name(),
			Path:  "./downloads/" + e.Name(),
		})
	}

	return result
}
