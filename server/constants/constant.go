package constants

var Constants = Constant{
	RootPath:      "/var/lib/mist",
	LogPath:       "/var/lib/mist/logs",
	AvatarDirPath: "/var/lib/mist/uploads/avatar",
}

type Constant struct {
	RootPath      string
	LogPath       string
	AvatarDirPath string
}
