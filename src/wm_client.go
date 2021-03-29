package main

type WmClient struct{
	box_id XWindowID
	app_id XWindowID
}

type WmTransparent struct{
	window 	XWindowID
	surface *CairoSfc
	ctx 	*CairoCtx
}

func wm_client_new(host *WmHost){

}