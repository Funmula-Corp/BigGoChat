package settings

import (
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func stringsToOptions(in []string) []*model.PostActionOptions {
	out := make([]*model.PostActionOptions, len(in))
	for i, o := range in {
		out[i] = &model.PostActionOptions{
			Text:  o,
			Value: o,
		}
	}
	return out
}
