// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {Role} from '@mattermost/types/roles';
import type {GlobalState} from '@mattermost/types/store';
import type {UserProfile} from '@mattermost/types/users';

import {createSelector} from 'mattermost-redux/selectors/create_selector';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/common';
import { General } from 'mattermost-redux/constants';

export type PermissionsOptions = {
    channel?: string;
    team?: string;
    permission: string;
};

export function getRoles(state: GlobalState) {
    return state.entities.roles.roles;
}

export const getMySystemRoles: (state: GlobalState) => Set<string> = createSelector(
    'getMySystemRoles',
    getCurrentUser,
    (user: UserProfile) => {
        if (user) {
            return new Set<string>(user.roles.split(' '));
        }

        return new Set<string>();
    },
);

export const getMySystemPermissions: (state: GlobalState) => Set<string> = createSelector(
    'getMySystemPermissions',
    getMySystemRoles,
    getRoles,
    (mySystemRoles: Set<string>, allRoles: any) => {
        return getPermissionsForRoles(allRoles, mySystemRoles);
    },
);

export function haveISystemPermission(state: GlobalState, options: PermissionsOptions) {
    return getMySystemPermissions(state).has(options.permission);
}

export function getPermissionsForRoles(allRoles: Record<string, Role>, roleSet: Set<string>) {
    const permissions = new Set<string>();

    for (const roleName of roleSet) {
        const role = allRoles[roleName];

        if (!role) {
            continue;
        }

        for (const permission of role.permissions) {
            permissions.add(permission);
        }
    }

    return permissions;
}

export function haveIVerified(state: GlobalState) {
    return getCurrentUser(state)
        .roles
        .split(' ')
        .some(role =>
            role == General.SYSTEM_ADMIN_ROLE ||
            role == General.TEAM_ADMIN_ROLE ||
            role == General.CHANNEL_ADMIN_ROLE ||
            role == General.SYSTEM_VERIFIED_ROLE ||
            role == General.TEAM_VERIFIED_ROLE ||
            role == General.CHANNEL_VERIFIED_ROLE
        );
}
