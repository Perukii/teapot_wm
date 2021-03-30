package main

func (host *WmHost) wm_event_loop_map_notify(){
	var xmap XMapEvent
	xmap = *(*XMapEvent)(host.event.wm_event_get_pointer())
	if xmap.window == XWindowID(XNone) { return }
	if xmap.override_redirect == 1 { return }

	address := host.wm_client_allocate_from_host()
	if address == 0 { return }

	host.wm_client_setup(&host.client[address], xmap.window)

}


func (host *WmHost) wm_event_loop_unmap_notify(){
	var xunmap XUnmapEvent
	xunmap = *(*XUnmapEvent)(host.event.wm_event_get_pointer())
	if xunmap.window == XWindowID(XNone) { return }
	if xunmap.send_event == 0 { return }

	address := host.wm_client_search(xunmap.window)
	if address == 0 { return }

	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.box.window)
	host.wm_host_reparent_window(clt.app, host.wm_host_query_parent(clt.box.window),
								int(attr.x), int(attr.y))

	host.wm_client_withdraw(address)

}

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
	host.grab_window = XWindowID(XNone)
}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var motion XMotionEvent
	motion = *(*XMotionEvent)(host.event.wm_event_get_pointer())
	
	//if motion.subwindow == XWindowID(XNone) { return }
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
