package contexti

import "net/url"

type DeviceInfo struct {
	//设备
	Device     string `json:"device" gorm:"size:255"`
	Os         string `json:"os" gorm:"size:255"`
	AppCode    string `json:"appCode" gorm:"size:255"`
	AppVersion string `json:"appVersion" gorm:"size:255"`
	IP         string `json:"ip" gorm:"size:255"`
	Lng        string `json:"lng" gorm:"type:numeric(10,6)"`
	Lat        string `json:"lat" gorm:"type:numeric(10,6)"`
	Area       string `json:"area" gorm:"size:255"`
	UserAgent  string `json:"userAgent" gorm:"size:255"`
}

func Device(infoHeader, area, localHeader, userAgent, ip string) *DeviceInfo {
	unknow := true
	var info DeviceInfo
	//Device-Info:device,osInfo,appCode,appVersion
	if infoHeader != "" {
		unknow = false
		var n, m int
		for i, c := range infoHeader {
			if c == '-' {
				switch n {
				case 0:
					info.Device = infoHeader[m:i]
				case 1:
					info.Os = infoHeader[m:i]
				case 2:
					info.AppCode = infoHeader[m:i]
				case 3:
					info.AppVersion = infoHeader[m:i]
				}
				m = i + 1
				n++
			}
		}
	}
	// area:xxx
	// location:1.23456,2.123456
	if area != "" {
		unknow = false
		info.Area, _ = url.PathUnescape(area)
	}
	if localHeader != "" {
		unknow = false
		var n, m int
		for i, c := range localHeader {
			if c == '-' {
				switch n {
				case 0:
					info.Lng = localHeader[m:i]
				case 1:
					info.Lat = localHeader[m:i]
				}
				m = i + 1
				n++
			}
		}

	}

	if userAgent != "" {
		unknow = false
		info.UserAgent = userAgent
	}
	if ip != "" {
		unknow = false
		info.IP = ip
	}
	if unknow {
		return nil
	}
	return &info
}
