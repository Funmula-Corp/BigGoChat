// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @channels @channel

import {getAdminAccount} from '../../../support/env';
import {getRandomId} from '../../../utils';

describe('Archived channels', () => {
    let testTeam;
    let testUser;
    let testAdmin;

    before(() => {
        cy.apiUpdateConfig({
            TeamSettings: {
                ExperimentalViewArchivedChannels: true,
            },
        });

        cy.apiInitSetup({promoteNewUserAsAdmin: true}).then(({team, user}) => {
            testTeam = team;
            testAdmin = user;

            cy.apiCreateUser().then(({user}) => {
                testUser = user;
                cy.apiPatchUserRoles(user.id, ["system_verified"]);
                cy.apiAddUserToTeam(team.id, testUser.id);
            });
        });
    });
    
    it('MM-T1682-1 User cannot join an archived public channel by selecting a permalink to one of its posts', () => {
        // # Log in as another user
        cy.apiAdminLogin();

        // # Create a new channel
        cy.apiCreateChannel(testTeam.id, 'channel', 'channel').then(({channel}) => {
            // # Make a post in the new channel and get a permalink for it
            cy.postMessageAs({
                sender: getAdminAccount(),
                message: 'post',
                channelId: channel.id,
            }).then((post) => {
                const permalink = `/${testTeam.name}/pl/${post.id}`;

                // # Visit the channel
                cy.visit(`/${testTeam.name}/channels/${channel.name}`);

                // # Archive the channel
                cy.uiArchiveChannel();

                // # Log out and back in as the test user
                cy.apiLogin(testUser);
                cy.reload();

                // * Verify that we've logged in as the test user
                verifyUsername(testUser.username);

                // # Visit the permalink
                cy.visit(permalink);

                cy.findByText('Permalink belongs to a deleted message or to a channel to which you do not have access.').should('be.visible');
                cy.visit('/');
            });
        });
    });

    it('MM-T1682-2 Admin can join an archived public channel by selecting a permalink to one of its posts', () => {
        // # Log in as another user
        cy.apiAdminLogin();

        // # Create a new channel
        cy.apiCreateChannel(testTeam.id, 'channel', 'channel').then(({channel}) => {
            // # Make a post in the new channel and get a permalink for it
            cy.postMessageAs({
                sender: getAdminAccount(),
                message: 'post',
                channelId: channel.id,
            }).then((post) => {
                const permalink = `/${testTeam.name}/pl/${post.id}`;

                // # Visit the channel
                cy.visit(`/${testTeam.name}/channels/${channel.name}`);

                // # Archive the channel
                cy.uiArchiveChannel();

                // # Log out and back in as the test user
                cy.apiLogin(testAdmin);
                cy.reload();

                // * Verify that we've logged in as the test user
                verifyUsername(testAdmin.username);

                // # Visit the permalink
                cy.visit(permalink);

                verifyViewingArchivedChannel(channel);
            });
        });
    });

    it('MM-T1683-1 User cannot join an archived channel by selecting a link to channel', () => {
        // # Log in as another user
        cy.apiAdminLogin();

        // # Create a new channel
        cy.apiCreateChannel(testTeam.id, 'channel', 'channel').then(({channel}) => {
            const channelLink = `/${testTeam.name}/channels/${channel.name}`;

            // # Visit the channel and archive it
            cy.visit(channelLink);
            cy.uiArchiveChannel();

            // # Visit off-topic
            cy.visit(`/${testTeam.name}/channels/off-topic`);

            // # Make a post linking to the archived channel
            const linkText = `link ${getRandomId()}`;
            cy.getCurrentChannelId().then((currentChannelId) => {
                cy.postMessageAs({
                    sender: getAdminAccount(),
                    message: `This is a link: [${linkText}](${channelLink})`,
                    channelId: currentChannelId,
                });
            });

            // # Log out and back in as the test user
            cy.apiLogin(testUser);
            cy.reload();

            // # Visit off-topic
            cy.visit(`/${testTeam.name}/channels/off-topic`);

            // * Verify that we've logged in as the test user
            verifyUsername(testUser.username);

            // * Verify that the link exists and then click on it
            cy.contains('a', linkText).should('be.visible').click();

            cy.get('.loading__content', { timeout: 5000 }).should('be.visible');
        });
    });

    it('MM-T1683-2 Admin can join an archived channel by selecting a link to channel', () => {
        // # Log in as another user
        cy.apiAdminLogin();

        // # Create a new channel
        cy.apiCreateChannel(testTeam.id, 'channel', 'channel').then(({channel}) => {
            const channelLink = `/${testTeam.name}/channels/${channel.name}`;

            // # Visit the channel and archive it
            cy.visit(channelLink);
            cy.uiArchiveChannel();

            // # Visit off-topic
            cy.visit(`/${testTeam.name}/channels/off-topic`);

            // # Make a post linking to the archived channel
            const linkText = `link ${getRandomId()}`;
            cy.getCurrentChannelId().then((currentChannelId) => {
                cy.postMessageAs({
                    sender: getAdminAccount(),
                    message: `This is a link: [${linkText}](${channelLink})`,
                    channelId: currentChannelId,
                });
            });

            // # Log out and back in as the test user
            cy.apiLogin(testAdmin);
            cy.reload();

            // # Visit off-topic
            cy.visit(`/${testTeam.name}/channels/off-topic`);

            // * Verify that we've logged in as the test user
            verifyUsername(testAdmin.username);

            // * Verify that the link exists and then click on it
            cy.contains('a', linkText).should('be.visible').click();

            verifyViewingArchivedChannel(channel);
        });
    });
});

function verifyViewingArchivedChannel(channel) {
    // * Verify that we've switched to the correct channel and that the header contains the archived icon
    cy.get('#channelHeaderTitle').should('contain', channel.display_name);
    cy.get('#channelHeaderInfo .icon__archive').should('be.visible');

    // * Verify that the channel is visible in the sidebar with the archived icon
    cy.get(`#sidebarItem_${channel.name}`).should('be.visible').
        find('.icon-archive-outline').should('be.visible');

    // * Verify that the archived channel banner is visible at the bottom of the channel view
    cy.get('#channelArchivedMessage').should('be.visible');
}

function verifyUsername(username) {
    // * Verify that we've logged in as the test user
    cy.uiOpenUserMenu().findByText(`@${username}`);

    // # Close the user menu
    cy.uiGetSetStatusButton().click();
}
