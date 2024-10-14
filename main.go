package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/codec/openh264"
	// "github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/prop"

	_ "github.com/getrec/whip-go/driver/audiotest"
	_ "github.com/getrec/whip-go/driver/videotest"
	"github.com/pion/webrtc/v3"
)

func main() {
	iceServer := flag.String("i", "stun:stun.l.google.com:19302", "ice server")
	token := flag.String("t", "whip-go", "publishing token")
	videoCodec := flag.String("vc", "h264", "video codec vp8|h264")
	audioBitrate := flag.Int("ab", 160_000, "video bitrate in bits per second")
	videoBitrate := flag.Int("vb", 3_000_000, "video bitrate in bits per second")
	videoWidth := flag.Int("vw", 1280, "video width in pixels")
	videoHeight := flag.Int("vh", 720, "video height in pixels")
	videoFrameRate := flag.Float64("vf", 30, "video frame rate in frames per second")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Invalid number of arguments, pass the publishing url as the first argument")
	}

	mediaEngine := webrtc.MediaEngine{}
	whip := NewWHIPClient(flag.Args()[0], *token)

	// configure codec specific parameters
	vpxParams, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vpxParams.BitRate = *videoBitrate

	opusParams, err := opus.NewParams()
	if err != nil {
		panic(err)
	}
	opusParams.BitRate = *audioBitrate

	openh264Params, err := openh264.NewParams()
	if err != nil {
		panic(err)
	}
	openh264Params.BitRate = *videoBitrate
	// openh264Params.Preset = openh264.PresetSuperfast

	var videoCodecSelector mediadevices.CodecSelectorOption
	if *videoCodec == "vp8" {
		videoCodecSelector = mediadevices.WithVideoEncoders(&vpxParams)
	} else {
		videoCodecSelector = mediadevices.WithVideoEncoders(&openh264Params)
	}
	var stream mediadevices.MediaStream

	codecSelector := mediadevices.NewCodecSelector(
		videoCodecSelector,
		mediadevices.WithAudioEncoders(&opusParams),
	)
	codecSelector.Populate(&mediaEngine)

	stream, err = mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.Width = prop.Int(*videoWidth)
			constraint.Height = prop.Int(*videoHeight)
			constraint.FrameRate = prop.Float(*videoFrameRate)
		},
		Audio: func(constraint *mediadevices.MediaTrackConstraints) {},
		Codec: codecSelector,
	})
	if err != nil {
		log.Fatal("Unexpected error capturing test source. ", err)
	}

	iceServers := []webrtc.ICEServer{
		{
			URLs: []string{*iceServer},
		},
	}

	whip.Publish(stream, mediaEngine, iceServers, true)

	fmt.Println("Press 'Enter' to finish...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	whip.Close(true)
}
