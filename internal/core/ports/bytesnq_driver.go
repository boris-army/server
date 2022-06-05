package ports

type DriverTextNSDeliverFn = func(recipient string, data []byte) error

type DriverTextNSRenderFn = func(dst []byte) error

type DriverTextNQ interface {
	// Submit the text message for delivery. Returns whether the message
	// will be processed shortly or not (depends on a queue load).
	// Only reports render errors.
	Submit(renderFn DriverTextNSRenderFn, deliverFn DriverTextNSDeliverFn) (shortly bool, errRender error)
}
