// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {memo, useEffect, useState} from 'react';
import {useIntl} from 'react-intl';
import {useDispatch, useSelector} from 'react-redux';

import type {FileInfo} from '@mattermost/types/files';

import {getFilePublicLink} from 'mattermost-redux/actions/files';
import {getFilePublicLink as selectFilePublicLink} from 'mattermost-redux/selectors/entities/files';

import CopyButton from 'components/copy_button';
import ExternalLink from 'components/external_link';
import WithTooltip from 'components/with_tooltip';

import Constants, {FileTypes} from 'utils/constants';
import {copyToClipboard, getFileType} from 'utils/utils';

import type {GlobalState} from 'types/store';

import {isFileInfo} from '../types';
import type {LinkInfo} from '../types';

import './file_preview_modal_main_actions.scss';
import OverlayTrigger from 'components/overlay_trigger';
import Tooltip from 'components/tooltip';

interface Props {
    usedInside?: 'Header' | 'Footer';
    showOnlyClose?: boolean;
    showClose?: boolean;
    showPublicLink?: boolean;
    filename: string;
    fileURL: string;
    fileInfo: FileInfo | LinkInfo;
    enablePublicLink: boolean;
    canDownloadFiles: boolean;
    canCopyContent: boolean;
    handleModalClose: () => void;
    content: string;
}

const FilePreviewModalMainActions: React.FC<Props> = (props: Props) => {
    const intl = useIntl();

    const tooltipPlacement = props.usedInside === 'Header' ? 'bottom' : 'top';
    const selectedFilePublicLink = useSelector((state: GlobalState) => selectFilePublicLink(state)?.link);
    const dispatch = useDispatch();
    const [publicLinkCopied, setPublicLinkCopied] = useState(false);

    useEffect(() => {
        if (isFileInfo(props.fileInfo) && props.enablePublicLink) {
            dispatch(getFilePublicLink(props.fileInfo.id));
        }
    }, [props.fileInfo, props.enablePublicLink]);
    const copyPublicLink = () => {
        copyToClipboard(selectedFilePublicLink ?? '');
        setPublicLinkCopied(true);
    };
    const openPublicLink = () => {
        if (selectedFilePublicLink) {
            globalThis.window.open(selectedFilePublicLink, '_blank')
        }
    };

    const closeMessage = intl.formatMessage({
        id: 'full_screen_modal.close',
        defaultMessage: 'Close',
    });
    const closeButton = (
        <WithTooltip
            id='close-icon-tooltip'
            title={closeMessage}
            placement={tooltipPlacement}
            key='publicLink'
        >
            <button
                className='file-preview-modal-main-actions__action-item'
                onClick={props.handleModalClose}
                aria-label={closeMessage}
            >
                <i className='icon icon-close'/>
            </button>
        </WithTooltip>
    );

    let publicTooltipMessage;
    if (publicLinkCopied) {
        publicTooltipMessage = intl.formatMessage({
            id: 'file_preview_modal_main_actions.public_link-copied',
            defaultMessage: 'Public link copied',
        });
    } else {
        publicTooltipMessage = intl.formatMessage({
            id: 'view_image_popover.publicLink',
            defaultMessage: 'Get a public link',
        });
    }
    const publicLink = (
        <WithTooltip
            id='link-variant-icon-tooltip.text'
            key='filePreviewPublicLink'
            placement={tooltipPlacement}
            title={publicTooltipMessage}
            shouldUpdatePosition={true}
            onExit={() => setPublicLinkCopied(false)}
        >
            <a
                href='#'
                className='file-preview-modal-main-actions__action-item'
                onClick={copyPublicLink}
                aria-label={publicTooltipMessage}
            >
                <i className='icon icon-link-variant'/>
            </a>
        </WithTooltip>
    );
    // TODO i18n
    const publicLinkPage = (
        <OverlayTrigger
            delayShow={Constants.OVERLAY_TIME_DELAY}
            key='link'
            placement={tooltipPlacement}
            overlay={
                <Tooltip id='link-variant-icon-tooltip'>
                    Open image in new tab
                </Tooltip>
            }
        >
            <a
                href='#'
                className='file-preview-modal-main-actions__action-item'
                onClick={openPublicLink}
            >
                <i className='icon icon-magnify-plus'/>
            </a>
        </OverlayTrigger>
    );

    const downloadMessage = intl.formatMessage({
        id: 'view_image_popover.download',
        defaultMessage: 'Download',
    });
    const download = (
        <WithTooltip
            id='download-icon-tooltip.text'
            key='download'
            placement={tooltipPlacement}
            title={downloadMessage}
        >
            <ExternalLink
                href={props.fileURL}
                className='file-preview-modal-main-actions__action-item'
                location='file_preview_modal_main_actions'
                download={props.filename}
                aria-label={downloadMessage}
            >
                <i className='icon icon-download-outline'/>
            </ExternalLink>
        </WithTooltip>
    );

    const copy = (
        <CopyButton
            className='file-preview-modal-main-actions__action-item'
            isForText={getFileType(props.fileInfo.extension) === FileTypes.TEXT}
            placement={tooltipPlacement}
            content={props.content}
        />
    );
    return (
        <div className='file-preview-modal-main-actions__actions'>
            {!props.showOnlyClose && props.canCopyContent && copy}
            {!props.showOnlyClose && props.enablePublicLink && props.showPublicLink && !!selectedFilePublicLink && publicLinkPage}
            {!props.showOnlyClose && props.enablePublicLink && props.showPublicLink && publicLink}
            {!props.showOnlyClose && props.canDownloadFiles && download}
            {props.showClose && closeButton}
        </div>
    );
};

FilePreviewModalMainActions.defaultProps = {
    showOnlyClose: false,
    usedInside: 'Header',
    showClose: true,
    showPublicLink: true,
};

export default memo(FilePreviewModalMainActions);
