package certutils

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func GetASCertificateFilename(configDir, isdAs string) (string, error) {
	parts := strings.Split(isdAs, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid ISD-AS format")
	}

	asnInFile := strings.ReplaceAll(parts[1], ":", "_")
	return filepath.Join(configDir, "crypto", "as", fmt.Sprintf("ISD%s-AS%s.pem", parts[0], asnInFile)), nil
}

func GetASPrivateKeyFilename(configDir string) string {
	return filepath.Join(configDir, "crypto", "as", "cp-as.key")
}

func GetTwoLatestTRCsForISD(listFiles func(string, string, string) ([]string, error), configDir, isd string) (string, string, error) {
	var laterstTRC string
	var penultimateTRC string
	trcDir := filepath.Join(configDir, "certs")
	trcFiles, err := listFiles(trcDir, "ISD"+isd, ".trc")
	if err != nil {
		return "", "", err
	}
	if len(trcFiles) == 0 {
		return "", "", fmt.Errorf("No TRC files found in %s", trcDir)
	}
	sort.Slice(trcFiles, func(i, j int) bool {
		filenameI := filepath.Base(trcFiles[i])
		filenameJ := filepath.Base(trcFiles[j])

		var bI, sI, bJ, sJ, isdI, isdJ int
		fmt.Sscanf(filenameI, "ISD%d-B%d-S%d.trc", &isdI, &bI, &sI)
		fmt.Sscanf(filenameJ, "ISD%d-B%d-S%d.trc", &isdJ, &bJ, &sJ)

		if bI != bJ {
			return bI < bJ
		}
		return sI < sJ
	})

	laterstTRC = trcFiles[len(trcFiles)-1] // Get the latest TRC version
	if len(trcFiles) > 1 {
		// XXX: This parameter is input to the scion-pki command, which internally
		// verifies that the penultimate TRC is in the grace period of the latest TRC.
		penultimateTRC = trcFiles[len(trcFiles)-2] // Get the penultimate TRC version
	}
	return laterstTRC, penultimateTRC, nil
}
