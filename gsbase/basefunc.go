package gsbase

/*
func GetRaceState() bool {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Error.")
	}
	for _, s := range info.Settings {
		if s.Key == "-race" && s.Value == "true" {
			return true
		}
	}
	return false
}
*/

func GetRaceState() bool {
	return gsRaceState
}
