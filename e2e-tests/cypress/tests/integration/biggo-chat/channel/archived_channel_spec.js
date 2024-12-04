// Copyright (c) 2024-present Funmula, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @biggo-chat

describe('Archived Channel', () => {
    let testTeam;
    let testUser;
    let otherUser;

    before(() => {
        cy.apiInitSetup().then(({team, user}) => {
            testTeam = team;
            testUser = user;
            cy.apiUpdateTeamMemberSchemeRole(team.id, user.id, {scheme_admin: true, scheme_user: true});

            cy.apiCreateUser().then(({user: secondUser}) => {
                cy.apiPatchUserRoles(secondUser.id, ['system_user', 'system_verified']);
                cy.apiAddUserToTeam(testTeam.id, secondUser.id);
                otherUser = secondUser;
            });
        });
    });

    it('BC-T5 - admin can archive channel', () => {
        cy.apiLogin(testUser);
        cy.visit(`/${testTeam.name}`);

        cy.uiCreateChannel({});
        cy.uiAddUsersToCurrentChannel([otherUser.username]);

        cy.apiLogin(otherUser).reload();
        cy.get('#channelHeaderDropdownButton').click();
        cy.get('#channelArchiveChannel').should('not.exist');

        cy.getCurrentChannelId().then((channelId) => {
            cy.externalRequest({user: testUser, method: 'put', path: `channels/${channelId}/members/${otherUser.id}/schemeRoles`, data: {scheme_user: true, scheme_admin: true, scheme_verified: true}});
        });

        cy.get('#channelArchiveChannel').click();
        cy.get('#deleteChannelModalDeleteButton').click();
    });
});
