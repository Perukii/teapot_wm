package main

func (host *WmHost) wm_event_loop_key_press(){
	
	var xkey XKeyEvent
	xkey = *(*XKeyEvent)(host.event.wm_event_get_pointer())

	xkeytype := host.event.wm_event_get_type()
	keycode := int(xkey.keycode)

	keyflag := xkeytype == XKeyPress

	switch keycode{
	case 133:
		host.press_menu = keyflag
	case 111:
		host.press_up = keyflag
	case 113:
		host.press_left = keyflag
	case 116:
		host.press_down = keyflag
	case 114:
		host.press_right = keyflag
	}

	var address WmClientAddress
	address = len(host.client)-1

	if keycode == 113 && keyflag == true && host.press_menu == true {
		host.wm_client_harf_maximize(address, false)
	}
	if keycode == 114 && keyflag == true && host.press_menu == true {
		host.wm_client_harf_maximize(address, true)
	}
	if keycode == 111 && keyflag == true && host.press_menu == true {
		host.wm_client_toggle_maxmize(address)
	}

	host.wm_client_raise_mask(address)

}


func (host *WmHost) wm_event_loop_button_press(){
	var xbutton XButtonEvent
	xbutton = *(*XButtonEvent)(host.event.wm_event_get_pointer())
	if xbutton.subwindow == XWindowID(XNone) { return }

	address := host.wm_client_search(xbutton.subwindow)
	if address == 0 { return }

	clt := host.client[address]

	is_mask := (xbutton.subwindow == clt.mask.window)

	if is_mask{
		host.wm_host_set_focus_to_client(address)

		attr := host.wm_host_get_window_attributes(clt.mask.window)

		host.wm_host_update_button_mode(int(xbutton.x), int(xbutton.y),
				int(attr.x), int(attr.y), int(attr.width), int(attr.height),
		)

		host.wm_host_draw_client(address)

		if host.mask_button != WM_BUTTON_NONE { return }
	}

	attr := host.wm_host_get_window_attributes(clt.app)
	if int(xbutton.x) >= int(attr.x) &&
	   int(xbutton.x) <= int(attr.x) + int(attr.width) &&
	   int(xbutton.y) >= int(attr.y) &&
	   int(xbutton.y) <= int(attr.y) + int(attr.height) {
		return
	}

	host.grab_window = clt.app
	host.grab_button = int(xbutton.button)
	host.grab_root_x = int(xbutton.x_root)
	host.grab_root_y = int(xbutton.y_root)
	host.grab_x = int(attr.x)
	host.grab_y = int(attr.y)
	host.grab_w = int(attr.width)
	host.grab_h = int(attr.height)

	host.grab_comp_left = 0
	host.grab_comp_right = 0

	if clt.maximize_mode != WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_RIGHT &&
	   clt.maximize_mode != WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_LEFT { return }

	border_width := host.setting.client_border_overall_width
	shadow_width := host.setting.client_border_shadow_width
	bs_diff := border_width-shadow_width
	   
	for i := len(host.client)-1; i > 0; i-- {
		mclt := &host.client[i]
		if address == i { continue }
		if mclt.maximize_mode == clt.maximize_mode { continue }
		mattr := host.wm_host_get_window_attributes(mclt.mask.window)
		if mclt.maximize_mode == WM_CLIENT_MAXIMIZE_MODE_NEUTRAL_RIGHT {
			diff := int(mattr.x) - (host.grab_x + host.grab_w + bs_diff)
			if diff*diff <= bs_diff*bs_diff {
				host.grab_comp_right = i
				return
			}
		} else {
			diff := int(mattr.x) + int(mattr.width) - (host.grab_x - bs_diff)
			if diff*diff <= bs_diff*bs_diff {
				host.grab_comp_left = i
				return
			}
		}
		
		
	}
	
}

func (host *WmHost) wm_event_loop_button_release(){
	host.grab_window = XWindowID(XNone)

	var xbutton XButtonEvent
	xbutton = *(*XButtonEvent)(host.event.wm_event_get_pointer())
	if xbutton.subwindow == XWindowID(XNone) { return }

	address := host.wm_client_search(xbutton.subwindow)
	if address == 0 { return }

	clt := host.client[address]
	if xbutton.subwindow != clt.mask.window { return }

	if host.mask_button == WM_BUTTON_EXIT {
		host.wm_host_send_delete_event(clt.app)
	}

	if host.mask_button == WM_BUTTON_MAXIMIZE {
		host.wm_client_toggle_maxmize(address)
	} 

	host.mask_button = WM_BUTTON_NONE

	host.wm_host_draw_client(address)

}

