// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import type {Dispatch} from 'redux';

import type {Channel} from '@mattermost/types/channels';

import {getChannels, getArchivedChannels, joinChannel, getChannelsMemberCount, searchAllChannels, getAllPrivateChannels} from 'mattermost-redux/actions/channels';
import {RequestStatus} from 'mattermost-redux/constants';
import {createSelector} from 'mattermost-redux/selectors/create_selector';
import {getChannelsInCurrentTeam, getMyChannelMemberships, getChannelsMemberCount as getChannelsMemberCountSelector} from 'mattermost-redux/selectors/entities/channels';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getCurrentTeam, getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';
import {Permissions} from 'mattermost-redux/constants';

import {setGlobalItem} from 'actions/storage';
import {openModal, closeModal} from 'actions/views/modals';
import {closeRightHandSide} from 'actions/views/rhs';
import {getIsRhsOpen, getRhsState} from 'selectors/rhs';
import {makeGetGlobalItem} from 'selectors/storage';

import Constants, {StoragePrefixes} from 'utils/constants';

import type {GlobalState} from 'types/store';

import BrowseChannels from './browse_channels';
import { haveICurrentTeamPermission } from 'mattermost-redux/selectors/entities/roles';

const getChannelsWithoutPrivateArchived = createSelector(
    'getChannelsWithoutArchived',
    getChannelsInCurrentTeam,
    (channels: Channel[]) => channels && channels.filter((c) => c.type !== Constants.PRIVATE_CHANNEL),
);

const getArchivedOtherChannels = createSelector(
    'getArchivedOtherChannels',
    getChannelsInCurrentTeam,
    (channels: Channel[]) => channels && channels.filter((c) => c.delete_at !== 0),
);

const getPrivateChannelsSelector = createSelector(
    'getPrivateChannelsSelector',
    getChannelsInCurrentTeam,
    (channels: Channel[]) => channels && channels.filter((c) => c.type === Constants.PRIVATE_CHANNEL),
);

function mapStateToProps(state: GlobalState) {
    const team = getCurrentTeam(state);
    const getGlobalItem = makeGetGlobalItem(StoragePrefixes.HIDE_JOINED_CHANNELS, 'false');

    return {
        channels: getChannelsWithoutPrivateArchived(state) || [],
        archivedChannels: getArchivedOtherChannels(state) || [],
        privateChannels: getPrivateChannelsSelector(state) || [],
        currentUserId: getCurrentUserId(state),
        teamId: getCurrentTeamId(state),
        teamName: team?.name,
        channelsRequestStarted: state.requests.channels.getChannels.status === RequestStatus.STARTED,
        canShowArchivedChannels: (getConfig(state).ExperimentalViewArchivedChannels === 'true'),
        canShowAllPrivateChannels: haveICurrentTeamPermission(state, Permissions.MANAGE_TEAM),
        myChannelMemberships: getMyChannelMemberships(state) || {},
        shouldHideJoinedChannels: getGlobalItem(state) === 'true',
        rhsState: getRhsState(state),
        rhsOpen: getIsRhsOpen(state),
        channelsMemberCount: getChannelsMemberCountSelector(state),
    };
}

function mapDispatchToProps(dispatch: Dispatch) {
    return {
        actions: bindActionCreators({
            getChannels,
            getAllPrivateChannels,
            getArchivedChannels,
            joinChannel,
            searchAllChannels,
            openModal,
            closeModal,
            setGlobalItem,
            closeRightHandSide,
            getChannelsMemberCount,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(BrowseChannels);
