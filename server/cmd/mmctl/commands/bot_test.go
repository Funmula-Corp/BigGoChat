// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"context"
	"errors"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	gomock "github.com/golang/mock/gomock"

	"git.biggo.com/Funmula/BigGoChat/server/v8/cmd/mmctl/printer"

	"github.com/spf13/cobra"
)

func (s *MmctlUnitTestSuite) TestBotCreateCmd() {
	s.Run("Should create a bot", func() {
		printer.Clean()

		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		mockBot := model.Bot{Username: botArg, DisplayName: "some-name", Description: "some-text"}

		s.client.
			EXPECT().
			CreateBot(context.TODO(), &mockBot).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		err := botCreateCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should create a bot with an access token", func() {
		printer.Clean()

		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		cmd.Flags().Bool("with-token", true, "")
		mockBot := model.Bot{Username: botArg, DisplayName: "some-name", Description: "some-text"}
		mockToken := model.UserAccessToken{Token: "token-id", Description: "autogenerated"}

		s.client.
			EXPECT().
			CreateBot(context.TODO(), &mockBot).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("no user found with the given email")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(model.UserFromBot(&mockBot), &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			CreateUserAccessToken(context.TODO(), mockBot.UserId, "autogenerated").
			Return(&mockToken, &model.Response{}, nil).
			Times(1)

		err := botCreateCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 2)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
		s.Require().Equal(&mockToken, printer.GetLines()[1])
	})

	s.Run("Should error when creating a bot", func() {
		printer.Clean()

		botArg := "a-bot"
		mockBot := model.Bot{Username: botArg, DisplayName: "", Description: ""}

		s.client.
			EXPECT().
			CreateBot(context.TODO(), &mockBot).
			Return(nil, &model.Response{}, errors.New("some-error")).
			Times(1)

		err := botCreateCmdF(s.client, &cobra.Command{}, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "could not create bot")
	})
}

func (s *MmctlUnitTestSuite) TestBotUpdateCmd() {
	s.Run("Should update a bot", func() {
		printer.Clean()

		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().String("username", "new-username", "")
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		cmd.Flags().Lookup("username").Changed = true
		cmd.Flags().Lookup("display-name").Changed = true
		cmd.Flags().Lookup("description").Changed = true
		mockBot := model.Bot{Username: "new-username", DisplayName: "some-name", Description: "some-text"}
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			PatchBot(context.TODO(), mockUser.Id, gomock.Any()).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		err := botUpdateCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should error when user not found bot", func() {
		printer.Clean()

		botArg := "a-bot"
		cmd := &cobra.Command{}
		cmd.Flags().String("username", "bot-username", "")
		cmd.Flags().Lookup("username").Changed = true

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUser(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botUpdateCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "unable to find user 'a-bot'")
	})

	s.Run("Should error when updating bot", func() {
		printer.Clean()

		botArg := "a-bot"
		cmd := &cobra.Command{}
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		cmd.Flags().Lookup("display-name").Changed = true
		cmd.Flags().Lookup("description").Changed = true
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			PatchBot(context.TODO(), mockUser.Id, gomock.Any()).
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botUpdateCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "could not update bot")
	})
}

