// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @channels @menu

describe('Main menu', () => {
    let testTeam;
    let testUser;

    before(() => {
        cy.apiInitSetup().then(({team, user}) => {
            testTeam = team;
            testUser = user;
        });
    });

    // BigGoChat does not have a feature to open the profile modal
    it.skip('MM-T711_1 - Click on menu item should toggle the menu', () => {
        cy.apiLogin(testUser);
        cy.visit(`/${testTeam.name}/channels/town-square`);

        cy.uiOpenProfileModal('Profile Settings');
        cy.findByRole('set status').should('not.exist');
    });

    it('MM-T711_2 - Click on menu divider shouldn\'t toggle the menu', () => {
        cy.apiLogin(testUser);
        cy.visit(`/${testTeam.name}/channels/town-square`);

        cy.uiOpenUserMenu();

        cy.findByRole('menu').should('exist').find('.menu-divider:visible').last().click();
        cy.findByRole('menu').should('exist');
    });

    it('should show integrations option for system administrator', () => {
        cy.apiAdminLogin();
        cy.visit(`/${testTeam.name}/channels/town-square`);

        cy.uiOpenProductMenu();

        cy.findByRole('menu').findByText('Integrations').should('be.visible');
    });

    it('should not show integrations option for team member without permissions', () => {
        cy.apiLogin(testUser);
        cy.visit(`/${testTeam.name}/channels/town-square`);

        cy.uiOpenProductMenu();

        cy.findByRole('menu').findByText('Integrations').should('not.exist');
    });
});
