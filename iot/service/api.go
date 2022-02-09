package service

func InitListener() {
	go CheckListener()
	go PushListener()
	go DeviceStatusUpdate()
}
