// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package featureflag

import (
	"fmt"

	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
)

type splitLogger struct {
	wrappedLog *mlog.Logger
}

func (s *splitLogger) Error(msg ...any) {
	s.wrappedLog.Error(fmt.Sprint(msg...))
}

func (s *splitLogger) Warning(msg ...any) {
	s.wrappedLog.Warn(fmt.Sprint(msg...))
}

// Ignoring more verbose messages from split
func (s *splitLogger) Info(msg ...any) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}

func (s *splitLogger) Debug(msg ...any) {
	//s.wrappedLog.Debug(fmt.Sprint(msg...))
}

func (s *splitLogger) Verbose(msg ...any) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}
