// Copyright (c) 2024-present Funmula, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @biggo-chat

describe('Channle Creation', () => {
    let testTeam;
    let testUser;
    let otherUser;

    before(() => {
        cy.apiInitSetup().then(({team, user}) => {
            testTeam = team;
            testUser = user;
            cy.apiUpdateTeamMemberSchemeRole(team.id, user.id, {scheme_moderator: true, scheme_user: true});
            cy.apiCreateUser().then(({user: secondUser}) => {
                cy.apiPatchUserRoles(secondUser.id, ['system_user', 'system_verified']);
                cy.apiAddUserToTeam(testTeam.id, secondUser.id);
                otherUser = secondUser;
            });

            cy.visit(`/${team.name}/channels/town-square`);
        });
    });

    it('BC-T3 - only team admin can create public channel', () => {
        cy.apiLogin(otherUser).reload();

        cy.get('#SidebarContainer .AddChannelDropdown').should('be.visible').click();
        cy.get('#showNewChannel').should('not.exist');
        cy.get('#showMoreChannels').should('be.visible').click();

        cy.apiLogin(testUser).reload();

        cy.get('#SidebarContainer .AddChannelDropdown').should('be.visible').click();
        cy.get('#showNewChannel').should('be.visible').click();

        cy.get('#input_new-channel-modal-name').should('be.enabled').type('test public channel');
        cy.get('#new-channel-modal [type=submit]').should('be.enabled').click();

        cy.get('#new-channel-modal').should('not.exist');
    });

    it('BC-T4 - only team admin can create private channel', () => {
        cy.apiLogin(otherUser).reload();

        cy.get('#SidebarContainer .AddChannelDropdown').should('be.visible').click();
        cy.get('#showNewChannel').should('not.exist');
        cy.get('#showMoreChannels').should('be.visible').click();

        cy.apiLogin(testUser).reload();

        cy.get('#SidebarContainer .AddChannelDropdown').should('be.visible').click();
        cy.get('#showNewChannel').should('be.visible').click();

        cy.get('#input_new-channel-modal-name').should('be.enabled').type('test private channel');
        cy.get('#public-private-selector-button-P').should('be.enabled').click();
        cy.get('#new-channel-modal [type=submit]').should('be.enabled').click();

        cy.get('#new-channel-modal').should('not.exist');
    });
});
