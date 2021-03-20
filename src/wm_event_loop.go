package main

import "log"

func (host *WmHost) wm_event_loop_key_press(){
}

func (host *WmHost) wm_event_loop_button_press(){
	var button XButtonEvent
	button = *(*XButtonEvent)(host.event.wm_event_get_pointer())
	if button.subwindow == XWindowID(XNone) { return }
	attr := host.wm_host_get_window_attributes(button.subwindow)
	host.grab_window = button.subwindow
	host.grab_button = int(button.button)
	host.grab_root_x = int(button.x_root)
	host.grab_root_y = int(button.y_root)
	host.grab_x = int(attr.x)
	host.grab_y = int(attr.y)
	host.grab_w = int(attr.width)
	host.grab_h = int(attr.height)
}

func (host *WmHost) wm_event_loop_button_release(){
	log.Println("release")
	host.grab_window = XWindowID(XNone)
}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var motion XMotionEvent
	motion = *(*XMotionEvent)(host.event.wm_event_get_pointer())
	
	if motion.subwindow == XWindowID(XNone) { return }
	if host.grab_window == XWindowID(XNone) { return }

	xdiff := int(motion.x_root) - host.grab_root_x
	ydiff := int(motion.y_root) - host.grab_root_y

	if host.grab_button == 1 {
		host.wm_host_move_window(host.grab_window, host.grab_x + xdiff, host.grab_y + ydiff)
	}

	if host.grab_button == 3 {
		host.wm_host_resize_window(host.grab_window, host.grab_w + xdiff, host.grab_h + ydiff)
	}

}
