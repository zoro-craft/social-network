// otherParticipant returns the id the payload targets for the viewer: for a DM it is the
// participant who is not the viewer (the sender for the recipient, the recipient for the
// sender); for a group message it is the group id.
export function otherParticipant(payload, me) {
  return payload.type === 'private'
    ? payload.sender_id === me
      ? payload.ref_id
      : payload.sender_id
    : payload.ref_id;
}

// conversationKeyFromPayload derives the per-viewer conversation key for a chat envelope payload.
export function conversationKeyFromPayload(payload, me) {
  return `${payload.type}:${otherParticipant(payload, me)}`;
}