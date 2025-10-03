package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func UpdateDotenv(kv map[string]string) error {
	f, err := os.Open(".env")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		f.Close()
		return err
	}

	f.Close()

	f, err = os.Create(".env")
	if err != nil {
		return err
	}

	for line := range strings.Lines(string(data)) {
		var found bool
		for key, value := range kv {
			if strings.HasPrefix(line, fmt.Sprintf("%s=", key)) {
				fmt.Fprintf(f, "%s=%s\n", key, value)
				found = true
				break
			}
		}

		if !found {
			f.WriteString(line)
		}
	}

	return f.Close()
}
