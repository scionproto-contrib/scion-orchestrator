package osutils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// GetJournalLogs fetches logs for a service between start and end positions
func GetJournalLogs(service string, lines int) (string, error) {
	const lastWindow = "10m"

	predicate := fmt.Sprintf(
		`process == "%[1]s" OR subsystem == "%[1]s" OR senderImagePath CONTAINS[c] "%[1]s"`,
		service,
	)

	cmd := exec.Command(
		"log", "show",
		"--style", "syslog",
		"--predicate", predicate,
		"--last", lastWindow,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return tailLastN(out.Bytes(), lines), nil
}

// tailLastN returns the last n lines of a byte slice
func tailLastN(b []byte, n int) string {
	if n <= 0 || len(b) == 0 {
		return ""
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	const maxLine = 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxLine)

	ring := make([]string, n)
	idx, count := 0, 0
	for sc.Scan() {
		ring[idx] = sc.Text()
		idx = (idx + 1) % n
		if count < n {
			count++
		}
	}

	var res bytes.Buffer
	start := (idx + (n - count)) % n
	for i := 0; i < count; i++ {
		pos := (start + i) % n
		res.WriteString(ring[pos])
		res.WriteByte('\n')
	}
	return res.String()
}
