// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import channelIntroPublic from './image/channel_intro_public.png';

const PNG = (props: Partial<React.ImgHTMLAttributes<HTMLImageElement>>) => (
    <img
        {...props}
        width={props.width ?? '250'}
        height={props.height ?? '250'}
        src={channelIntroPublic}
    />
);

export default PNG;
