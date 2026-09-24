import { apiRequest } from "@/utils/api";

// Accept an invitation using the notification's ref_id (migration 000024 renamed invitation_id).
export const handleAcceptInvitation = async (notification) => {
  try {
    await apiRequest(`/invitations/${notification.ref_id}/accept`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

export const handleRejectInvitation = async (notification) => {
  try {
    await apiRequest(`/invitations/${notification.ref_id}/reject`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

// Accept a join request using the notification's ref_id (migration 000024 renamed request_id).
export const handleAcceptRequest = async (notification) => {
  try {
    await apiRequest(`/join_requests/${notification.ref_id}/accept`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

export const handleRejectRequest = async (notification) => {
  try {
    await apiRequest(`/join_requests/${notification.ref_id}/reject`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

// A follow request is resolved through the followers endpoint, keyed by the acting user id.
export const handleAcceptFollow = async (notification) => {
  try {
    await apiRequest(`/follow/${notification.ref_id}/accept`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

export const handleDeclineFollow = async (notification) => {
  try {
    await apiRequest(`/follow/${notification.ref_id}/decline`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};

// Mark an informational event notification read.
export const handleMarkRead = async (notification) => {
  try {
    await apiRequest(`/notifications/${notification.id}/read`, "POST");
    return { ok: true };
  } catch (error) {
    return { ok: false, status: error.status };
  }
};