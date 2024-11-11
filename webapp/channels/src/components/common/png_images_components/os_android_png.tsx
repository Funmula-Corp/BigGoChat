// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSAndroid from './image/os_android.png';

export function PNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={48}
                height={48}
                src={OSAndroid}
            />
        </span>
    );
}

export default PNG;
