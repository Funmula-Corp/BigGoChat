// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelIntroPrivate from './image/channel_intro_private.png';

type PngProps = {
    width?: number;
    height?: number;
}

const Png = (props: PngProps) => (
    <img
        width={props.width ?? '151'}
        height={props.height ?? '149'}
        src={channelIntroPrivate}
    />
);

export default Png;
