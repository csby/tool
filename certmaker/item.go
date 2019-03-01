package main

import (
	"github.com/lxn/walk/declarative"
	"strings"
)

var (
	Kinds = []*Kind{
		{
			Name:        "server",
			DisplayName: "服务器",
		},
		{
			Name:        "client",
			DisplayName: "客户端",
		},
		{
			Name:        "user",
			DisplayName: "用户",
		},
	}
)

type Kind struct {
	Name        string
	DisplayName string
}

func kindDisplayName(name string) string {
	for _, item := range Kinds {
		if strings.ToLower(name) == item.Name {
			return item.DisplayName
		}
	}
	return name
}

func newLabel(text string, row, column int) declarative.Composite {
	return declarative.Composite{
		Row:    row,
		Column: column,
		Layout: declarative.VBox{
			Margins: declarative.Margins{
				Top: 8,
			},
		},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text:          text,
				TextAlignment: declarative.AlignHFarVNear,
			},
		},
	}
}
