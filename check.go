// Browsercheck implements a simple way to check whether a specific 
// browser and installed plugins are outdated.
// The check is performed only using the useragent including the installed
// plugins, an example is provided in the subdirectory example. 
package browsercheck

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

// Defines a single application
type Application struct {
	ReadableName         string `xml:"readable-name"` // Readable name which can be shown to the user
	StringRegex          string `xml:"regex"`         // Regex to match the application and the version
	Regex                *regexp.Regexp
	LastSecureVersion    string `xml:"last-secure-version"` // Defines the last secure version
	lastSecureVersionInt []int
	UpdateUrl            string `xml:"update-url"` // Link to a document stating how to install the new version
}

type applications struct {
	Applications []Application `xml:"application"`
}

// This XML defines the last secure version of an application
// TODO: Check whether it is possible to move this to an external file.
const definitions = `
<data>
	<application>
		<readable-name>Mac OS X 10.9</readable-name>
		<regex>Mac OS X (10)_(9)_(\d+)</regex>
		<last-secure-version>10.9.1</last-secure-version>
		<update-url>http://support.apple.com/kb/ht1338</update-url>
	</application>
	<application>
		<readable-name>Mac OS X</readable-name>
		<regex>Mac OS X (10)_[0-8]_(\d+)</regex>
		<last-secure-version>10.9</last-secure-version>
		<update-url>https://www.apple.com/osx/how-to-upgrade/</update-url>
	</application>
	<application>
		<readable-name>Windows</readable-name>
		<regex>(Windows (?:NT 5\.(\d+)))</regex>
		<last-secure-version>5.3</last-secure-version>
		<update-url>https://www.microsoft.com/en-us/windows/enterprise/endofsupport.aspx</update-url>
	</application>
	<application>
		<readable-name>Google Chrome</readable-name>
		<regex>(Chromium|Chrome)/(\d+)\.(\d+)\.(\d+)\.(\d+)</regex>
		<last-secure-version>32.0.1700.68</last-secure-version>
		<update-url>https://support.google.com/chrome/answer/95414</update-url>
	</application>
	<application>
		<readable-name>Adobe Flash</readable-name>
		<regex>Shockwave Flash (\d+)\.(\d+) r(\d+)</regex>
		<last-secure-version>11.9.900</last-secure-version>
		<update-url>https://get.adobe.com/de/flashplayer/</update-url>
	</application>
	<application>
		<readable-name>Microsoft Silverlight</readable-name>
		<regex>Silverlight Plug\-In(\d+)\.(\d+).(\d+)</regex>
		<last-secure-version>5.1.20913</last-secure-version>
		<update-url>https://www.microsoft.com/getsilverlight/Get-Started/Install/Default.aspx</update-url>
	</application>
	<application>
		<readable-name>Internet Explorer</readable-name>
		<regex>Trident\/(\d+).(\d+)</regex>
		<last-secure-version>5.0</last-secure-version>
		<update-url>http://windows.microsoft.com/en-us/internet-explorer/download-ie</update-url>
	</application>
</data>`

var apps applications

// Initializes the library
func init() {
	xml.Unmarshal([]byte(definitions), &apps)

	for i, app := range apps.Applications {
		// Parse the regexes
		apps.Applications[i].Regex = regexp.MustCompile(app.StringRegex)

		// Split the last secure version into single strings and convert them to []int
		// e.g. 5.8.3 => [5 8 3]
		lastSecureVersion := strings.Split(app.LastSecureVersion, ".")
		for versionI := range lastSecureVersion {
			version, _ := strconv.Atoi(lastSecureVersion[versionI])
			apps.Applications[i].lastSecureVersionInt = append(apps.Applications[i].lastSecureVersionInt, version)
		}
	}
}

// Checks for insecure applications in an UA, returns nil or an list of insecure applications
// The UA may contain installed plugins, e.g. Adobe Flash
func Check(ua string) []Application {

	// Initialize a slice of potential outdated applications
	var outdatedApplications []Application

	// Loop over every regex
	for _, app := range apps.Applications {
		// Check whether the application has matched the regex
		foundApp := app.Regex.FindStringSubmatch(ua)

		// An regex has been matched, we have to determine determine
		// whether the version is older or newer than the last secure
		// version
		if len(foundApp) > 0 {

			// Parse the sent sent field
			foundVersion := make([]int, len(app.lastSecureVersionInt))
			for i := range app.lastSecureVersionInt {
				foundVersion[i], _ = strconv.Atoi(foundApp[i+1])
			}

			// Read the last fields and compare the version
			for i := 0; i < len(app.lastSecureVersionInt); i++ {

				// Mark applications as secure if a release is newer
				var secureApp bool

				// Verify whether the version is newer
				if foundVersion[i] >= app.lastSecureVersionInt[i] {
					secureApp = true
				} else {
					secureApp = false
				}

				// Mark as outdated if the passed version is lower than the last secure version
				if !secureApp && i == len(app.lastSecureVersionInt)-1 {
					outdatedApplications = append(outdatedApplications, app)
					break
				}
			}
		}
	}

	return outdatedApplications
}
