package version

const (
	Major = 0
	Minor = 1
	Patch = 6

	AppName = "termfolio"
	AppDesc = "SSH-based interactive portfolio application served over SSH, built with Go, Wish, and Bubble Tea."
)

func Version() string {
	return AppName + " v" + VersionString()
}

func VersionString() string {
	return formatVersion(Major, Minor, Patch)
}

func formatVersion(major, minor, patch int) string {
	return toString(major) + "." + toString(minor) + "." + toString(patch)
}

func toString(i int) string {
	if i < 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune(i%10+48)) + s
		i /= 10
	}
	if s == "" {
		return "0"
	}
	return s
}

func VersionInfo() string {
	return Version() + " - " + AppDesc
}
