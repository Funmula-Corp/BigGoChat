// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSLinux from './image/os_linux.png';

export function OSLinuxPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={48}
                height={48}
                src={OSLinux}
            />
        </span>
    );
}
