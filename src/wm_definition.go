package main

import "os"

type WmClientAddress = int

type WmClient struct{
	box WmTransparent
	app XWindowID
}

type WmTransparent struct{
	window 	XWindowID
	surface *CairoSfc
	drawtype int
}

type WmConfig struct{
	client_drawable_range_border_width int
}

type WmHost struct{
	display 	*XDisplay
	root_window XWindowID
	event		XEvent

	grab_window XWindowID
	grab_root_x	int
	grab_root_y	int
	grab_x		int
	grab_y		int
	grab_button int
	grab_w		int
	grab_h		int

	log_file *os.File

	client []WmClient

	config WmConfig
}

type WmWindowRelation struct{
	window 		XWindowID
	root_window XWindowID
	parent 		XWindowID
	children 	[]XWindowID
}

const(
	WM_DRAW_TYPE_NONE = 0
	WM_DRAW_TYPE_BOX = 1
)