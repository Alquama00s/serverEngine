package main

import (
	"github.com/Alquama00s/serverEngine/lib/DI"
	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
)

func main() {
	DI.InitialiseContextBuilder(".")
	autoconfigure.GetAppContextBuilder().BootStrap()
}
