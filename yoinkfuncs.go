package yoinkfuncs

import (
	"path/filepath"
)

// Required Entries
// Ref: https://learn.microsoft.com/en-us/windows/win32/menurc/versioninfo-resource
	type FVInfo struct {
		CompanyName		string
		FileDescription    string
		FileVersion        string
		InternalName		string
		ProductName        string
		ProductVersion     string
		OriginalFilename   string
	// Extended Entries
	// Not Required But Can Be Commonplace
		Copyright       string
		Trademark		string
		Language		string
		Comments		string
}

func UNU(x ...interface{}) {}

func addPrefixToFileName(fp string) string {
	dir := filepath.Dir(fp)
	baseName := filepath.Base(fp)
	newBaseName := "yoinked-" + baseName
	return filepath.Join(dir, newBaseName)
}


