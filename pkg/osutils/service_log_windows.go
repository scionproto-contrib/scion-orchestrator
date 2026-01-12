package osutils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetJournalLogs fetches logs for a service from the Windows Event Log.
func GetJournalLogs(service string, lines int) (string, error) {
	// The log name where service events are typically stored.
	// This might need to be adjusted based on the specific service.
	logName := "System"

	// PowerShell command to get the last 'lines' events for the specified service.
	psCommand := fmt.Sprintf(
		"Get-WinEvent -LogName %s -ProviderName %s -MaxEvents %d | Format-List | Out-String",
		logName,
		service,
		lines,
	)

	// Execute the PowerShell command.
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", psCommand)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get logs for service %s: %v\n%s", service, err, stderr.String())
	}

	return out.String(), nil
}
