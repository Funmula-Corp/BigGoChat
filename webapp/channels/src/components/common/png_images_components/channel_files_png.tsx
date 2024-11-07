// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelFiles from './image/channel_files.png';

export function ChannelFilesPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={17}
                height={17}
                src={channelFiles}
            />
        </span>
    );
}
