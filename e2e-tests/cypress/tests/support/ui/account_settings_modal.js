// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

Cypress.Commands.add('uiOpenProfileModal', (section = '') => {
    // # Open profile settings modal
    cy.uiOpenUserMenu('Profile');

    const profileSettingsModal = () => cy.findByRole('dialog', {name: 'Profile'}).should('be.visible');

    if (!section) {
        return profileSettingsModal();
    }

    // # Click on a particular section
    cy.findByRoleExtended('tab', {name: section}).should('be.visible').click();

    return profileSettingsModal();
});

// BigGoChat profile opens in a new window, so use uiOpenProfilePage instead of uiOpenProfileModal
Cypress.Commands.add('uiOpenProfilePage', () => {
    cy.uiOpenUserMenu('Profile').then(() => {
        cy.window().then((win) => {
            cy.stub(win, 'open').callsFake((url) => {
                expect(url).to.eq('https://account.biggo.com/setting/');
                return null;
            });
        });
    });
});

Cypress.Commands.add('verifyAccountNameSettings', (firstname, lastname) => {
    // # Go to Profile
    cy.uiOpenProfileModal('Profile Settings');

    // * Check name value
    cy.get('#nameDesc').should('have.text', `${firstname} ${lastname}`);
    cy.uiClose();
});

Cypress.Commands.add('uiChangeGenericDisplaySetting', (setting, option) => {
    cy.uiOpenSettingsModal('Display');
    cy.get(setting).scrollIntoView();
    cy.get(setting).click();
    cy.get('.section-max').scrollIntoView();

    cy.get(option).check().should('be.checked');

    cy.uiSaveAndClose();
});

/*
 * Change the message display setting
 * @param {String} setting - as 'STANDARD' or 'COMPACT'
 */
Cypress.Commands.add('uiChangeMessageDisplaySetting', (setting = 'STANDARD') => {
    const SETTINGS = {STANDARD: '#message_displayFormatA', COMPACT: '#message_displayFormatB'};
    cy.uiChangeGenericDisplaySetting('#message_displayTitle', SETTINGS[setting]);
});

/*
 * Change the collapsed reply threads display setting
 * @param {String} setting - as 'OFF' or 'ON'
 */
Cypress.Commands.add('uiChangeCRTDisplaySetting', (setting = 'OFF') => {
    const SETTINGS = {
        ON: '#collapsed_reply_threadsFormatA',
        OFF: '#collapsed_reply_threadsFormatB',
    };

    cy.uiChangeGenericDisplaySetting('#collapsed_reply_threadsTitle', SETTINGS[setting]);
});
