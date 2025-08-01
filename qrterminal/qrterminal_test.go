package qrterminal

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	Generate("https://gitee.com/quant1x/gox/qrterminal", L, os.Stdout)
}

func TestGenerateWithConfig(t *testing.T) {
	config := Config{
		Level:     M,
		Writer:    os.Stdout,
		BlackChar: WHITE, // Inverted
		WhiteChar: BLACK,
		QuietZone: QUIET_ZONE,
	}
	GenerateWithConfig("https://gitee.com/quant1x/gox/qrterminal", config)
}

func TestGenerateHalfBlock(t *testing.T) {
	GenerateHalfBlock("https://gitee.com/quant1x/gox/qrterminal", L, os.Stdout)
}

func TestGenerateWithHalfBlockConfig(t *testing.T) {
	config := Config{
		Level:          M,
		Writer:         os.Stdout,
		HalfBlocks:     true,
		BlackChar:      BLACK_BLACK,
		WhiteBlackChar: WHITE_BLACK,
		WhiteChar:      WHITE_WHITE,
		BlackWhiteChar: BLACK_WHITE,
		QuietZone:      3,
	}
	GenerateWithConfig("https://gitee.com/quant1x/gox/qrterminal", config)
}
