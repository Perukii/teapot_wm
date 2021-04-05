package main

func (host *WmHost) wm_event_loop_map_notify(){
	var xmap XMapEvent
	xmap = *(*XMapEvent)(host.event.wm_event_get_pointer())
	if xmap.window == XWindowID(XNone) { return }
	if xmap.override_redirect == 1 { return }

	if host.wm_client_search(xmap.window) != 0 { return }

	address := host.wm_client_allocate_from_host()
	if address == 0 { return }

	host.wm_client_setup(address, xmap.window)

	host.wm_host_restack_clients()

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

	host.wm_client_withdraw(address)

}

func (host *WmHost) wm_event_loop_configure_notify(){
	/*
	var xconfig XConfigureEvent
	xconfig = *(*XConfigureEvent)(host.event.wm_event_get_pointer())
	if xconfig.window == XWindowID(XNone) { return }

	address := host.wm_client_search(xconfig.window)
	if address == 0 { return }

	clt := host.client[address]
	if xconfig.window != clt.mask.window { return }
	*/
}


func (host *WmHost) wm_event_loop_configure_request(){
	var xcfgreq XConfigureRequestEvent
	xcfgreq = *(*XConfigureRequestEvent)(host.event.wm_event_get_pointer())
	if xcfgreq.window == XWindowID(XNone) { return }

	var configure_normal = func(){
		host.wm_host_move_window(xcfgreq.window, int(xcfgreq.x), int(xcfgreq.y))
		host.wm_host_resize_window(xcfgreq.window, int(xcfgreq.width), int(xcfgreq.height))
	}

	address := host.wm_client_search(xcfgreq.window)
	if address == 0 {
		configure_normal()
		return
	}

	clt := host.client[address]
	if xcfgreq.window != clt.app {
		return
	}

	host.wm_client_configure(address,
		int(xcfgreq.x),
		int(xcfgreq.y),
		int(xcfgreq.width),
		int(xcfgreq.height),
		false,
	)
	

}


func (host *WmHost) wm_event_loop_enter_notify(){
	
}

func (host *WmHost) wm_event_loop_leave_notify(){

}

func (host *WmHost) wm_event_loop_property_notify(){
	var xproper XPropertyEvent
	xproper = *(*XPropertyEvent)(host.event.wm_event_get_pointer())
	if xproper.window == XWindowID(XNone) { return }
	address := host.wm_client_search(xproper.window)
	if address == 0 { return }
	clt := &host.client[address]
	if xproper.window != clt.app { return }

	if xproper.atom == host.wm_host_intern_atom("NET_WM_NAME") ||
	   xproper.atom == host.wm_host_intern_atom("WM_NAME") {
		clt.title = host.wm_host_get_window_title(clt.app)

		host.wm_host_draw_client(address)
	}


}