// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Group: @channels @enterprise @accessibility

// BigGoChat does not have a feature to open the profile modal

import * as TIMEOUTS from '../../../../fixtures/timeouts';

describe('Verify Accessibility Support in Modals & Dialogs', () => {
    let testTeam;
    let testChannel;
    let testUser;

    before(() => {
        // * Check if server has license for Guest Accounts
        cy.apiRequireLicenseForFeature('GuestAccounts');

        cy.apiInitSetup({userPrefix: 'user000a'}).then(({team, channel, user}) => {
            testTeam = team;
            testChannel = channel;
            testUser = user;

            for (let i = 0; i < 20; i++) {
                cy.apiCreateUser().then(({user: newUser}) => {
                    cy.apiAddUserToTeam(testTeam.id, newUser.id).then(() => {
                        cy.apiAddUserToChannel(testChannel.id, newUser.id);
                    });
                });
            }
        });
    });

    beforeEach(() => {
        // # Login as sysadmin and visit the town-square
        cy.apiAdminLogin();
        cy.visit(`/${testTeam.name}/channels/town-square`);
    });

    it('MM-T1454 Accessibility Support in Different Modals and Dialog screen', () => {
        // * Verify the accessibility support in Profile Dialog
        // verifyUserMenuModal('Profile');

        // * Verify the accessibility support in Team Settings Dialog
        verifyMainMenuModal('Team Settings');

        // * Verify the accessibility support in Manage Members Dialog
        verifyMainMenuModal('Manage Members', `${testTeam.display_name} Members`);

        cy.visit(`/${testTeam.name}/channels/off-topic`);

        // * Verify the accessibility support in Channel Edit Header Dialog
        verifyChannelMenuModal('Edit Channel Header', 'Edit Header for Off-Topic');

        cy.wait(TIMEOUTS.TWO_SEC);

        // * Verify the accessibility support in Channel Edit Purpose Dialog
        verifyChannelMenuModal('Edit Channel Purpose', 'Edit Purpose for Off-Topic');

        // * Verify the accessibility support in Rename Channel Dialog
        verifyChannelMenuModal('Rename Channel');
    });

    it('MM-T1487 Accessibility Support in Manage Channel Members Dialog screen', () => {
        // # Visit test team and channel
        cy.visit(`/${testTeam.name}/channels/off-topic`);

        // # Open Channel Members Dialog
        cy.get('#channelHeaderDropdownIcon').click();
        cy.findByText('Manage Members').click().wait(TIMEOUTS.FIVE_SEC);

        // * Verify the accessibility support in Manage Members Dialog
        cy.get('#sidebar-right').within(() => {
            cy.contains('span', 'Managing Members');
            cy.contains('span', 'Off-Topic');

            // # Set focus on search input
            cy.findByPlaceholderText('Search members').
                focus().
                type(' {backspace}').
                wait(TIMEOUTS.HALF_SEC).
                tab({shift: true}).tab();
            cy.wait(TIMEOUTS.HALF_SEC);

            // # Press tab and verify focus on first user's profile image
            cy.focused().tab();
            cy.focused().within(() => {
                cy.findByAltText('sysadmin profile image').should('exist');
            });

            // # Press tab and verify focus on first user's username
            cy.focused().tab();
            cy.focused().within(() => {
                cy.findByText('sysadmin').should('exist');
            });

            // # Press tab and verify focus on second user's profile image
            cy.focused().tab();
            cy.focused().within(() => {
                cy.findByAltText(`${testUser.username} profile image`).should('exist');
            });

            // # Press tab and verify focus on second user's username
            cy.focused().tab();
            cy.focused().within(() => {
                cy.findByText(`${testUser.username}`).should('exist');
            })
            .should('have.class', 'dropdown-toggle')
            .and('contain', 'Member');

            // * Verify accessibility support in search total results
            cy.get('[aria-live="polite"]').should('exist');
        });
    });
});

function verifyMainMenuModal(menuItem, modalName) {
    cy.uiGetLHSHeader().click();
    verifyModal(menuItem, modalName);
}

function verifyChannelMenuModal(menuItem, modalName) {
    cy.get('#channelHeaderDropdownIcon').click();
    verifyModal(menuItem, modalName);
}

function verifyUserMenuModal(menuItem) {
    cy.uiGetSetStatusButton().click();
    verifyModal(menuItem);
}

function verifyModal(menuItem, modalName) {
    // * Verify that menu is open
    cy.findByRole('menu');

    // # Click menu item
    cy.findByText(menuItem).click();

    // * Verify the modal
    const name = modalName || menuItem;
    cy.findAllByRole('dialog').eq(1).within(() => {
        cy.get('.modal-title').contains(name);
        cy.uiClose();
    });
}
