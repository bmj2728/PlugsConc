package logger

// Color represents a string type used for ANSI color codes.
// Supported color codes: Black, BlackBackground, Red, RedBackground, Green, GreenBackground, Yellow, YellowBackground,
// Blue, BlueBackground, Magenta, MagentaBackground, Cyan, CyanBackground, White, WhiteBackground, Default,
// DefaultBackground, BrightBlack, BrightBlackBackground, BrightRed, BrightRedBackground, BrightGreen,
// BrightGreenBackground, BrightYellow, BrightYellowBackground, BrightBlue, BrightBlueBackground, BrightMagenta,
// BrightMagentaBackground, BrightCyan, BrightCyanBackground, and ResetColor.
type Color string

// ColorSetting represents a combination of text and background colors using ANSI color codes.
type ColorSetting struct {
	Foreground Color
	Background Color
}

func NewColorSetting(color Color) ColorSetting {
	return ColorSetting{
		Foreground: color,
		Background: DefaultBackground,
	}
}

func (c ColorSetting) WithBackground(background Color) ColorSetting {
	c.Background = background
	return c
}

func NewColorSettingWithBackground(color Color, background Color) ColorSetting {
	return ColorSetting{
		Foreground: color,
		Background: background,
	}
}

const (
	// Black defines the ANSI color code for black text.
	Black = Color("\033[1;30m")
	// BlackBackground defines the ANSI color code for a black background.
	BlackBackground = Color("\033[1;40m")
	// Red defines the ANSI color code for red text.
	Red = Color("\033[1;31m")
	// RedBackground defines the ANSI color code for a red background.
	RedBackground = Color("\033[1;41m")
	// Green defines the ANSI color code for green text.
	Green = Color("\033[1;32m")
	// GreenBackground defines the ANSI color code for a green background.
	GreenBackground = Color("\033[1;42m")
	// Yellow defines the ANSI color code for yellow text.
	Yellow = Color("\033[1;33m")
	// YellowBackground defines the ANSI color code for a yellow background.
	YellowBackground = Color("\033[1;43m")
	// Blue defines the ANSI color code for blue text.
	Blue = Color("\033[1;34m")
	// BlueBackground defines the ANSI color code for a blue background.
	BlueBackground = Color("\033[1;44m")
	// Magenta defines the ANSI color code for magenta text.
	Magenta = Color("\033[1;35m")
	// MagentaBackground defines the ANSI color code for a magenta background.
	MagentaBackground = Color("\033[1;45m")
	// Cyan defines the ANSI color code for cyan text.
	Cyan = Color("\033[1;36m")
	// CyanBackground defines the ANSI color code for a cyan background.
	CyanBackground = Color("\033[1;46m")
	// White defines the ANSI color code for white text.
	White = Color("\033[1;37m")
	// WhiteBackground defines the ANSI color code for a white background.
	WhiteBackground = Color("\033[1;47m")
	// Default defines the ANSI color code for default the text color.
	Default = Color("\033[1;39m")
	// DefaultBackground defines the ANSI color code for a default background.
	DefaultBackground = Color("\033[1;49m")
	// BrightBlack defines the ANSI color code for bright black text.
	BrightBlack = Color("\033[1;90m")
	// BrightBlackBackground defines the ANSI color code for a bright black background.
	BrightBlackBackground = Color("\033[1;100m")
	// BrightRed defines the ANSI color code for bright red text.
	BrightRed = Color("\033[1;91m")
	// BrightRedBackground defines the ANSI color code for a bright red background.
	BrightRedBackground = Color("\033[1;101m")
	// BrightGreen defines the ANSI color code for bright green text.
	BrightGreen = Color("\033[1;92m")
	// BrightGreenBackground defines the ANSI color code for a bright green background.
	BrightGreenBackground = Color("\033[1;102m")
	// BrightYellow defines the ANSI color code for bright yellow text.
	BrightYellow = Color("\033[1;93m")
	// BrightYellowBackground defines the ANSI color code for a bright yellow background.
	BrightYellowBackground = Color("\033[1;103m")
	// BrightBlue defines the ANSI color code for bright blue text.
	BrightBlue = Color("\033[1;94m")
	// BrightBlueBackground defines the ANSI color code for a bright blue background.
	BrightBlueBackground = Color("\033[1;104m")
	// BrightMagenta defines the ANSI color code for bright magenta text.
	BrightMagenta = Color("\033[1;95m")
	// BrightMagentaBackground defines the ANSI color code for a bright magenta background.
	BrightMagentaBackground = Color("\033[1;105m")
	// BrightCyan defines the ANSI color code for bright cyan text.
	BrightCyan = Color("\033[1;96m")
	// BrightCyanBackground defines the ANSI color code for a bright cyan background.
	BrightCyanBackground = Color("\033[1;106m")
	// ResetColor defines the ANSI escape code to reset all color attributes.
	ResetColor = Color("\033[0m")
)
