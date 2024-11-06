// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/// <reference types="cypress" />

// ***************************************************************
// Each command should be properly documented using JSDoc.
// See https://jsdoc.app/index.html for reference.
// Basic requirements for documentation are the following:
// - Meaningful description
// - Each parameter with `@params`
// - Return value with `@returns`
// - Example usage with `@example`
// Custom command should follow naming convention of having `external` prefix, e.g. `externalActivateUser`.
// ***************************************************************

declare namespace Cypress {
    interface Chainable {

        /**
         * Makes an external request as a sysadmin and activate/deactivate a user directly via API
         * @param {String} userId - The user ID
         * @param {Boolean} active - Whether to activate or deactivate - true/false
         *
         * @example
         *   cy.externalActivateUser('user-id', false);
         */
        externalActivateUser(userId: string, activate: boolean): Chainable;

        /**
         * Makes an external request as a sysadmin and send an invite to a user
         * @param {String} email - The email of the user to invite
         * @param {String} teamId - The team ID
         *
         * @example
         *   cy.externalSendInvite('user-email', 'team-id');
         */
        externalSendInvite(email: string, teamId: string): Chainable;

        /**
         * Makes an external request as a sysadmin and update user roles
         * @param {String} userId - The user ID
         * @param {Array<String>} roles - The roles to update
         *
         * @example
         *   cy.externalPatchUserRoles('user-id', ['system_user', 'system_verified']);
         */
        externalPatchUserRoles(userId: string, roles: string[]): Chainable;
    }
}
