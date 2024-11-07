// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import Mentions from './image/mentions.png';

export function MentionsPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={17}
                height={17}
                src={Mentions}
            />
        </span>
    );
}
