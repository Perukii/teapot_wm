package main

import "os"

type WmClientAddress = int

type WmClient struct{
	app XWindowID
	mask WmTransparent
}

type WmTransparent struct{
	window 	XWindowID
	surface *CairoSfc
	drawtype int
}

type WmConfig struct{
	client_drawable_range_border_width int
	client_grab_area_resize_width int
	client_button_width int
	client_button_margin_width int
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
	grab_mode_1 int
	grab_mode_2 int
	grab_mode_configured int

	cursor      int

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
	WM_DRAW_TYPE_MASK = 2

	WM_RESIZE_MODE_NONE = 0
	WM_RESIZE_MODE_TOP = 1
	WM_RESIZE_MODE_BOTTOM = 2
	WM_RESIZE_MODE_RIGHT = 3
	WM_RESIZE_MODE_LEFT = 4


)