func (host *WmHost) wm_event_loop_motion_notify(){
	
	var xmotion XMotionEvent
	xmotion = *(*XMotionEvent)(host.event.wm_event_get_pointer())

	grab_mode_top_when_maximized := func(address WmClientAddress) bool{
		return host.client[address].maximize_mode != WM_CLIENT_MAXIMIZE_MODE_NORMAL &&
				host.grab_mode_2 == WM_RESIZE_MODE_TOP
	}

	if host.grab_window != XWindowID(XNone) {

		address := host.wm_client_search(host.grab_window)
		clt := &host.client[address]


		xdiff := int(xmotion.x_root) - host.grab_root_x
		ydiff := int(xmotion.y_root) - host.grab_root_y

		expx := host.grab_x
		expy := host.grab_y
		expw := host.grab_w
		exph := host.grab_h

		if host.grab_mode_1 == WM_RESIZE_MODE_NONE && host.grab_mode_2 == WM_RESIZE_MODE_NONE {

			if clt.maximize_mode != WM_CLIENT_MAXIMIZE_MODE_NORMAL {
				host.wm_client_app_reverse_size(address)
				app_reverse_w_before := clt.app_reverse_w
				border_width := host.setting.client_border_overall_width
				attr := host.wm_host_get_window_attributes(clt.app)
				host.grab_x = int(xmotion.x)-app_reverse_w_before/2+border_width
				host.grab_y = int(xmotion.y)
				host.grab_w = int(attr.width)
				host.grab_h = int(attr.height)
				return
			}

			expx = host.grab_x + xdiff
			expy = host.grab_y + ydiff

		} else {

			if grab_mode_top_when_maximized(address) { return }

			if host.grab_mode_1 == WM_RESIZE_MODE_NONE { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_NONE { ydiff = 0 }
			if host.grab_mode_1 == WM_RESIZE_MODE_LEFT { xdiff = -xdiff }
			if host.grab_mode_2 == WM_RESIZE_MODE_TOP  { ydiff = -ydiff }

			expw = host.grab_w + xdiff
			exph = host.grab_h + ydiff

			if expw < clt.app_min_w {
				expw = clt.app_min_w
				xdiff = expw - host.grab_w
			}
			if clt.app_max_w != XNone && expw > clt.app_max_w { expx = clt.app_max_w }
			if exph < clt.app_min_h {
				exph = clt.app_min_h
				ydiff = exph - host.grab_h
			}
			if clt.app_max_h != XNone && exph > clt.app_max_h { exph = clt.app_max_h }

			if host.grab_mode_1 == WM_RESIZE_MODE_RIGHT  { xdiff = 0 }
			if host.grab_mode_2 == WM_RESIZE_MODE_BOTTOM { ydiff = 0 }

			expx = host.grab_x - xdiff
			expy = host.grab_y - ydiff

		}

		if host.grab_comp_right + host.grab_comp_left > 0 {
			rattr := host.wm_host_get_window_attributes(host.root_window)
			border_width := host.setting.client_border_overall_width
			shadow_width := host.setting.client_border_shadow_width
			bs_diff := border_width-shadow_width

			comp_w := int(rattr.width)-expw-bs_diff*4
	
			if host.grab_comp_right != 0 {
				cminw := host.client[host.grab_comp_right].app_min_w
				if comp_w < cminw {
					comp_w = cminw
					expw = int(rattr.width)-comp_w-bs_diff*4
					if expw < clt.app_min_w { return }
				}
				
				host.wm_client_configure(host.grab_comp_right, expx+expw+bs_diff*2, expy,
										 comp_w, exph, true)
			}

			if host.grab_comp_left != 0 {
				cminw := host.client[host.grab_comp_left].app_min_w
				if comp_w < cminw {
					expw = int(rattr.width)-cminw-bs_diff*4
					expx = expx + (cminw-comp_w)
					comp_w = cminw
					if expw < clt.app_min_w { return }
				}
				host.wm_client_configure(host.grab_comp_left, bs_diff, expy,
										 comp_w, exph, true)
			}
		}



		host.wm_client_configure(address, expx, expy, expw, exph, true)
		host.wm_host_update_cursor()


		return
	}



	address := host.wm_client_search(xmotion.subwindow)
	if address == 0 {
		host.wm_host_define_cursor(XCLeftPtr)
		return
	}

	clt := host.client[address]
	attr := host.wm_host_get_window_attributes(clt.mask.window)
	
	host.wm_host_update_button_mode(int(xmotion.x), int(xmotion.y),
								    int(attr.x), int(attr.y), int(attr.width), int(attr.height),
	)
	if host.mask_button != WM_BUTTON_NONE {
		host.grab_mode_1 = WM_RESIZE_MODE_NONE
		host.grab_mode_2 = WM_RESIZE_MODE_NONE
		host.wm_host_define_cursor(XCLeftPtr)
		return
	}

	host.wm_host_update_grab_mode(int(xmotion.x), int(xmotion.y),
								  int(attr.x), int(attr.y), int(attr.width), int(attr.height),
	)

	if grab_mode_top_when_maximized(address) {
		host.wm_host_define_cursor(XCLeftPtr)
		return
	}
	
	host.wm_host_update_cursor()


}
