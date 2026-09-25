import Link from "next/link";
import Avatar from "@/components/Avatar";
import { getProfile } from "./getProfile";
import { ProfileModals, ProfileModalsProvider } from "./ProfileModals";
import FollowButton from "./FollowButton";
import StatItem from "./StatItem";
import PostGrid from "./PostGrid";

export default async function ProfilePage({ params }) {
  const { id } = await params;

  // Fetch data
  const { profile, error } = await getProfile(id);

  // Loading & Errors
  if (error) {
    return (
      <div className="profile-wrapper">
        <div className="glass-card">
          <h2>Profile unavailable</h2>
          <p>{error}</p>
          <Link href="/posts">Back to posts</Link>
        </div>
      </div>
    );
  }

  if (!profile) {
    return <div className="profile-wrapper">Loading...</div>;
  }

  // Extract data
  const { user, relation, followers_count, following_count, posts } = profile;

  // Check if profile is private
  const isPrivateHidden = user.is_private && relation !== "self" && relation !== "accepted";

  // Render UI
  return (
    <ProfileModalsProvider id={id} isPrivateHidden={isPrivateHidden}>
      <div className="profile-wrapper">
        <div className="profile-layout">

          {/* Left Sidebar: User Info */}
          <div className="glass-card left-sidebar">
            <Avatar src={user.avatar} name={user.first_name} size={112} />
            <h2>{user.first_name} {user.last_name}</h2>
            <p>@{user.nickname || user.first_name}</p>
            <p>{user.is_private ? "Private Profile" : "Public Profile"}</p>

            {/* Show if not private */}
            {!isPrivateHidden && (
              <>
                <p>Born: {user.date_of_birth?.split("T")[0]}</p>
                <p>{user.about_me}</p>
              </>
            )}

            <FollowButton
              id={id}
              relation={relation}
              isPrivate={user.is_private}
            />
          </div>

          {/* Right Area: Stats & Posts */}
          <div className="right-area">

            {/* Stats */}
            <div className="stats-grid">
              <StatItem type="followers" label="Followers" value={followers_count || 0} />

              <StatItem type="following" label="Following" value={following_count || 0} />

              <div className="glass-card stat-item">
                <p>Joined</p>
                <p>{user.created_at ? user.created_at.split("-")[0] : "Recently"}</p>
              </div>
            </div>

            {/* Posts */}
            <div className="glass-card">
              <h3>Posts</h3>
              <PostGrid posts={posts} isPrivateHidden={isPrivateHidden} />
            </div>

          </div>
        </div>

        {/* Modals */}
        <ProfileModals user={user} />
      </div>
    </ProfileModalsProvider>
  );
}
