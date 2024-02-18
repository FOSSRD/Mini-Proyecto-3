package main

import (
	"fmt"
	"os"
	exec "os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Disk struct {
	name        string
	blockDevice string
	size        string
}

type DiskList struct {
	disks    []Disk
	cursor   int
	selected int
}

func (l DiskList) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func initialModel() DiskList {
	disksRawCmd := exec.Command("lsblk", "-o", "MODEL,SERIAL,SIZE,STATE,PATH", "--nodeps")
	diskRawOutput, err := disksRawCmd.Output()
	if err != nil {
		panic("no sirve: " + err.Error())
	}

	diskRawString := string(diskRawOutput)
	diskRawList := strings.Split(diskRawString, "\n")
	modelPos := strings.Index(diskRawList[0], "SERIAL")
	diskRawList = diskRawList[1 : len(diskRawList)-1]
	disks := []Disk{}
	for _, diskLine := range diskRawList {
		diskLineSep := strings.Fields(string(diskLine[modelPos:]))
		for _, val := range diskLineSep {
			fmt.Println(val)
		}
		disk := Disk{
			name:        string(diskLine[:modelPos]),
			blockDevice: diskLineSep[3],
			size:        diskLineSep[1],
		}
		disks = append(disks, disk)
	}
	fmt.Println(diskRawList)
	return DiskList{
		disks:    disks,
		selected: -1,
	}
}

func (m DiskList) View() string {
	// The header
	s := "Select disk to \n\n"

	// Iterate over our choices
	for i, choice := range m.disks {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if i == m.selected {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s@%s (%s)\n", cursor, checked, choice.name, choice.blockDevice, choice.size)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (l DiskList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if l.cursor < len(l.disks)-1 {
				l.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			if l.selected == l.cursor {
				l.selected = -1
			} else {
				l.selected = l.cursor
			}
		}

	}

	return l, nil
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
