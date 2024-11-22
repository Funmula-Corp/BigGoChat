// Copyright (c) 2024-present Funmula, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @biggo-chat

import {getRandomId} from '../../../utils';

describe('Team Creation', () => {
    let verifiedUser;
    let unverifiedUser;
    before(() => {
        cy.apiCreateUser().then(({user}) => {
            cy.apiPatchUserRoles(user.id, ['system_verified', 'system_user']);
            verifiedUser = user;
        });

        cy.apiCreateUser().then(({user}) => {
            unverifiedUser = user;
        });
    });

    it('BC-T1 - verified user can create teams', () => {
        cy.apiLogin(unverifiedUser);
        cy.reload();
        cy.visit('/');
        cy.get('#createNewTeamLink').should('not.exist');

        const suffix = getRandomId();
        cy.apiLogin(verifiedUser);
        cy.visit('/');
        cy.get('#createNewTeamLink').should('be.visible').click();
        cy.get('#teamNameInput').should('be.enabled').type(`test-team-${suffix}-0`);
        cy.get('#teamNameNextButton').click();
        cy.get('#teamURLFinishButton').click();
        cy.get('#SidebarContainer .title').should('have.text', `test-team-${suffix}-0`);
    });

    it('BC-T2 - verified user can only cretate 10 teams', () => {
        cy.apiLogin(verifiedUser);
        const suffix = getRandomId();
        for (let i = 1; i <= 9; i++) {
            cy.apiCreateTeam(`test-team${i}`, `test-team-${suffix}-${i}`);
        }

        cy.intercept('POST', '**/api/v4/teams', (req) => {
            req.alias = 'createTeam';
        });

        cy.visit('/create_team/display_name');
        cy.get('#teamNameInput').should('be.enabled').type(`test-team-${suffix}-10`);
        cy.get('#teamNameNextButton').click();
        cy.get('#teamURLFinishButton').click();

        // TODO: incomplete language pack, test response only
        cy.wait('@createTeam').its('response.statusCode').should('eq', 403);
    });
});
