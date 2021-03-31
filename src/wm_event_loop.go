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
	/*
	address := host.wm_client_search(xunmap.window)
	if address == 0 { return }
	if host.client[address].app != xunmap.window { return }

	host.wm_client_withdraw(address, false)
	host.wm_host_update_client_focus()
	*/
	
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

	if host.wm_host_check_n_of_queued_event() >= 1 { return }

	address := host.wm_client_search(xconfig.window)
	if address == 0 { return }
	
	clt := host.client[address]
	if xconfig.window != clt.box.window { return }

	host.wm_client_adjust_app_for_box(address)
	host.wm_client_adjust_mask_for_box(address)
	
	host.wm_host_resize_surface(clt.box.surface, int(xconfig.width), int(xconfig.height))
	host.wm_host_draw_transparent(clt.box)
	
}


func (host *WmHost) wm_event_loop_configure_request(){
	var xcfgreq XConfigureRequestEvent
	xcfgreq = *(*XConfigureRequestEvent)(host.event.wm_event_get_pointer())
	if xcfgreq.window == XWindowID(XNone) { return }

	address := host.wm_client_search(xcfgreq.window)
	if address == 0 { return }

	clt := host.client[address]
	if xcfgreq.window != clt.app { return }

	border_width := host.config.client_drawable_range_border_width

	host.wm_client_configure(address,
							 int(xcfgreq.x)-border_width,
							 int(xcfgreq.y)-border_width,
							 int(xcfgreq.width)+border_width*2,
							 int(xcfgreq.height)+border_width*2)

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

		host.wm_host_update_grab_mode(host.grab_window, int(xbutton.x), int(xbutton.y))
	}

}

func (host *WmHost) wm_event_loop_button_release(){
	host.grab_window = XWindowID(XNone)
}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var motion XMotionEvent
	motion = *(*XMotionEvent)(host.event.wm_event_get_pointer())


	if host.grab_window != XWindowID(XNone) {

		

		xdiff := int(motion.x_root) - host.grab_root_x
		ydiff := int(motion.y_root) - host.grab_root_y

		if host.grab_mode_1 == 0 && host.grab_mode_2 == 0 {
			host.wm_host_move_window(host.grab_window,
									 host.grab_x + xdiff,
									 host.grab_y + ydiff)
		} else {

			if host.grab_mode_1 == WM_RESIZE_MODE_NONE { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_NONE { ydiff = 0 }
			if host.grab_mode_1 == WM_RESIZE_MODE_LEFT { xdiff = -xdiff }
			if host.grab_mode_2 == WM_RESIZE_MODE_TOP  { ydiff = -ydiff }

			host.wm_host_resize_window(host.grab_window,
									   host.grab_w + xdiff,
									   host.grab_h + ydiff,
			)

			if host.grab_mode_1 == WM_RESIZE_MODE_RIGHT  { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_BOTTOM { ydiff = 0 }

			host.wm_host_move_window(host.grab_window,
									 host.grab_x - xdiff,
									 host.grab_y - ydiff)

		}
	
	}
}
