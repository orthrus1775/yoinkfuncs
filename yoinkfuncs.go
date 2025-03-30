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

