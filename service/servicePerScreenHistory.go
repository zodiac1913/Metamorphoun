package service

import "Metamorphoun/config"

func clonePerScreenPics(pics []config.PicHistory) []config.PicHistory {
	if len(pics) == 0 {
		return nil
	}
	clones := make([]config.PicHistory, 0, len(pics))
	for idx, pic := range pics {
		pic.PerScreenPics = nil
		pic.PicNum = int16(idx)
		clones = append(clones, pic)
	}
	return clones
}

func attachPerScreenPics(pic config.PicHistory, perScreenPics []config.PicHistory) config.PicHistory {
	pic.PerScreenPics = clonePerScreenPics(perScreenPics)
	return pic
}
