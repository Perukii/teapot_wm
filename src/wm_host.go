package main

import "os"
import "log"

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
}

func (host *WmHost) wm_host_init_x11_element(){
	host.display = wm_x11_open_display()
	host.root_window = wm_x11_get_root(host.display)
}

func (host *WmHost) wm_host_init_log_file(){
	wm_debug_log_file_init(host.log_file)
}

func (host *WmHost) wm_host_close_log_file(){
	host.log_file.Close()
}

func (host *WmHost) wm_host_grab_button(window XWindowID){
	wm_x11_grab_button(host.display, window)
}

func (host *WmHost) wm_host_raise_window(window XWindowID){
	wm_x11_raise_window(host.display, window)
}

func (host *WmHost) wm_host_get_window_attributes(window XWindowID) XWindowAttributes{
	return wm_x11_get_window_attributes(host.display, window)
}

func (host *WmHost) wm_host_move_window(window XWindowID, x int, y int){
	wm_x11_move_window(host.display, window, x, y)
}

func (host *WmHost) wm_host_resize_window(window XWindowID, w int, h int){
	wm_x11_resize_window(host.display, window, w, h)
}

func (host *WmHost) wm_host_setup_transparent(transparent *WmTransparent,
														x int, y int, w int, h int,
														){
	wm_x11_create_transparent_window(host.display, host.root_window, x, y, w, h, transparent)
}

func (host *WmHost) wm_host_map_window(window XWindowID){
	wm_x11_map_window(host.display, window)
}

func (host *WmHost) wm_host_run(){

	var trans WmTransparent
	host.wm_host_setup_transparent(&trans, 100, 100, 500, 500)
	host.wm_host_map_window(trans.window)
	host.wm_host_raise_window(trans.window)
	wm_x11_draw_transparent(host.display, trans)

	for{
		host.event = wm_x11_peek_event(host.display)
		switch host.event.wm_event_get_type(){
		case XMapNotify:
			log.Println("map")
		case XKeyPress:
			host.wm_event_loop_key_press()
		case XButtonPress:
			host.wm_event_loop_button_press()
		case XButtonRelease:
			host.wm_event_loop_button_release()
		case XMotionNotify:
			host.wm_event_loop_motion_notify()
		}
	}
}


