// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelIntroPublic from './image/channel_intro_public.png';

type PngProps = {
    width?: number;
    height?: number;
}

const Png = (props: PngProps) => (
    <img
        width={props.width ?? '151'}
        height={props.height ?? '149'}
        src={channelIntroPublic}
    />
);

export default Png;
