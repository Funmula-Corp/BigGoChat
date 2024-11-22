// Copyright (c) 2024-present Funmula, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @biggo-chat

describe('Channel Join', () => {
    let testTeam;
    let testUser;
    let otherUser;
    let otherUser2;

    let publicChannel;
    let privateChannel;

    before(() => {
        cy.apiInitSetup().then(({team, user}) => {
            testTeam = team;
            testUser = user;
            cy.apiUpdateTeamMemberSchemeRole(team.id, user.id, {scheme_admin: true, scheme_user: true});
            cy.apiCreateChannel(team.id, 'private-channel', 'Private Channel', 'P').then(({channel}) => {
                cy.apiAddUserToChannel(channel.id, testUser.id);
                privateChannel = channel;
            });
            cy.apiCreateChannel(team.id, 'public-channel', 'Public Channel', 'O').then(({channel}) => {
                publicChannel = channel;
            });

            cy.apiCreateUser().then(({user: secondUser}) => {
                cy.apiPatchUserRoles(secondUser.id, ['system_user', 'system_verified']);
                cy.apiAddUserToTeam(testTeam.id, secondUser.id);
                otherUser = secondUser;
            });

            cy.apiCreateUser().then(({user: secondUser}) => {
                cy.apiPatchUserRoles(secondUser.id, ['system_user', 'system_verified']);
                cy.apiAddUserToTeam(testTeam.id, secondUser.id);
                otherUser2 = secondUser;
            });
        });
    });

    it('BC-T7 - anyone can join public channel', () => {
        cy.apiLogin(otherUser);
        cy.visit(`/${testTeam.name}`);

        cy.uiBrowseOrCreateChannel('Browse channels').click();
        cy.get('#moreChannelsList').then(() => {
            cy.findAllByText(publicChannel.display_name).click();
        });

        cy.uiGetChannelHeaderButton().should('contain.text', publicChannel.display_name);
    });

    it('BC-T8 - only admin can add user to private channel', () => {
        cy.apiLogin(testUser);
        cy.reload();
        cy.uiClickSidebarItem(privateChannel.name);

        cy.uiGetChannelHeaderButton().should('contain.text', privateChannel.display_name);

        cy.uiAddUsersToCurrentChannel([otherUser.username]);

        cy.apiLogin(otherUser);
        cy.reload();
        cy.uiGetChannelHeaderButton().click();
        cy.get('#channelAddMembers').should('not.exist').then(() => {
            cy.uiGetChannelHeaderButton().click();
        });

        cy.externalRequest({user: testUser, method: 'put', path: `channels/${privateChannel.id}/members/${otherUser.id}/schemeRoles`, data: {scheme_user: true, scheme_admin: true, scheme_verified: true}});

        cy.uiAddUsersToCurrentChannel([otherUser2.username]);
        cy.uiGetChannelMemberButton().click();
        cy.get('#sidebar-right').should('contain.text', otherUser2.username);
    });
});
