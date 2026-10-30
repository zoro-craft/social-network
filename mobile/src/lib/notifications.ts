import { apiRequest } from '@/lib/api';
import { Notification } from '@/lib/types';

type ActionResult = { ok: boolean; status?: number };

async function run(fn: () => Promise<unknown>): Promise<ActionResult> {
  try {
    await fn();
    return { ok: true };
  } catch (error: any) {
    return { ok: false, status: error?.status };
  }
}

export const handleAcceptInvitation = (n: Notification) =>
  run(() => apiRequest(`/invitations/${n.ref_id}/accept`, { method: 'POST' }));

export const handleRejectInvitation = (n: Notification) =>
  run(() => apiRequest(`/invitations/${n.ref_id}/reject`, { method: 'POST' }));

export const handleAcceptRequest = (n: Notification) =>
  run(() => apiRequest(`/join_requests/${n.ref_id}/accept`, { method: 'POST' }));

export const handleRejectRequest = (n: Notification) =>
  run(() => apiRequest(`/join_requests/${n.ref_id}/reject`, { method: 'POST' }));

export const handleAcceptFollow = (n: Notification) =>
  run(() => apiRequest(`/follow/${n.ref_id}/accept`, { method: 'POST' }));

export const handleDeclineFollow = (n: Notification) =>
  run(() => apiRequest(`/follow/${n.ref_id}/decline`, { method: 'POST' }));

export const handleMarkRead = (n: Notification) =>
  run(() => apiRequest(`/notifications/${n.id}/read`, { method: 'POST' }));