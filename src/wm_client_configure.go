package main

func (host *WmHost) wm_client_configure(address WmClientAddress,
					x int, y int, w int, h int, from_motion bool){

	clt := &host.client[address]

	resize_needed := w != clt.app_latest_w || h != clt.app_latest_h

	if clt.config_wait > 0 && resize_needed && from_motion{
		clt.config_wait--
		return
	}
	clt.config_wait = host.config.max_config_wait

	border_width := host.config.client_drawable_range_border_width
	
	
	mask_x := x - border_width
	mask_y := y - border_width
	mask_w := w + border_width*2
	mask_h := h + border_width*2

	if resize_needed {
		host.wm_host_move_resize_window(clt.app, x, y, w, h)
		host.wm_host_move_resize_window(clt.mask.window, mask_x, mask_y, mask_w, mask_h)
		host.wm_host_resize_surface(clt.mask.surface, mask_w, mask_h)
		host.wm_host_draw_client(address)
	} else {
		host.wm_host_move_window(clt.app, x, y)
		host.wm_host_move_window(clt.mask.window, mask_x, mask_y)
	}


	host.wm_client_update_app_latest_geometry(address, x, y, w, h)
}

func (host *WmHost) wm_client_toggle_maxmize(address WmClientAddress){
	clt := &host.client[address]
	if clt.maximize_mode == WM_CLIENT_MAXIMIZE_MODE_REVERSE {
		host.wm_client_app_reverse_size(address)
	} else {
		host.wm_client_maximize(address)
	}
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
		//host.wm_client_update_app_latest_geometry_from_attributes(address)
		clt.app_reverse_x = clt.app_latest_x
		clt.app_reverse_y = clt.app_latest_y
		clt.app_reverse_w = clt.app_latest_w
		clt.app_reverse_h = clt.app_latest_h
	}
	
	{
		rattr := host.wm_host_get_window_attributes(host.root_window)
		border_width := host.config.client_drawable_range_border_width
		shadow_width := host.config.client_grab_area_resize_width
		
		bs_diff := border_width-shadow_width

		var conf_x, conf_w int
		if left != 1 || right != 1{
			conf_x = int(rattr.x)+
					(int(rattr.width)/2+bs_diff)*right*(1-left)+
					bs_diff*left*(1-right)
			conf_w = (int(rattr.width)-bs_diff*4)/2

		} else {
			conf_x = int(rattr.x)+bs_diff
			conf_w = int(rattr.width)-bs_diff*2
		}

		host.wm_client_configure(address,
			conf_x,
			int(rattr.y) + border_width,
			conf_w,
			int(rattr.height) - border_width - bs_diff,
			false,
		)
	}
	if left == 1 && right == 1{
		clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_REVERSE
	} else if left == 1 {
		clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_LEFT
	} else {
		clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_RIGHT
	}
	
}

func (host *WmHost) wm_client_app_reverse_size(address WmClientAddress){

	clt := &host.client[address]

	host.wm_client_configure(address,
		clt.app_reverse_x,
		clt.app_reverse_y,
		clt.app_reverse_w,
		clt.app_reverse_h,
		false,
	)
	
	clt.maximize_mode = WM_CLIENT_MAXIMIZE_MODE_NORMAL
}

func (host *WmHost) wm_client_update_app_latest_geometry_from_attributes(address WmClientAddress){
	
	clt := &host.client[address]
	attr := host.wm_host_get_window_attributes(clt.app)

	clt.app_latest_x = int(attr.x)
	clt.app_latest_y = int(attr.y)
	clt.app_latest_w = int(attr.width)
	clt.app_latest_h = int(attr.height)
}

func (host *WmHost) wm_client_update_app_latest_geometry(address WmClientAddress,
					x int, y int, w int, h int){
	
	clt := &host.client[address]

	clt.app_latest_x = x
	clt.app_latest_y = y
	clt.app_latest_w = w
	clt.app_latest_h = h
}