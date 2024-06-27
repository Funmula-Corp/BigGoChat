// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import type {CSSProperties} from 'react';
import {useIntl} from 'react-intl';

export default function MattermostLogo(props: React.HTMLAttributes<HTMLSpanElement>) {
    const {formatMessage} = useIntl();
    return (
        <span {...props}>
            <svg
                version='1.1'
                x='0px'
                y='0px'
                viewBox='0 0 24 24'
                enableBackground='new 0 0 24 24'
                role='img'
                aria-label={formatMessage({id: 'generic_icons.mattermost', defaultMessage: 'BigGo Logo'})}
            >
                <g>
                    <path style={style}  d="M11,10.7c1.6,0,3.1-0.6,3.2-1.8c0.1-0.4,0-0.7-0.3-0.9c-0.4-0.5-1.5-0.6-3.2-0.5h-0.1l-0.4,3.2l0.2,0
                        C10.6,10.7,10.8,10.7,11,10.7z"/>
                    <path style={style} d="M10.2,12.3l-0.1,0l-0.5,4.1l0.3,0c1,0,2.1,0,3-0.2c0.9-0.2,1.6-0.7,1.7-1.5c0.1-0.5-0.1-1-0.5-1.4
                        C13.3,12.4,12,12.1,10.2,12.3z"/>
                    <path style={style} d="M16.9,2H7.1C3.9,2,2,3.9,2,7.1v9.7C2,20.1,3.9,22,7.1,22h9.7c3.3,0,5.1-1.9,5.1-5.1V7.1C22,3.9,20.1,2,16.9,2z M16.6,16.8
                        c-1.3,1.3-3.5,1.6-6.9,1.6c-1,0-2,0-3.2-0.1H6.4L8,5.6l0.7,0c3-0.1,5.9-0.3,7.5,1.3c0.6,0.6,0.9,1.3,0.9,2.1
                        c-0.1,0.9-0.5,1.8-1.3,2.4c0.9,0.6,1.5,1.6,1.7,2.7C17.7,15.1,17.4,16.1,16.6,16.8z"/>
                </g>
            </svg>
        </span>
    );
}

const style: CSSProperties = {
    fillRule: 'evenodd',
    clipRule: 'evenodd',
};
