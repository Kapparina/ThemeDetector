//go:build windows

package app

type Appearance string

const (
	Dark              Appearance = "Dark"
	Light                        = "Light"
	DarkHighContrast             = "DarkHighContrast"
	LightHighContrast            = "LightHighContrast"
)

func GetAppearance(index uint64) Appearance {
	return [...]Appearance{Dark, Light, DarkHighContrast, LightHighContrast}[index]
}
