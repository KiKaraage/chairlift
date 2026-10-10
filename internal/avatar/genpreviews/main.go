// Command genpreviews writes the embedded preview PNGs the picker loads.
//
// For every catalog entry it transcodes the embedded full illustration at
// PreviewSize and writes it to assets/preview/<base>.png (bluefin.webp ->
// bluefin.png). Run it with `go generate ./internal/avatar`, which runs this
// program from the package directory so assets/preview resolves inside the
// package. The output must exist before the package compiles: the //go:embed
// pattern fails the build when the directory is empty.
package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/projectbluefin/chairlift/internal/avatar"
)

func main() {
	const outDir = "assets/preview"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("genpreviews: creating %s: %v", outDir, err)
	}

	for _, entry := range avatar.Catalog() {
		webp, err := avatar.FetchAvatar(context.Background(), entry.ID)
		if err != nil {
			log.Fatalf("genpreviews: fetching %s: %v", entry.ID, err)
		}
		png, err := avatar.TranscodePreviewWebPToPNG(bytes.NewReader(webp))
		if err != nil {
			log.Fatalf("genpreviews: transcoding %s: %v", entry.ID, err)
		}
		name := strings.TrimSuffix(path.Base(entry.Path), ".webp") + ".png"
		out := filepath.Join(outDir, name)
		if err := os.WriteFile(out, png, 0o644); err != nil {
			log.Fatalf("genpreviews: writing %s: %v", out, err)
		}
		log.Printf("genpreviews: wrote %s (%d bytes)", out, len(png))
	}
}
