// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {ThunkActionFunc} from 'mattermost-redux/types/actions';

import icon50 from 'images/icon50x50.png';
import iconWS from 'images/icon_WS.png';
import Constants from 'utils/constants';
import * as UserAgent from 'utils/user_agent';

export type NotificationResult = {
    status: 'error' | 'not_sent' | 'success' | 'unsupported';
    reason?: string;
    data?: string;
}

let requestedNotificationPermission = Boolean('Notification' in window && Notification.permission !== 'default');

// showNotification displays a platform notification with the configured parameters.
//
// If successful in showing a notification, it resolves with a callback to manually close the
// notification. If no error occurred but the user did not grant permission to show notifications, it
// resolves with a no-op callback. Notifications that do not require interaction will be closed automatically after
// the Constants.DEFAULT_NOTIFICATION_DURATION. Not all platforms support all features, and may
// choose different semantics for the notifications.

export interface ShowNotificationParams {
    title: string;
    body: string;
    requireInteraction: boolean;
    silent: boolean;
    onClick?: (this: Notification, e: Event) => any | null;
}

export function showNotification(
    {
        title,
        body,
        requireInteraction,
        silent,
        onClick,
    }: ShowNotificationParams = {
        title: '',
        body: '',
        requireInteraction: false,
        silent: false,
    },
): ThunkActionFunc<Promise<NotificationResult & {callback: () => void}>> {
    return async () => {
        let icon = icon50;
        if (UserAgent.isEdge()) {
            icon = iconWS;
        }

        if (!isNotificationAPISupported()) {
            throw new Error('Notification API is not supported');
        }

        if (Notification.permission !== 'granted') {
            if (requestedNotificationPermission) {
                // User didn't allow notifications
                return {status: 'not_sent', reason: 'notifications_permission_previously_denied', data: Notification.permission, callback: () => {}};
            }

            requestedNotificationPermission = true;

            let permission = await Notification.requestPermission();
            if (typeof permission === 'undefined') {
                // Handle browsers that don't support the promise-based syntax.
                permission = await new Promise((resolve) => {
                    Notification.requestPermission(resolve);
                });
            }

            if (permission !== 'granted') {
                // User has denied notification for the site
                return {status: 'not_sent', reason: 'notifications_permission_denied', data: permission, callback: () => {}};
            }
        }

        const notification = new Notification(title, {
            body,
            tag: body,
            icon,
            requireInteraction,
            silent,
        });

        if (onClick) {
            notification.onclick = onClick;
        }

        notification.onerror = () => {
            throw new Error('Notification failed to show.');
        };

        // Mac desktop app notification dismissal is handled by the OS
        if (!requireInteraction && !UserAgent.isMacApp()) {
            setTimeout(() => {
                notification.close();
            }, Constants.DEFAULT_NOTIFICATION_DURATION);
        }

        return {
            status: 'success',
            callback: () => {
                notification.close();
            },
        };
    };
}

export function isNotificationAPISupported() {
    return ('Notification' in window) && (typeof Notification.requestPermission === 'function');
}

export async function requestNotificationPermission() {
    if (!isNotificationAPISupported()) {
        return null;
    }

    try {
        const notificationPermission = await Notification.requestPermission();
        return notificationPermission;
    } catch (error) {
        return null;
    }
}
