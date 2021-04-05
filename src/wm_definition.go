package main

import "os"

type WmClientAddress = int

type WmClient struct{
	app XWindowID
	mask WmTransparent
	app_min_w int
	app_min_h int
	app_max_w int
	app_max_h int
	maximize_mode int
	app_reverse_x int
	app_reverse_y int
	app_reverse_w int
	app_reverse_h int
	app_latest_x int
	app_latest_y int
	app_latest_w int
	app_latest_h int

	resize_process_wait float64

	title string
}

type WmTransparent struct{
	window 	XWindowID
	surface *CairoSfc
	drawtype int
}

type WmSetting struct{
	client_border_overall_width int
	client_border_shadow_width int
	client_button_width int
	client_button_margin_width int
	client_text_margin_width int
	max_resize_process_wait float64
}

type WmJsonSetting struct{
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

	grab_comp_left  WmClientAddress
	grab_comp_right WmClientAddress
	
	mask_button int

	cursor      int

	press_menu bool
	press_up bool
	press_down bool
	press_left bool
	press_right bool

	log_file *os.File

	client []WmClient

	setting WmSetting
	json_setting WmJsonSetting

}

type WmGeometry struct {
	x, y, w, h int
}

type WmWindowRelation struct{
	window 		XWindowID
	root_window XWindowID
	parent 		XWindowID
	children 	[]XWindowID
}

const(
	WM_RESIZE_MODE_NONE = 0
	WM_RESIZE_MODE_TOP = 1
	WM_RESIZE_MODE_BOTTOM = 2
	WM_RESIZE_MODE_RIGHT = 3
	WM_RESIZE_MODE_LEFT = 4

	WM_BUTTON_NONE = 0
	WM_BUTTON_EXIT = 1
	WM_BUTTON_MINIMIZE = 2
	WM_BUTTON_MAXIMIZE = 3

	WM_CLIENT_MAXIMIZE_MODE_NORMAL  = 0
	WM_CLIENT_MAXIMIZE_MODE_REVERSE = 1
	WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_RIGHT = 2
	WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_LEFT = 3
)