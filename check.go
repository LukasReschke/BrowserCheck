// Browsercheck implements a simple check whether an user has insecure applications running
package browsercheck

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

type Application struct {
	ReadableName      string `xml:"readable-name"`
	StringRegex       string `xml:"regex"`
	Regex             *regexp.Regexp
	LastSecureVersion string `xml:"last-secure-version"`
	UpdateUrl         string `xml:"update-url"`
}

type Applications struct {
	Applications []Application `xml:"application"`
}

var q Applications

// Definitions
// This database should contain the newest definitions of applications with security updates.
// Feature updates should not be included in this list. 
// TODO: Move to data.xml
const definitions = `
<data>
	<application>
		<readable-name>Mac OS X 10.9</readable-name>
		<regex>(Mac OS X) (\d+)[_.]9(?:[_.](\d+))?</regex>
		<last-secure-version>10.9.1</last-secure-version>
		<update-url>http://support.apple.com/kb/ht1338</update-url>
	</application>
	<application>
		<readable-name>Mac OS X</readable-name>
		<regex>(Mac OS X) (\d+)[_.][0-8](?:[_.](\d+))?</regex>
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

// Preprocess the regexes
func init() {
	xml.Unmarshal([]byte(definitions), &q)

	// Parse the regexes
	for i, app := range q.Applications {
		q.Applications[i].Regex = regexp.MustCompile(app.StringRegex)
	}
}

// Checks for insecure applications in an UA, returns nil or an list of insecure applications
func Check(ua string) []Application {
	var insecureApplications []Application
	// Loop over every regex
	for _, app := range q.Applications {
		foundApp := app.Regex.FindStringSubmatch(ua)
		if len(foundApp) > 0 {
			// Parse the last secure version
			lastSecureVersion := strings.Split(app.LastSecureVersion, ".")

			// Convert the user agent version string to an integer value
			foundVersion := make([]int, len(lastSecureVersion))
			for i := 1; i < len(lastSecureVersion)+1; i++ {
				foundVersion[len(foundVersion)-i], _ = strconv.Atoi(foundApp[len(foundApp)-i])
			}

			// Read the last fields and compare the version
			for i := 0; i < len(lastSecureVersion); i++ {
				lastSecureVersionInt, _ := strconv.Atoi(lastSecureVersion[i])
				// Mark applications as secure if a release is newer
				var secureApp bool
				if foundVersion[i] >= lastSecureVersionInt {
					secureApp = true
				} else {
					secureApp = false
				}
				if !secureApp && i == len(lastSecureVersion)-1 {
					insecureApplications = append(insecureApplications, app)
					break
				}
			}
		}
	}
	return insecureApplications
}