func (s *MmctlUnitTestSuite) TestBotListCmd() {
	s.Run("Should list correctly all", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", true, "")
		mockBot := model.Bot{UserId: model.NewId(), Username: botArg, DisplayName: "some-name", Description: "some-text", OwnerId: model.NewId()}
		mockUser := model.User{Id: mockBot.OwnerId}

		s.client.
			EXPECT().
			GetBotsIncludeDeleted(context.TODO(), 0, 200, "").
			Return([]*model.Bot{&mockBot}, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUsersByIds(context.TODO(), []string{mockBot.OwnerId}).
			Return([]*model.User{&mockUser}, &model.Response{}, nil).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should list fail if one featching all bots requests fail", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", true, "")

		s.client.
			EXPECT().
			GetBotsIncludeDeleted(context.TODO(), 0, 200, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "Failed to fetch bots")
	})

	s.Run("Should list correctly orphaned", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", true, "")
		cmd.Flags().Bool("all", false, "")
		mockBot := model.Bot{UserId: model.NewId(), Username: botArg, DisplayName: "some-name", Description: "some-text", OwnerId: model.NewId()}
		mockUser := model.User{Id: mockBot.OwnerId}

		s.client.
			EXPECT().
			GetBotsOrphaned(context.TODO(), 0, 200, "").
			Return([]*model.Bot{&mockBot}, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUsersByIds(context.TODO(), []string{mockBot.OwnerId}).
			Return([]*model.User{&mockUser}, &model.Response{}, nil).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should list fail if one featching bots orphaned requests fail", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", true, "")
		cmd.Flags().Bool("all", false, "")

		s.client.
			EXPECT().
			GetBotsOrphaned(context.TODO(), 0, 200, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "Failed to fetch bots")
	})

	s.Run("Should list correctly bots", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", false, "")
		mockBot := model.Bot{UserId: model.NewId(), Username: botArg, DisplayName: "some-name", Description: "some-text", OwnerId: model.NewId()}
		mockUser := model.User{Id: mockBot.OwnerId}

		s.client.
			EXPECT().
			GetBots(context.TODO(), 0, 200, "").
			Return([]*model.Bot{&mockBot}, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUsersByIds(context.TODO(), []string{mockBot.OwnerId}).
			Return([]*model.User{&mockUser}, &model.Response{}, nil).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should list correctly bots with invalid ownerId", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", false, "")
		mockBot := model.Bot{UserId: model.NewId(), Username: botArg, DisplayName: "some-name", Description: "some-text", OwnerId: "Mr.Robot"}

		s.client.
			EXPECT().
			GetBots(context.TODO(), 0, 200, "").
			Return([]*model.Bot{&mockBot}, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUsersByIds(context.TODO(), []string{mockBot.OwnerId}).
			Return([]*model.User{}, &model.Response{}, nil).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should list fail if one fetching bots requests fail", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", false, "")

		s.client.
			EXPECT().
			GetBots(context.TODO(), 0, 200, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "Failed to fetch bots")
	})

	s.Run("Should list fail if fetching owners requests fail", func() {
		printer.Clean()
		botArg := "a-bot"

		cmd := &cobra.Command{}
		cmd.Flags().Bool("orphaned", false, "")
		cmd.Flags().Bool("all", false, "")
		mockBot := model.Bot{UserId: model.NewId(), Username: botArg, DisplayName: "some-name", Description: "some-text", OwnerId: model.NewId()}

		s.client.
			EXPECT().
			GetBots(context.TODO(), 0, 200, "").
			Return([]*model.Bot{&mockBot}, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUsersByIds(context.TODO(), []string{mockBot.OwnerId}).
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botListCmdF(s.client, cmd, []string{botArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "Failed to fetch bots")
	})
}

func (s *MmctlUnitTestSuite) TestBotDisableCmd() {
	s.Run("Should disable a bot", func() {
		printer.Clean()

		botArg := "a-bot"

		mockBot := model.Bot{Username: botArg, DisplayName: "some-name", Description: "some-text"}
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			DisableBot(context.TODO(), mockUser.Id).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		err := botDisableCmdF(s.client, &cobra.Command{}, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should error when user not found bot", func() {
		printer.Clean()

		botArg := "a-bot"

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUser(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botDisableCmdF(s.client, &cobra.Command{}, []string{botArg})
		s.Require().Error(err)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Contains(printer.GetErrorLines()[0], "can't find user 'a-bot'")
	})

	s.Run("Should error when disabling bot", func() {
		printer.Clean()

		botArg := "a-bot"
		cmd := &cobra.Command{}
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		cmd.Flags().Lookup("display-name").Changed = true
		cmd.Flags().Lookup("description").Changed = true
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			DisableBot(context.TODO(), mockUser.Id).
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botDisableCmdF(s.client, cmd, []string{botArg})
		s.Require().Error(err)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Contains(printer.GetErrorLines()[0], "could not disable bot 'a-bot'")
	})
}

func (s *MmctlUnitTestSuite) TestBotEnableCmd() {
	s.Run("Should enable a bot", func() {
		printer.Clean()

		botArg := "a-bot"

		mockBot := model.Bot{Username: botArg, DisplayName: "some-name", Description: "some-text"}
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			EnableBot(context.TODO(), mockUser.Id).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		err := botEnableCmdF(s.client, &cobra.Command{}, []string{botArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should error when user not found bot", func() {
		printer.Clean()

		botArg := "a-bot"

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUser(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botEnableCmdF(s.client, &cobra.Command{}, []string{botArg})
		s.Require().Error(err)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Contains(printer.GetErrorLines()[0], "can't find user 'a-bot'")
	})

	s.Run("Should error when enabling bot", func() {
		printer.Clean()

		botArg := "a-bot"
		cmd := &cobra.Command{}
		cmd.Flags().String("display-name", "some-name", "")
		cmd.Flags().String("description", "some-text", "")
		cmd.Flags().Lookup("display-name").Changed = true
		cmd.Flags().Lookup("description").Changed = true
		mockUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			EnableBot(context.TODO(), mockUser.Id).
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botEnableCmdF(s.client, cmd, []string{botArg})
		s.Require().Error(err)
		s.Require().Len(printer.GetErrorLines(), 1)
		s.Require().Contains(printer.GetErrorLines()[0], "could not enable bot 'a-bot'")
	})
}

func (s *MmctlUnitTestSuite) TestBotAssignCmd() {
	s.Run("Should assign a bot", func() {
		printer.Clean()

		botArg := "a-bot"
		userArg := "a-user"

		mockBot := model.Bot{Username: botArg, DisplayName: "some-name", Description: "some-text"}
		mockBotUser := model.User{Id: model.NewId()}
		mockNewOwner := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockBotUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), userArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), userArg, "").
			Return(&mockNewOwner, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			AssignBot(context.TODO(), mockBotUser.Id, mockNewOwner.Id).
			Return(&mockBot, &model.Response{}, nil).
			Times(1)

		err := botAssignCmdF(s.client, &cobra.Command{}, []string{botArg, userArg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 1)
		s.Require().Equal(&mockBot, printer.GetLines()[0])
	})

	s.Run("Should error when bot user not found", func() {
		printer.Clean()

		botArg := "a-bot"
		userArg := "a-user"

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUser(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botAssignCmdF(s.client, &cobra.Command{}, []string{botArg, userArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "unable to find user 'a-bot'")
	})

	s.Run("Should error when new owner not found", func() {
		printer.Clean()

		botArg := "a-bot"
		userArg := "a-user"

		mockBotUser := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockBotUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), userArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUser(context.TODO(), userArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), userArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botAssignCmdF(s.client, &cobra.Command{}, []string{botArg, userArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "unable to find user 'a-user'")
	})

	s.Run("Should error when assigning bot", func() {
		printer.Clean()

		botArg := "a-bot"
		userArg := "a-user"

		mockBotUser := model.User{Id: model.NewId()}
		mockNewOwner := model.User{Id: model.NewId()}

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), botArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), botArg, "").
			Return(&mockBotUser, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			GetUserByEmail(context.TODO(), userArg, "").
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(context.TODO(), userArg, "").
			Return(&mockNewOwner, &model.Response{}, nil).
			Times(1)

		s.client.
			EXPECT().
			AssignBot(context.TODO(), mockBotUser.Id, mockNewOwner.Id).
			Return(nil, &model.Response{}, errors.New("mock error")).
			Times(1)

		err := botAssignCmdF(s.client, &cobra.Command{}, []string{botArg, userArg})
		s.Require().NotNil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Contains(err.Error(), "can not assign bot 'a-bot' to user 'a-user'")
	})
}
