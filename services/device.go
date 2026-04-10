package services

import (
	"regexp"
	"strings"

	"location-tracking-shortlink/models"
)

type DeviceService struct{}

func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

var (
	osPatterns = []struct {
		regex   *regexp.Regexp
		name    string
		version groupExtractor
	}{
		{regexp.MustCompile(`Windows NT 10\.0`), "Windows", extractVersion(`Windows NT 10\.0`)},
		{regexp.MustCompile(`Windows NT 6\.3`), "Windows", extractVersion(`Windows NT 6\.3`)},
		{regexp.MustCompile(`Windows NT 6\.2`), "Windows", extractVersion(`Windows NT 6\.2`)},
		{regexp.MustCompile(`Windows NT 6\.1`), "Windows", extractVersion(`Windows NT 6\.1`)},
		{regexp.MustCompile(`Mac OS X (\d+[._]\d+)`), "macOS", extractMacOSVersion},
		{regexp.MustCompile(`iPhone OS (\d+_\d+)`), "iOS", extractiOSVersion},
		{regexp.MustCompile(`iPad.*OS (\d+_\d+)`), "iPadOS", extractiOSVersion},
		{regexp.MustCompile(`Android (\d+(\.\d+)?)`), "Android", extractAndroidVersion},
		{regexp.MustCompile(`Linux`), "Linux", nil},
		{regexp.MustCompile(`CrOS`), "Chrome OS", nil},
	}

	browserPatterns = []struct {
		regex   *regexp.Regexp
		name    string
		version groupExtractor
	}{
		{regexp.MustCompile(`Edg/(\d+\.\d+\.\d+\.\d+)`), "Edge", nil},
		{regexp.MustCompile(`Chrome/(\d+\.\d+\.\d+\.\d+)`), "Chrome", nil},
		{regexp.MustCompile(`Firefox/(\d+\.\d+)`), "Firefox", nil},
		{regexp.MustCompile(`Safari/(\d+\.\d+)`), "Safari", extractSafariVersion},
		{regexp.MustCompile(`Opera|OPR`), "Opera", extractOperaVersion},
		{regexp.MustCompile(`MSIE (\d+\.\d+)`), "Internet Explorer", nil},
		{regexp.MustCompile(`Trident/.*rv:(\d+\.\d+)`), "Internet Explorer", nil},
	}

	mobileRegex = regexp.MustCompile(`Mobile|iPhone|iPad|iPod|Android.*Mobile|BlackBerry|IEMobile|Opera Mini|webOS`)
	tabletRegex = regexp.MustCompile(`iPad|Tablet|Kindle|Silk|PlayBook|Nexus 10|Xoom|SM-T`)
)

type groupExtractor func(match []string) string

func extractVersion(pattern string) groupExtractor {
	re := regexp.MustCompile(pattern)
	return func(match []string) string {
		return re.FindString(match[0])
	}
}

func extractMacOSVersion(match []string) string {
	if len(match) > 1 {
		return strings.ReplaceAll(match[1], "_", ".")
	}
	return ""
}

func extractiOSVersion(match []string) string {
	if len(match) > 1 {
		return strings.ReplaceAll(match[1], "_", ".")
	}
	return ""
}

func extractAndroidVersion(match []string) string {
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractSafariVersion(match []string) string {
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractOperaVersion(match []string) string {
	if len(match) > 0 {
		return match[0]
	}
	return ""
}

func (s *DeviceService) Parse(userAgent string) *models.DeviceInfo {
	if userAgent == "" {
		return &models.DeviceInfo{
			OS:         "Unknown",
			Browser:    "Unknown",
			DeviceType: "pc",
		}
	}

	info := &models.DeviceInfo{
		DeviceType: s.detectDeviceType(userAgent),
	}

	info.OS, info.OSVersion = s.parseOS(userAgent)
	info.Browser, info.BrowserVer = s.parseBrowser(userAgent, info.OS)

	return info
}

func (s *DeviceService) parseOS(userAgent string) (string, string) {
	for _, pattern := range osPatterns {
		if pattern.regex.MatchString(userAgent) {
			name := pattern.name
			var version string
			if pattern.version != nil {
				version = pattern.version(pattern.regex.FindStringSubmatch(userAgent))
			}
			return name, version
		}
	}
	return "Unknown", ""
}

func (s *DeviceService) parseBrowser(userAgent string, os string) (string, string) {
	for _, pattern := range browserPatterns {
		if pattern.regex.MatchString(userAgent) {
			name := pattern.name
			var version string
			if pattern.version != nil {
				version = pattern.version(pattern.regex.FindStringSubmatch(userAgent))
			} else {
				matches := pattern.regex.FindStringSubmatch(userAgent)
				if len(matches) > 1 {
					version = matches[1]
				}
			}
			return name, version
		}
	}
	return "Unknown", ""
}

func (s *DeviceService) detectDeviceType(userAgent string) string {
	if tabletRegex.MatchString(userAgent) {
		return "tablet"
	}
	if mobileRegex.MatchString(userAgent) {
		return "mobile"
	}
	return "pc"
}
