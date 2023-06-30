package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func CheckLinkName(linkName string) bool {
	EXCLUDE_LINK_NAME_PATTERNS := []string{
		"^nvme-eui.",
		"_\\d(-part\\d)*$",
		"^wwn-",
	}
	for _, pattern := range EXCLUDE_LINK_NAME_PATTERNS {
		// 编译正则表达式
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 判断字符串是否匹配正则表达式
		if re.MatchString(linkName) {
			return false
		}
	}
	return true
}

func Lsblk(args ...string) gjson.Result {
	args = append([]string{"--json"}, args...)
	cmd := exec.Command("lsblk", append([]string{"-b"}, args...)...)
	output, _ := cmd.Output()
	return gjson.ParseBytes(output)
}

type BlockDevice struct {
	Name string
	Type string
	Size int
}

func Smartctl(args ...string) gjson.Result {
	args = append([]string{"--json=c"}, args...)
	cmd := exec.Command("smartctl", args...)
	output, _ := cmd.Output()
	return gjson.ParseBytes(output)
}

type Device struct {
	Name      string
	Path      string
	LabelPath string
	Type      string
	Protocol  string
}

func GetDevices() []*Device {
	// Get device labels form lsblk
	blockDevices := []*BlockDevice{}
	Lsblk("-d").Get("blockdevices").ForEach(func(k, v gjson.Result) bool {
		blockDevice := &BlockDevice{}
		json.Unmarshal([]byte(v.Raw), blockDevice)
		blockDevices = append(blockDevices, blockDevice)
		return true
	})

	// Get all links under /dev/disk/by-id
	linkFiles, err := ioutil.ReadDir("/dev/disk/by-id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	devices := []*Device{}
	// Check name of each link fine
	for _, file := range linkFiles {
		// Get it's name and link path if it is a symbolic link
		if file.Mode()&os.ModeSymlink != 0 {
			linkName := file.Name()

			if CheckLinkName(linkName) {
				linkPath, err := filepath.EvalSymlinks("/dev/disk/by-id/" + linkName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//fmt.Println(linkName, linkPath)

				for _, blockDevice := range blockDevices {
					smartResult := Smartctl("-i", linkPath)

					// Check if it is included by lsblk
					if linkPath == "/dev/"+blockDevice.Name {
						// Check if device support S.M.A.R.T.
						exitStatus := smartResult.Get("smartctl.exit_status").Int()
						doSupportSmart := exitStatus == 0
						if doSupportSmart {
							device := &Device{}
							json.Unmarshal([]byte(smartResult.Get("device").Raw), device)
							device.Name = linkName
							device.Path = "/dev/disk/by-id/" + linkName
							device.LabelPath = linkPath

							devices = append(devices, device)
						} else {
							fmt.Println(linkName + " doesn't support S.M.A.R.T.")
						}
						break
					}
				}
			}
		}
	}
	//result:=Smartctl("--scan-open")
	//result.Get("devices").ForEach(func(k, v gjson.Result) bool {
	//	device := &Device{}
	//	json.Unmarshal([]byte(v.Raw), device)
	//	devices = append(devices, device)
	//	return true
	//})

	return devices
}

type Result struct {
	ModelName       string `json:"model_name"`
	SerialNumber    string `json:"serial_number"`
	FirmwareVersion string `json:"firmware_version"`
	Passed          bool
	Attributes      map[string]float64
}

func GetAll(dev *Device) *Result {
	j := Smartctl("--xall", dev.Path)

	r := &Result{}
	json.Unmarshal([]byte(j.Raw), r)

	r.Passed = j.Get("smart_status.passed").Bool()

	r.Attributes = map[string]float64{}
	r.Attributes["temperature"] = j.Get("temperature.current").Float()
	r.Attributes["power_on_hours"] = j.Get("power_on_time.hours").Float()

	switch dev.Type {
	case "nvme":
		j.Get("nvme_smart_health_information_log").ForEach(func(k, v gjson.Result) bool {
			if v.Type == gjson.JSON {
				return true
			}
			r.Attributes[k.Str] = v.Float()
			return true
		})
	case "sat":
		j.Get("ata_smart_attributes.table").ForEach(func(k, v gjson.Result) bool {
			name := strings.ReplaceAll(strings.ToLower(v.Get("name").Str), "-", "_")
			rawStr := v.Get("raw.string").Str
			value := float64(v.Get("value").Int())

			rawFloat, err := strconv.ParseFloat(rawStr, 64)
			if err != nil {
				idx := strings.IndexByte(rawStr, '(')
				if idx == -1 {
					return true
				}

				rawStr = strings.TrimSpace(rawStr[:idx])
				if rawFloat, err = strconv.ParseFloat(rawStr, 64); err != nil {
					return true
				}
			}

			r.Attributes[name] = rawFloat
			r.Attributes[name+"_value"] = value

			if name == "total_lbas_written" {
				r.Attributes["data_units_written"] = rawFloat
			}

			return true
		})
	}

	return r
}
