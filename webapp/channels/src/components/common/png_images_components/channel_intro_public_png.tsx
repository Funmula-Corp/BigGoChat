// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelIntroPublic from './image/channel_intro_public.png';

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
        src={channelIntroPublic}
    />
);

export default Png;
