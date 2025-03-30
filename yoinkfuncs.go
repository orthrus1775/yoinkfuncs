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

func LoadAllResourcesFromPath(path string) *winres.ResourceSet {

	fr, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open PE file: %v", err)
	}
	defer fr.Close()


	res, err := winres.LoadFromEXE(fr)
	if err != nil {
		log.Fatalf("Failed to load EXE resources: %v", err)
	}

	return res

}

func LoadAnIconResourceFromPath(path string) *winres.ResourceSet {

	fr, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open PE file: %v", err)
	}
	defer fr.Close()


	res, err := winres.LoadFromEXESingleType(fr, winres.RT_ICON)
	if err != nil {
		log.Fatalf("Failed to load ICO EXE resources: %v", err)
	}

	return res

}

// Purely a stop-gap function. Will need to port in my PE parser in future update. As of now will only work with "our current" appset.
func SearchForCommonICOGroups(res *winres.ResourceSet) *winres.Icon {

	var icos_names []string = []string{"IDI_APPLICATION", "IDR_MAINFRAME", "IDR_MAINFRAME_2", "IDR_MAINFRAME_3", "IDR_MAINFRAME_4"}
	var icos_numbs []int  = []int{
		1,
		2,3,4,
		5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 107, 108, 
		184,
		201, 
		1300, 1301, 1302, 1303, 1304, 1305, 1306, 1307,1308, 1309, 1310, 1311, 1312, 1313}
	fmt.Println("Performing a search for known common icon groups.")

	
	for _, e := range icos_numbs {
		ico, err := res.GetIcon(winres.ID(e))
		if err != nil {
			//fmt.Printf("Fail match lookup of %d.\n", e)
			continue
		} else if err == nil {
			fmt.Printf("ID MATCHED on lookup of %d.\n", e)
			return ico
		}
	}
		
	fmt.Println("No Number Matches. Performing Name Lookups")

	for _, e := range icos_names {
		ico, err := res.GetIcon(winres.Name(e))
		if err != nil {
			//fmt.Printf("Fail match lookup of %s.\n", e)
			continue
		} else if err == nil {
			fmt.Printf("NAME MATCHED on lookup of %s.\n", e)
			return ico
		}
	}

	panic("Failed to find matchable number or name. Consider manual specification.")
}

func performResPatch(rs2 winres.ResourceSet, inTarget string) {
	// How this works is the input target file is what we're modifying to be the output
	// target= rewriteblank.exe  src= C:\\Users\\crt\\Desktop\\pe-ops\\winres\\hwblank.exe
	purfile, err := os.Open(inTarget)
	if err != nil {
		log.Fatalf("Failed to open input PE file: %v", err)
	}

	outTarget := addPrefixToFileName(inTarget)

	rwfile, err := os.Create(outTarget)
	if err != nil {
		log.Fatalf("Could not open dst location: %v", err)
	}
	defer purfile.Close()
	defer rwfile.Close()

	err = rs2.WriteToEXE(rwfile, purfile, winres.WithAuthenticode(2))
	if err != nil {
		log.Fatalf("Failed to write new EXE: %v", err)
	}

	fmt.Println("Successfully wrote file: ", outTarget)

}
