package yoinkfuncs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
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

func RequestNewFileInfoForm(fvi FVInfo) FVInfo {

		form := huh.NewForm(
			// huh.NewGroup(huh.NewNote().
			// 	Title("File Version Info").
			// 	Description("Customize the File Version Info.\n\n").
			// 	Next(true).
			// 	NextLabel("Next"),
			// ),
	
			// huh.NewGroup(huh.NewNote().
			// 	Title("File Info").
			// 	Description("Input New File Info.\n\n"),
			// ),
	
			huh.NewGroup(
				huh.NewInput().
				Title("File Info").
				Description("Input New File Info.\n\n").
	
				Value(&fvi.CompanyName).
				Description("CompanyName").
				Placeholder(fvi.CompanyName),
	
				huh.NewInput().
				Value(&fvi.FileDescription).
				Description("FileDescription").
				Placeholder(fvi.FileDescription),
	
				huh.NewInput().
				Value(&fvi.FileVersion).
				Description("FileVersion").
				Placeholder(fvi.FileVersion),
	
				huh.NewInput().
				Value(&fvi.InternalName).
				Description("InternalName").
				Placeholder(fvi.InternalName),
	
				huh.NewInput().
				Value(&fvi.ProductName).
				Description("ProductName").
				Placeholder(fvi.ProductName),
	
				huh.NewInput().
				Value(&fvi.ProductVersion).
				Description("ProductVersion").
				Placeholder(fvi.ProductVersion),
	
				huh.NewInput().
				Value(&fvi.OriginalFilename).
				Description("OriginalFilename").
				Placeholder(fvi.OriginalFilename),
	
			),
	
			huh.NewGroup(
				huh.NewInput().
				Value(&fvi.Copyright).
				Description("Copyright").
				Placeholder(fvi.Copyright),
	
				huh.NewInput().
				Value(&fvi.Trademark).
				Description("Trademark").
				Placeholder(fvi.Trademark),
	
				huh.NewInput().
				Value(&fvi.Language).
				Description("Language").
				Placeholder(fvi.Language),
	
				huh.NewInput().
				Value(&fvi.Comments).
				Description("Comments").
				Placeholder(fvi.Comments),
			),
	
		).WithLayout(huh.LayoutStack)

		form.Run()

		return fvi

}
