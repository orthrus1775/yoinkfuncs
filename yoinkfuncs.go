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

const PKGVERSION = 0.1
const WINICON = winres.RT_ICON

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

func GetSpecRawResTypeData(rs1 *winres.ResourceSet, typeID int) []byte {

	unsafeCheckResTypeIdx(rs1, typeID)
	//refactor to to err checks?
	return unsafeGetResData(rs1, typeID)

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

func PerformResPatch(rs2 winres.ResourceSet, inTarget string) {
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


func modupFileVersionData(ogfvi FVInfo) FVInfo {
	ogfvi.Copyright = "NananaBooBoo"
	ogfvi.CompanyName = "Luna Industries"
	ogfvi.FileDescription = "Snake Generator"
	ogfvi.ProductVersion = "434.13"
	ogfvi.InternalName = "IDontWannaHearIt.exe"
	return ogfvi
}

func GetRawVersionInfo(rawfi []byte) *version.Info {
	vi, err := version.FromBytes(rawfi)
	if err != nil {
		log.Fatalf("Failed to marshal out the File Information. %v", err)

	}

	return vi
}  

func GetVersionInfoAsJSON(vi interface{}) []byte {
	// Cast the interface{} back to version.Info
	// Since GetRawVersionInfo returns 'any', we need to type assert
	versionInfo, ok := vi.(version.Info)
	if !ok {
		log.Fatalf("Failed to convert interface{} to version.Info")
		return nil
	}
	
	// Marshal to JSON using the existing MarshalJSON method
	jsonData, err := versionInfo.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal version info to JSON: %v", err)
		return nil
	}
	
	return jsonData
}

func GetSrcFileVersionData(jdata []byte) FVInfo {

	// Reexpports
	type FixedInfo struct {
		FileVersion    string `json:"file_version"`
		ProductVersion string `json:"product_version"`
	}

	type LocaleInfo struct {
		CompanyName     string `json:"CompanyName"`
		FileDescription string `json:"FileDescription"`
		FileVersion     string `json:"FileVersion"`
		InternalName    string `json:"InternalName"`
		OriginalFilename string `json:"OriginalFilename"`
		ProductName     string `json:"ProductName"`
		ProductVersion  string `json:"ProductVersion"`
		LegalCopyright  string `json:"LegalCopyright"`
		LegalTrademark	string `json:"LegalTrademark"`
		Language		string `json:"Language"`
		Comments		string `json:"Comments"`
	}

	type Data struct {
		Fixed FixedInfo           `json:"fixed"`
		Info  map[string]LocaleInfo `json:"info"`
	}

	///////////////////

	var dstruct Data
	json.Unmarshal(jdata, &dstruct)

	enUS_info, ok := dstruct.Info["0409"]
	if !ok {
		log.Fatal("The PE doesnt have an en-US string table. Suggest using manual tool")	
	}

	var fvi FVInfo
	fvi.CompanyName =  enUS_info.CompanyName
	fvi.FileDescription = enUS_info.FileDescription
	fvi.FileVersion = enUS_info.FileVersion
	fvi.InternalName = enUS_info.InternalName
	fvi.ProductName = enUS_info.ProductName
	fvi.ProductVersion = enUS_info.ProductVersion
	fvi.OriginalFilename = enUS_info.OriginalFilename

	fvi.Comments = enUS_info.Comments
	fvi.Copyright = enUS_info.LegalCopyright
	fvi.Trademark = enUS_info.LegalTrademark

	return fvi
}


// Update the data with custom updates
func SetDstFileInfoData(vi *version.Info, fvi FVInfo)  {

	vi.Set(0, "CompanyName",     fvi.CompanyName           )
	vi.Set(0, "FileDescription", fvi.FileDescription    )
	vi.Set(0, "FileVersion",     fvi.FileVersion            )
	vi.Set(0, "InternalName",    fvi.InternalName          )
	vi.Set(0, "ProductName",     fvi.ProductName            )
	vi.Set(0, "ProductVersion",  fvi.ProductVersion      )
	vi.Set(0, "OriginalFilename",fvi.OriginalFilename  )
								
	vi.Set(0, "Comments",        fvi.Comments                 )
	vi.Set(0, "LegalCopyright",  fvi.Copyright           )
	vi.Set(0, "LegalTrademarks", fvi.Trademark           )
	vi.Timestamp = setTimeZoneData("UTC")
	vi.SetFileVersion(fvi.FileVersion)
	vi.SetProductVersion(fvi.ProductVersion)
	return
}

func setTimeZoneData(date string) time.Time {

	bits := make([]int, 6)
	ansidt := time.ANSIC
	fmt.Sprintf(ansidt, date)
	date = convSystemToLocalTZArray(date)

	endbits := 0x7CF
	shl := endbits << 1
	shb := shl >> 1
	endbits = (shb ^ 0x0) | 0x0
	bits = append(bits, endbits)

	midbits := 0xC         
	shlm := midbits << 2 
	shbm := shlm >> 2
	midbits = (shbm ^ 0x0) & 0xFF

	st := midbits 
	st1 := st << 2      
	st2 := st1 - (0xC + 0x9) 
	st3 := st2 + 0xF     
	dy := st3 | 0x0
	bits = append(bits, (((dy << 1) << 0 ) - 1))
	bits = append(bits, (23^9) + 29)
	bits = append(bits, 10*10)
	dt := time.Date(endbits, 100 >> 3, dy, bits[1], bits[1], bits[1], 0, time.UTC)

	return dt
}

func convSystemToLocalTZArray(entry string) string {
    base := make([]uint8, 35)
	
	for i, j := range entry {

		if i <= 35 {

			base[i] = uint8(j)
		}

	}

    r1 := uint8(123 ^ 45)
    r2 := uint8(255 & 0)

    base = append(base, uint8(0))
    base = append(base, uint8(85 ^ 10))
    
    temp := uint8(99 ^ 33)
    temp = temp * r1
	r3 := uint8(r2 << 3)
    cyc := uint8(temp & r3)

    base = append(base, uint8(83 & 9))
    base = append(base, uint8(84 ^ 0))

    sentinel := uint8(47 & 3)   
    sentinel = sentinel * cyc         

    base = append(base, uint8(17 ^ 0))
    base = append(base, uint8(67 & 0))
    base = append(base, uint8(110 >> 7))
    base = append(base, uint8(161 >> 12))
    base = append(base, uint8(9 << 4))
    base = append(base, uint8(0))

    mid := string([]byte{base[23], base[r2], base[sentinel]})
	pref := strings.ToUpper(string([]byte{base[36], base[38], base[40]}))
	suf := strings.ToLower(string(uint8(42 ^ 42)))
	fmt.Sprintf("%s%s%x", pref, mid, suf)

    return pref
}


func unsafeGetResData(rs *winres.ResourceSet, sel int) []byte {

	rstr := reflect.ValueOf(rs)
	rstr = reflect.Indirect(rstr)
	rf := rstr.Field(0)

	if rf.Kind() == reflect.Map {
		typesFieldPtr := unsafe.Pointer(rf.UnsafeAddr())
		typesFieldRFValue := reflect.NewAt(rf.Type(), typesFieldPtr).Elem()
		mapKeys := typesFieldRFValue.MapKeys()

		var typeEntryPtr reflect.Value
		for _, key := range mapKeys {
			actualKey := key.Interface().(winres.Identifier)
				if (actualKey == winres.ID(16)) {
					typeEntryPtr = rf.MapIndex(key)
					break
				}
			}
			
		if !typeEntryPtr.IsValid() {
			log.Fatalf("Invalid entry pointer")
		}
		
		if typeEntryPtr.IsValid() && typeEntryPtr.Kind() == reflect.Ptr {
			typeEntryVal := typeEntryPtr.Elem()

			resourcesField := typeEntryVal.FieldByName("resources")
			if resourcesField.IsValid() && resourcesField.Kind() == reflect.Map {

				for _, resourceKey := range resourcesField.MapKeys() {
					resourcePtr := resourcesField.MapIndex(resourceKey)
					if resourcePtr.IsValid() && resourcePtr.Kind() == reflect.Ptr {
						resourceVal := resourcePtr.Elem()

						dataField := resourceVal.FieldByName("data")
						if dataField.IsValid() && dataField.Kind() == reflect.Map {

							for _, dataKey := range dataField.MapKeys() {
								dataPtr := dataField.MapIndex(dataKey)
								if dataPtr.IsValid() && dataPtr.Kind() == reflect.Ptr {
									dataVal := dataPtr.Elem()
									dataBytes := dataVal.FieldByName("data")
									if dataBytes.IsValid() {
										return dataBytes.Bytes()
									}
								}
							}
						}
					}
				}
			}
		}	
	}
	log.Fatalf("Did not grab %d resource was not located in the binary.", sel)
	panic(1)
}


func unsafeCheckResTypeIdx(rs *winres.ResourceSet, sel int) {

	rstr := reflect.ValueOf(rs)
	rstr = reflect.Indirect(rstr)
	rf := rstr.Field(0)

	if rf.Kind() == reflect.Map {
		typesFieldPtr := unsafe.Pointer(rf.UnsafeAddr())
		typesFieldRFValue := reflect.NewAt(rf.Type(), typesFieldPtr).Elem()
		mapKeys := typesFieldRFValue.MapKeys()

		for _, key := range mapKeys {
			actualKey := key.Interface().(winres.Identifier)
				if (actualKey == winres.ID(0x10)) {
					print("Confirmed availability File Information resource data.\n")
					return
				}	
		}

		log.Fatalf("Did not locate %d resource was not located in the binary.", sel)
		panic(1)
	}

	panic(1)
}

func dbgFVIColorPrint(fvi FVInfo) {

	color.Green("FVI Info for CompanyName = %s",  fvi.CompanyName                    )
	color.Green("FVI Info for FileDescription = %s", fvi.FileDescription             )
	color.Green("FVI Info for FileVersion = %s", fvi.FileVersion                     )
	color.Green("FVI Info for InternalName = %s", fvi.InternalName                   )
	color.Green("FVI Info for ProductName = %s", fvi.ProductName                     )
	color.Green("FVI Info for ProductVersion = %s", fvi.ProductVersion               )
	color.Green("FVI Info for OriginalFilename = %s", fvi.OriginalFilename           )																					 
	color.Green("FVI Info for Comments = %s", fvi.Comments                           )
	color.Green("FVI Info for Copyright = %s", fvi.Copyright                    )
	color.Green("FVI Info for Trademark = %s", fvi.Trademark                    )

}

func dbgJSONPrettyPrint(vi version.Info) {

		jsondata, _ := vi.MarshalJSON()
		fmt.Printf("%s", jsondata)
		fmt.Printf("Hex JSON: %x\n\n", jsondata)
		var jb bytes.Buffer
		json.Indent(&jb, jsondata, "", "  ")
		fmt.Printf("%s \n", jb)


		// Alternatively could initiate return to be used later for multiline editing.
		// return jb

}

func dbgRawUnMarshPrint(vi version.Info) {

	var unmarshjson []byte
	err := vi.UnmarshalJSON(unmarshjson)
	if err != nil {
		log.Fatalf("Unmarshal out err. %v", err)
	}
	fmt.Printf("Unmarshed: %x\n", 	string(unmarshjson))

}

func unsafeDbgExplResType(rs *winres.ResourceSet) {

	rstr := reflect.ValueOf(rs)
	rstr = reflect.Indirect(rstr)
	rf := rstr.Field(0)

	if rf.Kind() == reflect.Map {
		
		for _, key := range rf.MapKeys() {

			typeEntryPtr := rf.MapIndex(key)
			if typeEntryPtr.IsValid() && typeEntryPtr.Kind() == reflect.Ptr {
				typeEntryVal := typeEntryPtr.Elem()

				resourcesField := typeEntryVal.FieldByName("resources")
				if resourcesField.IsValid() && resourcesField.Kind() == reflect.Map {
					// Iterate the map specfically
					for _, resourceKey := range resourcesField.MapKeys() {
						resourcePtr := resourcesField.MapIndex(resourceKey)
						if resourcePtr.IsValid() && resourcePtr.Kind() == reflect.Ptr {
							resourceVal := resourcePtr.Elem()

							dataField := resourceVal.FieldByName("data")
							if dataField.IsValid() && dataField.Kind() == reflect.Map {

								for _, dataKey := range dataField.MapKeys() {
									dataPtr := dataField.MapIndex(dataKey)
									if dataPtr.IsValid() && dataPtr.Kind() == reflect.Ptr {
										dataVal := dataPtr.Elem()
										dataBytes := dataVal.FieldByName("data")
										if dataBytes.IsValid() {
											fmt.Printf("Data for %v: %s\n", dataKey, string(dataBytes.Bytes()))
											fmt.Printf("Byte Repr on %v: %x\n", dataKey, dataBytes.Bytes())
										}
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		fmt.Println("Expected 'types' mapping but was not found.")
		//crashable??
	}

	fmt.Println("Explorer End")
}
