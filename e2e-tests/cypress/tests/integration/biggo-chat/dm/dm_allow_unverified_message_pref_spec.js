// Copyright (c) 2024-present Funmula, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @biggo-chat

describe('DM allow_unverified_message pref', () => {
    let userA;
    let unverifiedUser;
    let userC;
    let team1;
    let testChannelUrl;

    before(() => {
        cy.apiInitSetup().then(({team, user}) => {
            userA = user;
            team1 = team;
            cy.apiPatchUserRoles(userA.id, ["system_verified"]);
            testChannelUrl = `/${team.name}/channels/town-square`;

            cy.apiCreateUser().then(({user: otherUser}) => {
                unverifiedUser = otherUser;
                cy.apiAddUserToTeam(team.id, unverifiedUser.id);
            });
            cy.apiCreateUser().then(({user: otherUser}) => {
                userC = otherUser;
                cy.apiSaveAllowUnverifiedMessage(userC.id, true);
            });
        });
    });

    it('unverified user cannot send dm by default', () => {
        // # Log in as Unverified User
        cy.apiLogin(unverifiedUser);
        cy.visit(testChannelUrl);

        cy.goToDm(userA.username).get('.AdvancedTextEditor__verified-button').should('be.visible');
        cy.goToDm(userC.username).get('.AdvancedTextEditor__verified-button').postMessage('userC allow me to post.');
    });
});
