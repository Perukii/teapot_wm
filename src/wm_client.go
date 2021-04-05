package main

func (host *WmHost) wm_client_allocate_from_host() WmClientAddress{
	var clt WmClient
	host.client = append(host.client, clt)
	return len(host.client)-1
}

func (host *WmHost) wm_client_search(window XWindowID) WmClientAddress{
	for i := 0; i < len(host.client); i++ {
		if host.client[i].app == window { return i }
		if host.client[i].mask.window == window { return i }
	}

	return 0
}

func (host *WmHost) wm_client_withdraw(address WmClientAddress){
	host.wm_host_remove_transparent(host.client[address].mask)
	host.client[address] = host.client[len(host.client)-1]
	host.client = host.client[:len(host.client)-1]
}

func (host *WmHost) wm_client_setup(address WmClientAddress, xapp XWindowID){

	clt := &host.client[address]
	clt.app = xapp

	attr := host.wm_host_get_window_attributes(clt.app)

	border_width := host.config.client_border_overall_width

	// ---Mask---
	host.wm_host_setup_transparent(&clt.mask, host.root_window,
								   int(attr.x)-border_width,
								   int(attr.y)-border_width,
								   int(attr.width)+border_width*2,
								   int(attr.height)+border_width*2)

	host.wm_host_select_input(clt.mask.window,
				XSubstructureNotifyMask | XKeyPressMask | XKeyReleaseMask)
	host.wm_host_select_input(clt.app,
				XPropertyChangeMask | XKeyPressMask | XKeyReleaseMask)
	host.wm_host_map_window(clt.mask.window)

	size_hints := host.wm_host_get_size_hints(clt.app)
	clt.app_min_w = int(size_hints.min_width)
	clt.app_min_h = int(size_hints.min_height)
	clt.app_max_w = int(size_hints.max_width)
	clt.app_max_h = int(size_hints.max_height)

	clt.title = host.wm_host_get_window_title(clt.app)

	host.wm_host_draw_client(address)

	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NORMAL
}

func (host *WmHost) wm_client_raise_mask(address WmClientAddress){
	clt := host.client[address]
	host.wm_host_raise_window(clt.app)
	host.wm_host_raise_window(clt.mask.window)
	host.wm_host_draw_client(address)
}

func (host *WmHost) wm_client_raise_app(address WmClientAddress){
	clt := host.client[address]
	host.wm_host_raise_window(clt.mask.window)
	host.wm_host_raise_window(clt.app)
	
}

func (host *WmHost) wm_client_get_mask_geometry_from_app(x int, y int, w int, h int) (int,int,int,int) {
	
	border_width := host.config.client_border_overall_width
	mask_x := x - border_width
	mask_y := y - border_width
	mask_w := w + border_width*2
	mask_h := h + border_width*2
	return mask_x, mask_y, mask_w, mask_h

	// example
	// mask_x, mask_y, mask_w, mask_h := host.wm_client_get_mask_geometry_from_app(x, y, w, h)
}


func (host *WmHost) wm_client_get_app_geometry_from_mask(x int, y int, w int, h int) (int,int,int,int) {
	
	border_width := host.config.client_border_overall_width
	app_x := x + border_width
	app_y := y + border_width
	app_w := w - border_width*2
	app_h := h - border_width*2
	return app_x, app_y, app_w, app_h

	// example
	// app_x, app_y, app_w, app_h := host.wm_client_get_mask_geometry_from_app(x, y, w, h)
}