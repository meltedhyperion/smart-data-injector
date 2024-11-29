package util

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func PadStringTo(v string, n int) string {
	if len(v) >= n {
		return v[:n]
	}
	return fmt.Sprintf("%-*s", n, v)
}

func IsValidFileType(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".json" || ext == ".csv"
}

func GenerateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(originalName, ext)
	uniqueID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), GenerateUUID())
	return fmt.Sprintf("%s-%s%s", base, uniqueID, ext)
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
