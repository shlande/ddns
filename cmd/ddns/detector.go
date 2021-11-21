package main

import (
	"github.com/shlande/ddns"
	"github.com/shlande/ddns/detector"
	"github.com/sirupsen/logrus"
	"strings"
)

type DetectorBuilder func(input string) ddns.Detector
type DetectorV6Builder func(input string) ddns.DetectorV6

var detectors = map[string]DetectorBuilder{
	"xd": func(input string) ddns.Detector {
		return detector.Xd{}
	},
	"public": func(input string) ddns.Detector {
		return detector.Public{}
	},
	"device": func(input string) ddns.Detector {
		detector, err := detector.NewDeviceGetter(input)
		if err != nil {
			logrus.Fatalf("无法创建detector:%v", err)
		}
		return detector
	},
}

var detectorV6s = map[string]DetectorV6Builder{
	"device": func(input string) ddns.DetectorV6 {
		detector, err := detector.NewDeviceGetter(input)
		if err != nil {
			logrus.Fatalf("无法创建detector:%v", err)
		}
		return detector
	},
}

func buildDetector(detect string) ddns.Detector {
	raw := strings.Split(detect, "=")
	if len(raw) < 1 || len(raw) > 2 {
		logrus.Fatal("参数数量应在1-2之间")
	}
	detectorName, input := raw[0], raw[1]
	for k, v := range detectors {
		if k == detectorName {
			return v(input)
		}
	}
	logrus.Fatal("没有找到detector")
	return nil
}

func buildDetectorV6(detect string) ddns.DetectorV6 {
	raw := strings.Split(detect, "=")
	var detectorName, input string
	switch len(raw) {
	case 1:
		detectorName = raw[0]
	case 2:
		detectorName, input = raw[0], raw[1]
	default:
		logrus.Fatal("参数数量应在1-2之间")
	}

	for k, v := range detectorV6s {
		if k == detectorName {
			return v(input)
		}
	}
	logrus.Fatal("没有找到detector或者该detector不支持当前网络类型")
	return nil
}
