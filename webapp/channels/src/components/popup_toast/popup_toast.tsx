// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import classNames from 'classnames';
import React, {useEffect, useCallback} from 'react';
import {CSSTransition} from 'react-transition-group';

import './popup_toast.scss';

type Props = {
    content: {
        message: string;
    };
    className?: string;
    onExited: () => void;
}

function PopupToast({content, onExited, className}: Props): JSX.Element {
    const closeToast = useCallback(() => {
        onExited();
    }, [onExited]);

    const toastContainerClassname = classNames('popup-toast', className);

    useEffect(() => {
        const timer = setTimeout(() => {
            onExited();
        }, 3000);

        return () => clearTimeout(timer);
    }, [onExited]);

    return (
        <CSSTransition
            in={Boolean(content)}
            classNames='toast'
            mountOnEnter={true}
            unmountOnExit={true}
            timeout={300}
            appear={true}
        >
            <div className={toastContainerClassname}>
                <span onClick={closeToast}>{content.message}</span>
            </div>
        </CSSTransition>
    );
}

export default React.memo(PopupToast);
