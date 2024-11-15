// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import noDraft from './image/no_draft.png';

type PngProps = {
    width?: number;
    height?: number;
    className?: string;
}

const PNG = (props: PngProps) => (
    <img
        {...props}
        width={props.width ?? 250}
        height={props.height ?? 250}
        src={noDraft}
    />
);

export default PNG;
