// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package storetest

import (
	"math/rand"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

func MakeEmail() string {
	return "success_" + model.NewId() + "@simulator.amazonses.com"
}

func GenerateTestMobilephone() string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "886" + string(b)
}
