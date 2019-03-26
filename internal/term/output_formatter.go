package term

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type (
	// TextFormatter formats logs into text
	TextFormatter struct {
	}

	metaMap map[string]string
)

var (
	metaFields = metaMap{
		"@Head": "",
	}
)

func isHeader(entry *log.Entry) bool {
	v, ok := entry.Data["@Head"]
	if ok {
		if b, ok := v.(bool); ok && b {
			return true
		}
	}

	return false
}

func dataAsString(entry *log.Entry) string {
	var sb strings.Builder
	for k, v := range entry.Data {
		if _, ok := metaFields[k]; ok {
			continue
		}
		if k == "error" {
			sb.WriteString(color.HiRedString(k))
		} else {
			sb.WriteString(color.New(color.Bold, color.FgCyan, color.Italic).Sprint(k))
		}
		sb.WriteString("=")
		sb.WriteString(color.YellowString(fmt.Sprintf("\"%v\" ", v)))
	}

	return sb.String()
}

func decoratedMessage(entry *log.Entry) string {
	return fmt.Sprintf("%s    %s", entry.Message, dataAsString(entry))
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *log.Entry) ([]byte, error) {

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefix := ""

	switch entry.Level {
	case log.PanicLevel, log.FatalLevel, log.ErrorLevel:
		prefix = color.RedString("\u2716 ")
	case log.InfoLevel:
		prefix = color.GreenString("\u2714 ")
	case log.WarnLevel:
		prefix = color.YellowString("\u2022 ")
	default:
		prefix = color.MagentaString("\u2022 ")
	}

	if isHeader(entry) {
		b.WriteString(color.HiWhiteString("\n    %s", strings.ToUpper(entry.Message)))
	} else {
		b.WriteString(fmt.Sprintf("    %s%s", prefix, decoratedMessage(entry)))
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
