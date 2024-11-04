// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @channels @messaging @plugin @not_cloud

import * as TIMEOUTS from '../../../fixtures/timeouts';
import { matterpollPlugin } from '../../../utils/plugins';

describe('Header', () => {
    let testUser;
    let testTeam;
    let testChannel;
    before(() => {

        // # Setup and visit town-square
        cy.apiInitSetup().then(({ team: team, townSquareUrl: tUrl }) => {
            testTeam = team;
            cy.apiCreateChannel(team.id, 'testchannel', 'Test Channel').then(({ channel }) => {
                testChannel = channel;
                cy.apiCreateUser().then(({ user }) => {
                    testUser = user;
                    cy.apiPatchUserRoles(testUser.id, ['system_verified']);
                    cy.apiAddUserToTeam(team.id, testUser.id);
                    cy.apiAddUserToChannel(testChannel.id, testUser.id);
                });
                // magic scheme id
                cy.apiUpdateChannelScheme(testChannel.id, 'biggosyyyyyyyyyyyyyyyyyyyd');
            });
        });
    });

    it('verified user should not see readonly header', () => {
        cy.visit(`/${testTeam.name}/channels/${testChannel.name}`);
        cy.get('#channelHeaderTitle').should('be.visible').and('contain', testChannel.display_name);
        // * Verify readonly icon does not exist
        cy.get('.material-icons-outlined').should('not.exist');
    });

    it('unverified user should see readonly header', () => {
        cy.apiLogin(testUser);
        cy.visit(`/${testTeam.name}/channels/${testChannel.name}`);
        cy.get('#channelHeaderTitle').should('be.visible').and('contain', testChannel.display_name);

        // * Verify readonly icon exists
        cy.get('.material-icons-outlined').should('be.visible');
    });
});
