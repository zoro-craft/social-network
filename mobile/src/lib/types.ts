// API payload shapes, mirrored from the Go backend models.

export interface User {
  id: number;
  first_name: string;
  last_name: string;
  nickname?: string;
  email: string;
  avatar?: string;
  about_me?: string;
  date_of_birth?: string;
  is_private?: boolean;
  created_at?: string;
  full_name?: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  image?: string | null;
  privacy?: string;
  group_id?: number | null;
  user_id: number;
  full_name: string;
  created_at: string;
  like_count: number;
  dislike_count: number;
  user_reaction?: string | null;
}

export interface Comment {
  id: number;
  post_id: number;
  group_id?: number | null;
  user_id: number;
  content: string;
  image?: string;
  first_name?: string;
  avatar?: string;
  created_at?: string;
}

// Stories are ephemeral (24h) image posts, grouped per author like Instagram's rail.
export interface Story {
  id: number;
  user_id: number;
  image: string;
  caption: string;
  created_at: number;
  seen: boolean;
}

export interface StoryGroup {
  user_id: number;
  full_name: string;
  avatar: string;
  seen: boolean;
  stories: Story[];
}

// Reels are short videos in a persistent vertical feed (unlike 24h stories), so
// every engagement counter is stored alongside them.
export interface Reel {
  id: number;
  user_id: number;
  full_name: string;
  avatar: string;
  video: string;
  caption: string;
  created_at: number;
  like_count: number;
  comment_count: number;
  view_count: number;
  save_count: number;
  share_count: number;
  liked: boolean;
  saved: boolean;
}

// Returned by every reel engagement endpoint so one reel can be patched in place.
export interface ReelStats {
  reel_id: number;
  like_count: number;
  comment_count: number;
  view_count: number;
  save_count: number;
  share_count: number;
  liked: boolean;
  saved: boolean;
}

export interface ReelComment {
  id: number;
  reel_id: number;
  user_id: number;
  content: string;
  first_name?: string;
  avatar?: string;
  created_at: number;
}

export interface ProfileData {
  user: User;
  relation: 'self' | 'accepted' | 'pending' | 'none';
  followers_count: number;
  following_count: number;
  posts: Post[];
}

export interface Notification {
  id: number;
  user_id: number;
  actor_id?: number | null;
  type: string;
  ref_id?: number | null;
  message: string;
  is_read: boolean;
  created_at: number;
}

export interface NotificationsResponse {
  items: Notification[];
  unread_count: number;
  has_more: boolean;
}

export interface Group {
  id: number;
  titel: string;
  description: string;
  image?: string;
  is_member?: boolean;
  is_pending?: boolean;
  member_count?: number;
  created_at?: string;
}

export interface Event {
  id: number;
  group_id: number;
  title: string;
  description: string;
  event_date: string;
  event_time: string;
  created_by: number;
  created_at: string;
  options: string[];
  vote_counts: { option: string; count: number }[];
  user_response?: string;
}

export interface ConversationUser {
  user_id: number;
  name: string;
  avatar?: string;
  online: boolean;
  last_message?: string | null;
  last_message_at?: number | null;
  can_chat: boolean;
}

export interface ConversationGroup {
  group_id: number;
  title: string;
  last_message?: string | null;
  last_message_at?: number | null;
}

export interface ConversationsResponse {
  users: ConversationUser[];
  groups: ConversationGroup[];
}

export interface ChatMessage {
  id: number;
  sender_id: number;
  sender_name: string;
  content: string;
  created_at: number;
  type: 'private' | 'group';
  ref_id: number;
}

export interface UserSearchResult {
  id: number;
  full_name: string;
  first_name?: string;
  last_name?: string;
  nickname?: string;
  avatar?: string;
}