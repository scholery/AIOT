package listener

func InitListener() {
	go CheckListener()
	go PushListener()
	go DeviceStatusUpdate()
}
