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
import { reUrl } from '../../../utils';
import { initial } from 'lodash';

describe('Header', () => {
    let testUser;
    let testUser2;
    let testUser3;
    let testTeam;
    let testTeam2;
    let testChannel;
    let dm;
    let adminUser;
    let permalink2;
    let permalink3;
    before(() => {
        // # Setup and visit town-square
        cy.apiInitSetup().then(({ team: team, townSquareUrl: tUrl, user: user }) => {
            adminUser = user;
            testTeam = team;
            cy.apiCreateTeam('new-team', 'New Team').then(({ team }) => {
                testTeam2 = team;
            }).apiCreateChannel(testTeam.id, 'test-channel', 'Test Channel', 'O').then(({ channel }) => {
                testChannel = channel;
            }).apiCreateUser({ prefix: 'one' }).then(({ user }) => {
                testUser = user;
                cy.apiAddUserToTeam(team.id, testUser.id);
                cy.apiCreateDirectChannel([testUser.id, adminUser.id]).then(({ channel }) => {
                    dm = channel;
                });
            }).apiCreateUser({ prefix: 'two' }).then(({ user }) => {
                testUser2 = user;
                cy.apiAddUserToTeam(testTeam.id, testUser2.id)
                cy.externalSendInvite(testUser2.email, testTeam2.id)
                    .getRecentEmail(testUser2).then((data) => {
                        console.log(data.body);
                        const matched = data.body[3].match(reUrl);
                        // assert(matched.length > 0);
                        permalink2 = matched[0];
                    })
            });
        });
    });

    it('verified user should be able to post after verified', () => {
        cy.apiLogout();
        cy.uiLogin(testUser).then(() => {
            cy.visit(`/${testTeam.name}/channels/town-square`)
                .get('.AdvancedTextEditor__verified-button').should('be.visible')
                .visit(`/${testTeam.name}/channels/${dm.name}`)
                .get('.AdvancedTextEditor__verified-button').should('be.visible')
                .visit(`/${testTeam.name}/channels/${testChannel.name}`)
                .get('.AdvancedTextEditor__verified-button').should('be.visible')
                .externalPatchUserRoles(testUser.id, ['system_user', 'system_verified']);
            // server sends a websocket event to the frontend to update the verified state
            // it always failed when the user is in town-square channel
            cy.get('.AdvancedTextEditor__verified-button').should('not.exist');
            cy.postMessage('i can post now');
        });
    });

    it('test user2 click the invite link after verified', () => {
        cy.apiLogout();
        cy.uiLogin(testUser2).visit(`/${testTeam.name}/channels/town-square`)
            .get('#channelHeaderTitle').should('be.visible')
        cy.get('.AdvancedTextEditor__verified-button').should('be.visible');
        cy.externalPatchUserRoles(testUser2.id, ['system_user', 'system_verified'])
        cy.visit(permalink2).get('#channelHeaderTitle').should('be.visible')
            .get('.AdvancedTextEditor__verified-button').should('not.exist');
        cy.goToDm(adminUser.username).postMessage('testuser2 can post to DMnow')
    });
});
