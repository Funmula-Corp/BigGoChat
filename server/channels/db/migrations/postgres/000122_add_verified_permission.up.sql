ALTER TABLE channelmembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultteamverifiedrole character varying(64);
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultchannelverifiedrole character varying(64);
ALTER TABLE teammembers ADD COLUMN IF NOT EXISTS schemeverified boolean;

INSERT INTO public.roles (id, name, displayname, description, createat, updateat, deleteat, permissions, schememanaged, builtin) VALUES
('biggoyyyyyyyyyyyyyyyyyyyyn', 'system_verified', 'authentication.roles.system_verified.name', 'authentication.roles.system_verified.description', 1716905618000, 1716905618000, 0, 'create_custom_group delete_emojis manage_custom_group_members create_direct_channel delete_custom_group join_public_teams create_group_channel create_emojis edit_custom_group view_members list_public_teams restore_custom_group create_team', true, true),
('biggoyyyyyyyyyyyyyyyyyyyyd', 'team_verified', 'authentication.roles.team_verified.name', 'authentication.roles.team_verified.description', 155291281668000, 1552912816680, 0, ' list_team_channels join_public_channels read_public_channel view_team create_public_channel manage_public_channel_properties delete_public_channel create_private_channel manage_private_channel_properties delete_private_channel invite_user add_user_to_team', true, true),
('biggoyyyyyyyyyyyyyyyyyyyyr', 'channel_verified', 'authentication.roles.channel_verified.name', 'authentication.roles.channel_verified.description', 1552912816680, 1552912816680, 0, ' read_channel add_reaction remove_reaction manage_public_channel_members upload_file get_public_link create_post use_slash_commands manage_private_channel_members delete_post edit_post', true, true);
