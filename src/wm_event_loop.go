package main

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
}

func (host *WmHost) wm_event_loop_map_request(){
	var xmapreq XMapRequestEvent
	xmapreq = *(*XMapRequestEvent)(host.event.wm_event_get_pointer())
	if xmapreq.window == XWindowID(XNone) { return }

	host.wm_host_map_window(xmapreq.window)

}

func (host *WmHost) wm_event_loop_destroy_notify(){
	var xdestroy XDestroyWindowEvent
	xdestroy = *(*XDestroyWindowEvent)(host.event.wm_event_get_pointer())
	if xdestroy.window == XWindowID(XNone) { return }

	address := host.wm_client_search(xdestroy.window)
	if address == 0 { return }
	if xdestroy.window != host.client[address].app { return }

	host.wm_client_withdraw(address, true)
	host.wm_host_update_client_focus()
	
}

func (host *WmHost) wm_event_loop_configure_notify(){
	
	var xconfig XConfigureEvent
	xconfig = *(*XConfigureEvent)(host.event.wm_event_get_pointer())
	if xconfig.window == XWindowID(XNone) { return }

	//if host.wm_host_check_n_of_queued_event() >= 1 { return }

	address := host.wm_client_search(xconfig.window)
	if address == 0 { return }

	clt := host.client[address]
	if xconfig.window != clt.mask.window { return }
	host.wm_host_draw_transparent(clt.mask)
	
}


func (host *WmHost) wm_event_loop_configure_request(){
	var xcfgreq XConfigureRequestEvent
	xcfgreq = *(*XConfigureRequestEvent)(host.event.wm_event_get_pointer())
	if xcfgreq.window == XWindowID(XNone) { return }

	host.wm_host_move_window(xcfgreq.window, int(xcfgreq.x), int(xcfgreq.y))
	host.wm_host_resize_window(xcfgreq.window, int(xcfgreq.width), int(xcfgreq.height))

	address := host.wm_client_search(xcfgreq.window)
	if address == 0 {

		return
	}

	


	clt := host.client[address]
	if xcfgreq.window != clt.app { return }


	host.wm_client_configure(address,
							 int(xcfgreq.x),
							 int(xcfgreq.y),
							 int(xcfgreq.width),
							 int(xcfgreq.height))
	

}

func (host *WmHost) wm_event_loop_key_press(){
}

func (host *WmHost) wm_event_loop_button_press(){
	var xbutton XButtonEvent
	xbutton = *(*XButtonEvent)(host.event.wm_event_get_pointer())
	if xbutton.subwindow == XWindowID(XNone) { return }

	address := host.wm_client_search(xbutton.subwindow)

	if address == 0 { return }

	is_mask := (xbutton.subwindow == host.client[address].mask.window)

	if is_mask{
		host.wm_host_set_focus_to_client(address)
	}

	attr := host.wm_host_get_window_attributes(host.client[address].app)
	host.grab_window = host.client[address].app
	host.grab_button = int(xbutton.button)
	host.grab_root_x = int(xbutton.x_root)
	host.grab_root_y = int(xbutton.y_root)
	host.grab_x = int(attr.x)
	host.grab_y = int(attr.y)
	host.grab_w = int(attr.width)
	host.grab_h = int(attr.height)

	host.wm_host_update_grab_mode(host.client[address].mask.window, int(xbutton.x), int(xbutton.y))


}

func (host *WmHost) wm_event_loop_button_release(){
	host.grab_window = XWindowID(XNone)
}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var motion XMotionEvent
	motion = *(*XMotionEvent)(host.event.wm_event_get_pointer())

	if host.wm_host_check_n_of_queued_event() >= 1 { return }

	address := host.wm_client_search(host.grab_window)
	clt := host.client[address]

	if host.grab_window != XWindowID(XNone) {
		xdiff := int(motion.x_root) - host.grab_root_x
		ydiff := int(motion.y_root) - host.grab_root_y

		expx := host.grab_x
		expy := host.grab_y
		expw := host.grab_w
		exph := host.grab_h

		if host.grab_mode_1 == 0 && host.grab_mode_2 == 0 {
			expx = host.grab_x + xdiff
			expy = host.grab_y + ydiff
		} else {

			if host.grab_mode_1 == WM_RESIZE_MODE_NONE { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_NONE { ydiff = 0 }
			if host.grab_mode_1 == WM_RESIZE_MODE_LEFT { xdiff = -xdiff }
			if host.grab_mode_2 == WM_RESIZE_MODE_TOP  { ydiff = -ydiff }

			host.wm_host_resize_window(clt.app,
									   host.grab_w + xdiff,
									   host.grab_h + ydiff,
			)
			expw = host.grab_w + xdiff
			exph = host.grab_h + ydiff

			if host.grab_mode_1 == WM_RESIZE_MODE_RIGHT  { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_BOTTOM { ydiff = 0 }

			expx = host.grab_x - xdiff
			expy = host.grab_y - ydiff

		}

		host.wm_client_configure(address, expx, expy, expw, exph)
	
	}
}
