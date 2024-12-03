// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import type {GlobalState} from 'types/store';

import DisplayName from './display_name';
import { haveIVerified } from 'mattermost-redux/selectors/entities/roles_helpers';

function mapStateToProps(state: GlobalState) {
    const isPhoneVerified = haveIVerified(state);

    return {
        isPhoneVerified,
    }
}

export default connect(mapStateToProps)(DisplayName);
