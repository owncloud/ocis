package useragent

// IsWindows shorthand function to check if OS == Windows
func (ua UserAgent) IsWindows() bool {
	return ua.OS == Windows
}

// IsAndroid shorthand function to check if OS == Android
func (ua UserAgent) IsAndroid() bool {
	return ua.OS == Android
}

// IsMacOS shorthand function to check if OS == MacOS
func (ua UserAgent) IsMacOS() bool {
	return ua.OS == MacOS
}

// IsIOS shorthand function to check if OS == IOS
func (ua UserAgent) IsIOS() bool {
	return ua.OS == IOS
}

// IsLinux shorthand function to check if OS == Linux
func (ua UserAgent) IsLinux() bool {
	return ua.OS == Linux
}

// IsChromeOS shorthand function to check if OS == CrOS
func (ua UserAgent) IsChromeOS() bool {
	return ua.OS == ChromeOS || ua.OS == "CrOS"
}

// IsBlackberryOS shorthand function to check if OS == BlackBerry
func (ua UserAgent) IsBlackberryOS() bool {
	return ua.OS == BlackBerry
}

// IsOpera shorthand function to check if Name == Opera
func (ua UserAgent) IsOpera() bool {
	return ua.Name == Opera
}

// IsOperaMini shorthand function to check if Name == Opera Mini
func (ua UserAgent) IsOperaMini() bool {
	return ua.Name == OperaMini
}

// IsChrome shorthand function to check if Name == Chrome
func (ua UserAgent) IsChrome() bool {
	return ua.Name == Chrome
}

// IsFirefox shorthand function to check if Name == Firefox
func (ua UserAgent) IsFirefox() bool {
	return ua.Name == Firefox
}

// IsInternetExplorer shorthand function to check if Name == Internet Explorer
func (ua UserAgent) IsInternetExplorer() bool {
	return ua.Name == InternetExplorer
}

// IsSafari shorthand function to check if Name == Safari
func (ua UserAgent) IsSafari() bool {
	return ua.Name == Safari
}

// IsEdge shorthand function to check if Name == Edge
func (ua UserAgent) IsEdge() bool {
	return ua.Name == Edge
}

// IsBlackBerry shorthand function to check if Name == BlackBerry
func (ua UserAgent) IsBlackBerry() bool {
	return ua.Name == BlackBerry
}

// IsGooglebot shorthand function to check if Name == Googlebot
func (ua UserAgent) IsGooglebot() bool {
	return ua.Name == Googlebot
}

// IsTwitterbot shorthand function to check if Name == Twitterbot
func (ua UserAgent) IsTwitterbot() bool {
	return ua.Name == Twitterbot
}

// IsFacebookbot shorthand function to check if Name == FacebookExternalHit
func (ua UserAgent) IsFacebookbot() bool {
	return ua.Name == FacebookExternalHit
}

// IsYandexbot shorthand function to check if Name == YandexBot
func (ua UserAgent) IsYandexbot() bool {
	return ua.Name == YandexBot
}

// IsUnknown returns true if the package can't determine the user agent reliably.
// Fields like Name, OS, etc. might still have values.
func (ua UserAgent) IsUnknown() bool {
	return !ua.Mobile && !ua.Tablet && !ua.Desktop && !ua.Bot
}
