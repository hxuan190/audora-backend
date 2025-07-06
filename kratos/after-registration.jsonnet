{
  identity: {
    id: payload.identity.id,
    traits: payload.identity.traits,
    created_at: payload.identity.created_at,
    updated_at: payload.identity.updated_at
  },
  flow: {
    id: payload.flow.id,
    type: payload.flow.type,
    expires_at: payload.flow.expires_at
  }
} 