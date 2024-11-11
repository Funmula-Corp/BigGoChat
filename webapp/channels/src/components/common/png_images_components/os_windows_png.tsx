// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSWindows from './image/os_windows.png';

export function OSWindowsPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={48}
                height={48}
                src={OSWindows}
            />
        </span>
    );
}
