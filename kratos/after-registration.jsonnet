local identity = std.extVar('identity');
local flow = std.extVar('flow');
local request_headers = std.extVar('request_headers');
local request_method = std.extVar('request_method');
local request_url = std.extVar('request_url');

{
  // User identity data
  identity: {
    id: identity.id,
    traits: {
      email: identity.traits.email,
      display_name: identity.traits.display_name,
      first_name: if std.objectHas(identity.traits, 'first_name') then identity.traits.first_name else null,
      last_name: if std.objectHas(identity.traits, 'last_name') then identity.traits.last_name else null,
      user_type: if std.objectHas(identity.traits, 'user_type') then identity.traits.user_type else 'listener',
      artist_name: if std.objectHas(identity.traits, 'artist_name') then identity.traits.artist_name else null,
      bio: if std.objectHas(identity.traits, 'bio') then identity.traits.bio else null,
      location: if std.objectHas(identity.traits, 'location') then identity.traits.location else null,
      profile_image: if std.objectHas(identity.traits, 'profile_image') then identity.traits.profile_image else null,
      preferences: if std.objectHas(identity.traits, 'preferences') then identity.traits.preferences else {
        email_notifications: true,
        marketing_emails: false
      }
    },
    schema_id: identity.schema_id,
    state: identity.state,
    created_at: identity.created_at,
    updated_at: identity.updated_at
  },
  
  // Flow information
  flow: {
    id: flow.id,
    type: flow.type,
    expires_at: flow.expires_at,
    issued_at: flow.issued_at,
    request_url: flow.request_url
  },
  
  // Request context for logging/analytics
  request_context: {
    method: request_method,
    url: request_url,
    user_agent: if std.objectHas(request_headers, 'User-Agent') then request_headers['User-Agent'] else null,
    ip_address: if std.objectHas(request_headers, 'X-Forwarded-For') then request_headers['X-Forwarded-For'] else if std.objectHas(request_headers, 'X-Real-IP') then request_headers['X-Real-IP'] else null,
    timestamp: std.toString(std.floor(std.time))
  }
}