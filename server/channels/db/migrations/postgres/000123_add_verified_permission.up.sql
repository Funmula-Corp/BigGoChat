ALTER TABLE channelmembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultteamverifiedrole character varying(64);
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultchannelverifiedrole character varying(64);
ALTER TABLE teammembers ADD COLUMN IF NOT EXISTS schemeverified boolean;

INSERT INTO public.roles (id, name, displayname, description, createat, updateat, deleteat, permissions, schememanaged, builtin) VALUES
('biggoryyyyyyyyyyyyyyyyyyyd', 'team_verified', 'authentication.roles.team_verified.name', 'authentication.roles.team_verified.description', 155291281668000, 1552912816680, 0, ' list_team_channels join_public_channels read_public_channel view_team create_public_channel manage_public_channel_properties delete_public_channel create_private_channel manage_private_channel_properties delete_private_channel invite_user add_user_to_team', true, true),
('biggoryyyyyyyyyyyyyyyyyyyr', 'channel_verified', 'authentication.roles.channel_verified.name', 'authentication.roles.channel_verified.description', 1552912816680, 1552912816680, 0, ' read_channel add_reaction remove_reaction manage_public_channel_members upload_file get_public_link create_post use_slash_commands manage_private_channel_members delete_post edit_post', true, true);
