package main


/*
#cgo pkg-config: x11 cairo
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <cairo/cairo-xlib.h>
#include <X11/cursorfont.h>
#include <stdlib.h>
#include "./c_wm_draw.h"
#include "./c_wm_x_access.h"
*/
import "C"

import "unsafe"

type (
	XEvent = C.XEvent
	XDisplay = C.Display
	XWindowID = C.Window
	XWindowAttributes = C.XWindowAttributes
	XKeyEvent = C.XKeyEvent
	XButtonEvent = C.XButtonEvent
	XMotionEvent = C.XMotionEvent
	XMapEvent = C.XMapEvent
	XUnmapEvent = C.XUnmapEvent
	XMapRequestEvent = C.XMapRequestEvent
	XDestroyWindowEvent = C.XDestroyWindowEvent
	XConfigureEvent = C.XConfigureEvent
	XConfigureRequestEvent = C.XConfigureRequestEvent
	XCrossingEvent = C.XCrossingEvent
	XSizeHints = C.XSizeHints
	XPropertyEvent = C.XPropertyEvent

	XAtom = C.Atom

	CairoSfc = C.cairo_surface_t
	CairoCtx = C.cairo_t

	CLong = C.long
)

const (
	XNone = C.None
	XKeyPress = int(C.KeyPress)
	XKeyRelease = int(C.KeyRelease)
	XButtonPress = int(C.ButtonPress)
	XButtonRelease = int(C.ButtonRelease)
	XMotionNotify = int(C.MotionNotify)
	XMapNotify = int(C.MapNotify)
	XUnmapNotify = int(C.UnmapNotify)
	XMapRequest = int(C.MapRequest)
	XDestroyNotify = int(C.DestroyNotify)
	XConfigureNotify = int(C.ConfigureNotify)
	XConfigureRequest = int(C.ConfigureRequest)
	XEnterNotify = int(C.EnterNotify)
	XLeaveNotify = int(C.LeaveNotify)
	XPropertyNotify = int(C.PropertyNotify)

	XSubstructureNotifyMask = CLong(C.SubstructureNotifyMask)
	XSubstructureRedirectMask = CLong(C.SubstructureRedirectMask)
	XButtonPressMask = CLong(C.ButtonPressMask)
	XButtonReleaseMask = CLong(C.ButtonReleaseMask)
	XPointerMotionMask = CLong(C.PointerMotionMask)
	XEnterWindowMask = CLong(C.EnterWindowMask)
	XLeaveWindowMask = CLong(C.LeaveWindowMask)
	XPropertyChangeMask = CLong(C.PropertyChangeMask)
	XKeyPressMask = CLong(C.KeyPressMask)
	XKeyReleaseMask = CLong(C.KeyReleaseMask)

	XCLeftPtr = int(C.XC_left_ptr)
	XCSideT = int(C.XC_top_side)
	XCSideB = int(C.XC_bottom_side)
	XCSideL = int(C.XC_left_side)
	XCSideR = int(C.XC_right_side)
	XCSideTL = int(C.XC_top_left_corner)
	XCSideTR = int(C.XC_top_right_corner)
	XCSideBL = int(C.XC_bottom_left_corner)
	XCSideBR = int(C.XC_bottom_right_corner)

)

func wm_x11_get_type_of_event(event *XEvent) int {
	return int(C.c_wm_x11_get_type_of_event(*event))
}

func wm_x11_query_tree(display *XDisplay, window XWindowID) WmWindowRelation{
	var target WmWindowRelation
	var children *XWindowID
	var nc C.uint
	C.XQueryTree(display, window, &target.root_window, &target.parent, &children, &nc)

	target.children = make([]XWindowID, int(nc), int(nc))

	for i := 0; i < int(nc); i++ {
		target.children[i] = C.c_wm_x11_query_window_from_array(children, C.int(i))
	}

	if children != nil {
		C.XFree(unsafe.Pointer(children))
	}

	return target
}
