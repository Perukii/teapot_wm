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

	border_width := host.config.client_drawable_range_border_width

	// ---Mask---
	clt.mask.drawtype = WM_DRAW_TYPE_BOX
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

func (host *WmHost) wm_client_configure(address WmClientAddress, x int, y int, w int, h int){
	clt := &host.client[address]

	border_width := host.config.client_drawable_range_border_width
	host.wm_host_move_resize_window(clt.app, x, y, w, h)
	
	mask_x := x - border_width
	mask_y := y - border_width
	mask_w := w + border_width*2
	mask_h := h + border_width*2

	host.wm_host_move_resize_window(clt.mask.window, mask_x, mask_y, mask_w, mask_h)
	host.wm_host_resize_surface(clt.mask.surface, mask_w, mask_h)
	host.wm_host_draw_client(address)
}

func (host *WmHost) wm_client_harf_maximize(address WmClientAddress, is_right bool){
	if is_right{
		host.wm_client_set_maximize(address, 0, 1)
	} else {
		host.wm_client_set_maximize(address, 1, 0)
	}
}

func (host *WmHost) wm_client_maximize(address WmClientAddress){
	host.wm_client_set_maximize(address, 1, 1)
}

func (host *WmHost) wm_client_set_maximize(address WmClientAddress, left int, right int){
	
	clt := &host.client[address]

	if clt.maximize_mode == WM_CLIENT_MAXIMIZE_MODE_NORMAL{
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
		
		conf_w := (int(attr.width)-(border_width)*2)/(3-left-right)
		conf_x := int(attr.x)+
					(int(attr.width)/2+shadow_width)*right*(1-left)+
					shadow_width*left*(1-right)

		host.wm_client_configure(address,
							conf_x,
							int(attr.y)+border_width,
							conf_w,
							int(attr.height)-border_width-shadow_width)
	}
	if left == 1 && right == 1{
		clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_REVERSE
	} else {
		clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NEUTRAL
	}
	
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