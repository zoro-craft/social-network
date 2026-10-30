import { ChatMessage } from '@/lib/types';

// Returns the id the message targets for the viewer: for a DM it is the other
// participant (the sender for the recipient, the recipient for the sender);
// for a group message it is the group id.
export function otherParticipant(payload: ChatMessage, me: number): number {
  if (payload.type === 'private') {
    return payload.sender_id === me ? payload.ref_id : payload.sender_id;
  }
  return payload.ref_id;
}

export function conversationKey(
  type: 'private' | 'group',
  id: number,
): string {
  return `${type}:${id}`;
}

export function conversationKeyFromPayload(payload: ChatMessage, me: number): string {
  return `private:${otherParticipant(payload, me)}`;
}