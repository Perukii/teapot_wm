package main

//import "log"

func (host *WmHost) wm_event_loop_map_notify(){
	var xmap XMapEvent
	xmap = *(*XMapEvent)(host.event.wm_event_get_pointer())
	if xmap.window == XWindowID(XNone) { return }
	if xmap.override_redirect == 1 { return }

	if host.wm_client_search(xmap.window) != 0 { return }

	address := host.wm_client_allocate_from_host()
	if address == 0 { return }

	host.wm_client_setup(&host.client[address], xmap.window)

	host.wm_host_update_client_focus()

}


func (host *WmHost) wm_event_loop_unmap_notify(){
	var xunmap XUnmapEvent
	xunmap = *(*XUnmapEvent)(host.event.wm_event_get_pointer())
	if xunmap.window == XWindowID(XNone) { return }
	if xunmap.send_event == 0 { return }
	
	address := host.wm_client_search(xunmap.window)
	if address == 0 { return }
	if host.client[address].app != xunmap.window { return }

	host.wm_client_withdraw(address, false)
	host.wm_host_update_client_focus()
}

func (host *WmHost) wm_event_loop_destroy_notify(){
	var xdestroy XDestroyWindowEvent
	xdestroy = *(*XDestroyWindowEvent)(host.event.wm_event_get_pointer())
	if xdestroy.window == XWindowID(XNone) { return }

	address := host.wm_client_search(xdestroy.window)
	if address == 0 { return }
	if host.client[address].app != xdestroy.window { return }

	host.wm_client_withdraw(address, true)
	host.wm_host_update_client_focus()
	
}

func (host *WmHost) wm_event_loop_key_press(){
}

func (host *WmHost) wm_event_loop_button_press(){
	var xbutton XButtonEvent
	xbutton = *(*XButtonEvent)(host.event.wm_event_get_pointer())
	if xbutton.subwindow == XWindowID(XNone) { return }

	address := host.wm_client_search(xbutton.subwindow)

	if address == 0 { return }

	is_box := (xbutton.subwindow == host.client[address].box.window)
	is_mask := (xbutton.subwindow == host.client[address].mask.window)

	if is_mask{
		host.wm_host_set_focus_to_client(address)
	}

	if is_box{
		attr := host.wm_host_get_window_attributes(host.client[address].box.window)
		host.grab_window = host.client[address].box.window
		host.grab_button = int(xbutton.button)
		host.grab_root_x = int(xbutton.x_root)
		host.grab_root_y = int(xbutton.y_root)
		host.grab_x = int(attr.x)
		host.grab_y = int(attr.y)
		host.grab_w = int(attr.width)
		host.grab_h = int(attr.height)
	}


	//host.wm_host_set_focus_to_client(address)




}

func (host *WmHost) wm_event_loop_button_release(){
	host.grab_window = XWindowID(XNone)
}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var motion XMotionEvent
	motion = *(*XMotionEvent)(host.event.wm_event_get_pointer())
	
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
