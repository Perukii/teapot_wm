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

func (host *WmHost) wm_client_setup(clt *WmClient, xapp XWindowID){

	clt.app = xapp

	attr := host.wm_host_get_window_attributes(clt.app)

	border_width := host.config.client_drawable_range_border_width

	// ---Mask---
	clt.mask.drawtype = WM_DRAW_TYPE_BOX
	host.wm_host_setup_transparent(&clt.mask, host.root_window,
								   int(attr.x)-border_width,
								   int(attr.y)-border_width,
								   int(attr.width)+border_width*2,
								   int(attr.height)+border_width*2)

	host.wm_host_select_input(clt.mask.window,
				XSubstructureNotifyMask)
	host.wm_host_map_window(clt.mask.window)

	host.wm_host_draw_transparent(clt.mask)

	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NORMAL
}

func (host *WmHost) wm_client_raise_mask(address WmClientAddress){
	clt := host.client[address]
	host.wm_host_raise_window(clt.app)
	host.wm_host_raise_window(clt.mask.window)
	host.wm_host_draw_transparent(clt.mask)
}

func (host *WmHost) wm_client_raise_app(address WmClientAddress){
	clt := host.client[address]
	host.wm_host_raise_window(clt.mask.window)
	host.wm_host_raise_window(clt.app)
}

func (host *WmHost) wm_client_configure(address WmClientAddress, x int, y int, w int, h int){
	clt := &host.client[address]

	border_width := host.config.client_drawable_range_border_width

	host.wm_host_move_window(clt.app, x, y)
	host.wm_host_resize_window(clt.app, w, h)

	host.wm_host_move_window  (clt.mask.window, x-border_width, y-border_width)
	host.wm_host_resize_transparent(clt.mask, w+border_width*2, h+border_width*2)
	host.wm_host_draw_transparent(clt.mask)

	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NORMAL
}

func (host *WmHost) wm_client_maximize(address WmClientAddress){
	
	clt := &host.client[address]

	{
		attr := host.wm_host_get_window_attributes(clt.mask.window)
		clt.reverse_x = int(attr.x)
		clt.reverse_y = int(attr.y)
		clt.reverse_w = int(attr.width)
		clt.reverse_h = int(attr.height)
	}
	
	{
		border_width := host.config.client_drawable_range_border_width
		shadow_width := host.config.client_grab_area_resize_width
	
		attr := host.wm_host_get_window_attributes(host.root_window)
	
		host.wm_client_configure(address,
							int(attr.x)+shadow_width,
							int(attr.y)+border_width,
							int(attr.width)-shadow_width*2,
							int(attr.height)-border_width-shadow_width)
	}

	
	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_REVERSE
}

func (host *WmHost) wm_client_reverse_size(address WmClientAddress){

	clt := &host.client[address]

	border_width := host.config.client_drawable_range_border_width

	host.wm_client_configure(address,
		clt.reverse_x+border_width,
		clt.reverse_y+border_width,
		clt.reverse_w-border_width*2,
		clt.reverse_h-border_width*2)
	
	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NORMAL
}