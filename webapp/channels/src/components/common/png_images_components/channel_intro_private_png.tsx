// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelIntroPrivate from './image/channel_intro_private.png';

type PngProps = {
    width?: number;
    height?: number;
    className?: string;
}

const Png = (props: PngProps) => (
    <img
        {...props}
        width={props.width ?? '250'}
        height={props.height ?? '250'}
        src={channelIntroPrivate}
    />
);

export default Png;
