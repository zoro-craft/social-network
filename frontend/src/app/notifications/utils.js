import { apiRequest } from "@/utils/api";

export const handleAcceptInvitation = async (notification,) => {
    try {
        await apiRequest(`/invitations/${notification.invitation_id}/accept`, 'POST');
        return true
    } catch (error) {
        return false
        console.error('Failed to accept invitation:', error);
    }
};

export const handleRejectInvitation = async (notification) => {
    try {
        await apiRequest(`/invitations/${notification.invitation_id}/reject`, 'POST');
        return true
    } catch (error) {
        console.error('Failed to reject invitation:', error);
        return false

    }
};

export const handleAcceptRequest = async (notification) => {
    try {
        await apiRequest(`/join_requests/${notification.request_id}/accept`, 'POST');
        return true
    } catch (error) {
        console.error('Failed to accept join request:', error);
        return false
    }
};

export const handleRejectRequest = async (notification) => {
    try {
        await apiRequest(`/join_requests/${notification.request_id}/reject`, 'POST');
        return true;

    } catch (error) {
        console.error('Failed to reject join request:', error);
        return false
    }
};

