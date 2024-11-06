package main

func main() {
	configs := []DisplayConfig{
		NewDisplayConfig("eDP-1", Position{x: 0, y: 0}, Resolution{x: 1280, y: 720}, 0.9, true),
		NewDisplayConfig("HDMI-1", Position{x: 1280, y: -1000}, Resolution{x: 1280, y: 720}, 1, false),
		NewDisplayConfig("HDMI-2", Position{x: 1280, y: 0}, Resolution{x: 3840, y: 2160}, 2, false),
		NewDisplayConfig("HDMI-4K", Position{x: -4096, y: 0}, Resolution{x: 4096, y: 2304}, 1, false),
	}

	Application(configs)
}
