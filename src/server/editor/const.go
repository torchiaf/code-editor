package editor

import "fmt"

const NAME_LABEL string = "app.kubernetes.io/name"
const INSTANCE_LABEL string = "app.kubernetes.io/instance"
const MATCH_LABEL string = "app.code-editor/path"

const (
	Enabled  string = "ENABLED"
	Updating string = "UPDATING"
	Disabled string = "DISABLED"
	Unknown  string = "UNKNOWN"
)

func matchLabel(s string) string {
	return fmt.Sprintf("%s=%s", MATCH_LABEL, s)
}
