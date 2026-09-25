"use client";

import { createContext, useContext, useState } from "react";
import { useRouter } from "next/navigation";
import { apiRequest } from "@/utils/api";
import { useToast } from "@/app/context/ToastContext";
import EditProfileModal from "./EditProfileModal";
import UsersModal from "./UsersModal";

const ProfileModalsContext = createContext(null);

export function ProfileModalsProvider({ id, isPrivateHidden, children }) {
  const router = useRouter();
  const showToast = useToast();

  // Modal states
  const [modalType, setModalType] = useState(null); // "followers" or "following"
  const [usersList, setUsersList] = useState([]);
  const [isEditing, setIsEditing] = useState(false);

  const reload = () => router.refresh();

  // Open modal function
  const openModal = async (type) => {
    if (isPrivateHidden) return; // Stop if private

    setModalType(type); // Show modal
    try {
      const data = await apiRequest(`/profile/${id}/${type}`);
      setUsersList(data || []); // Save users
    } catch (err) {
      setModalType(null); // Close on error
      showToast(err.message);
    }
  };

  return (
    <ProfileModalsContext.Provider
      value={{
        isEditing,
        modalType,
        usersList,
        reload,
        openModal,
        openEditProfile: () => setIsEditing(true),
        closeEditProfile: () => setIsEditing(false),
        closeUsersModal: () => setModalType(null),
        onProfileSaved: async () => {
          setIsEditing(false);
          await reload();
          showToast("Profile updated", false);
        },
      }}
    >
      {children}
    </ProfileModalsContext.Provider>
  );
}

export function useProfileModals() {
  return useContext(ProfileModalsContext);
}

export function ProfileModals({ user }) {
  const {
    isEditing,
    modalType,
    usersList,
    closeEditProfile,
    onProfileSaved,
    closeUsersModal,
  } = useProfileModals();

  return (
    <>
      {/* Edit Profile */}
      {isEditing && (
        <EditProfileModal
          user={user}
          onClose={closeEditProfile}
          onSaved={onProfileSaved}
        />
      )}

      {/* Users List */}
      {modalType && (
        <UsersModal
          type={modalType}
          users={usersList}
          onClose={closeUsersModal}
        />
      )}
    </>
  );
}
