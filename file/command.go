package file

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stack"
)

// Run is the interface to picking a file.
func (o Options) Run() error {
	if o.Path == "" {
		o.Path = "."
	}

	path, err := filepath.Abs(o.Path)
	if err != nil {
		return err
	}

	m := model{
		path:            path,
		cursor:          o.Cursor,
		selected:        0,
<<<<<<< HEAD
		showHidden:      o.All,
=======
>>>>>>> next
		autoHeight:      o.Height == 0,
		height:          o.Height,
		max:             0,
		min:             0,
		selectedStack:   stack.NewStack(),
		minStack:        stack.NewStack(),
		maxStack:        stack.NewStack(),
		cursorStyle:     o.CursorStyle.ToLipgloss(),
<<<<<<< HEAD
		symlinkStyle:    o.SymlinkStyle.ToLipgloss(),
=======
>>>>>>> next
		directoryStyle:  o.DirectoryStyle.ToLipgloss(),
		fileStyle:       o.FileStyle.ToLipgloss(),
		permissionStyle: o.PermissionsStyle.ToLipgloss(),
		selectedStyle:   o.SelectedStyle.ToLipgloss(),
		fileSizeStyle:   o.FileSizeStyle.ToLipgloss(),
	}

	tm, err := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).StartReturningModel()
	if err != nil {
		return err
	}

	m = tm.(model)

	if m.path == "" {
		os.Exit(1)
	}

	fmt.Println(m.path)

	return nil
}
