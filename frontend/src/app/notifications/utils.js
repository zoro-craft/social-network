export const handleAcceptInvitation = async (notification) => {
    try {
        await apiRequest(`/invitations/${notification.invitation_id}/accept`, 'POST',);

        removeNotification(notification.id);
    } catch (error) {
        console.error('Failed to accept invitation:', error);
    }
};

export const handleRejectInvitation = async (notification) => {
    try {
        await apiRequest(
            `/invitations/${notification.invitation_id}/reject`,
            {
                method: 'POST',
            }
        );

        removeNotification(notification.id);
    } catch (error) {
        console.error('Failed to reject invitation:', error);
    }
};

export const handleAcceptRequest = async (notification) => {
    try {
        await apiRequest(
            `/join_requests/${notification.request_id}/accept`,
            {
                method: 'POST',
            }
        );

        removeNotification(notification.id);
    } catch (error) {
        console.error('Failed to accept join request:', error);
    }
};

export const handleRejectRequest = async (notification) => {
    try {
        await apiRequest(
            `/join_requests/${notification.request_id}/reject`,
            {
                method: 'POST',
            }
        );

        removeNotification(notification.id);
    } catch (error) {
        console.error('Failed to reject join request:', error);
    }
};

// Remove notification from UI
const removeNotification = (notificationId) => {
    setNotifications((prev) =>
        prev.filter((notification) => notification.id !== notificationId)
    );
};
