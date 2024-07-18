// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Modal} from 'react-bootstrap';
import {FormattedMessage} from 'react-intl';

import {General} from 'mattermost-redux/constants';

import {trackEvent} from 'actions/telemetry_actions.jsx';

import FormattedMarkdownMessage from 'components/formatted_markdown_message';

import Constants from 'utils/constants';

type Props = {
    channelDisplayName: string;
    channelId: string;

    /**
     * Function injected by ModalController to be called when the modal can be unmounted
     */
    onExited: () => void;

    actions: {
        updateChannelPrivacy: (channelId: string, privacy: string) => void;
    };
}

type State = {
    show: boolean;
    confirmDisplayName: string;
}

export default class ConvertChannelModal extends React.PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            show: true,
            confirmDisplayName: '',
        };
    }

    canConvert = () => {
        return this.state.confirmDisplayName == this.props.channelDisplayName;
    };

    handleConvert = () => {
        const {actions, channelId} = this.props;
        if (channelId.length !== Constants.CHANNEL_ID_LENGTH) {
            return;
        }

        actions.updateChannelPrivacy(channelId, General.PRIVATE_CHANNEL);
        trackEvent('actions', 'convert_to_private_channel', {channel_id: channelId});
        this.onHide();
    };

    onUpdateConfirmName = (event: any) => {
        this.setState({
            confirmDisplayName: event.target.value,
        })
    };

    onHide = () => {
        this.setState({show: false});
    };

    render() {
        const {
            channelDisplayName,
            onExited,
        } = this.props;

        return (
            <Modal
                dialogClassName='a11y__modal'
                show={this.state.show}
                onHide={this.onHide}
                onExited={onExited}
                role='dialog'
                aria-labelledby='convertChannelModalLabel'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title
                        componentClass='h1'
                        id='convertChannelModalLabel'
                    >
                        <FormattedMessage
                            id='convert_channel.title'
                            defaultMessage='Confirm changing {display_name} to a private channel?'
                            values={{
                                display_name: channelDisplayName,
                            }}
                        />
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>
                        <FormattedMarkdownMessage
                            id='convert_channel.question1'
                            defaultMessage='When you convert **{display_name}** to a private channel, history and membership are preserved. Publicly shared files remain accessible to anyone with the link. Membership in a private channel is by invitation only.'
                        />
                    </p>
                    {/* todo i18n */}
                    <p style={{ marginTop: '25px' }}>
                        <div className='Input_wrapper'>
                            <input
                                className='Input form-control medium new-channel-modal-name-input channel-name-input-field'
                                placeholder='Enter channel name'
                                onChange={this.onUpdateConfirmName}
                                autoFocus
                            />
                        </div>
                    </p>
                    {/* todo i18n */}
                    <p style={{ marginTop: '25px' }}>
                        <div className='Input_wrapper'>
                            <input
                                className='Input form-control medium new-channel-modal-name-input channel-name-input-field'
                                placeholder='輸入頻道名稱'
                                onChange={this.onUpdateConfirmName}
                                autoFocus
                            />
                        </div>
                    </p>
                    <p style={{ fontSize: '12px', color: 'rgba(63, 67 89, 0.75)' }}>
                        請輸入此頻道名稱已確認變更
                    </p>
                </Modal.Body>
                <Modal.Footer>
                    <button
                        type='button'
                        className='btn btn-tertiary'
                        onClick={this.onHide}
                    >
                        <FormattedMessage
                            id='convert_channel.cancel'
                            defaultMessage='No, cancel'
                        />
                    </button>
                    <button
                        type='button'
                        className='btn btn-primary'
                        data-dismiss='modal'
                        onClick={this.handleConvert}
                        autoFocus={true}
                        disabled={!this.canConvert()}
                        data-testid='convertChannelConfirm'
                    >
                        <FormattedMessage
                            id='convert_channel.confirm'
                            defaultMessage='Yes, convert to private channel'
                        />
                    </button>
                </Modal.Footer>
            </Modal>
        );
    }
}
