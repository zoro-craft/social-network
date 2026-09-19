'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import { useToast } from '@/app/context/ToastContext';

export default function GroupPosts({ groupId }) {
    const router = useRouter();
    const showToast = useToast();
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!groupId) return;

        const fetchPosts = async () => {
            try {
                const res = await fetch(`${API_BASE}/groups/${groupId}/posts`, {
                    credentials: 'include',
                });
                if (res.ok) {
                    const data = await res.json();
                    setPosts(data || []);
                } else {
                    setPosts([]);
                }
            } catch {
                showToast('Failed to fetch posts');
                setPosts([]);
            } finally {
                setLoading(false);
            }
        };

        fetchPosts();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [groupId]);

    if (loading) {
        return <p className="glass-muted-text">Loading posts...</p>;
    }

    if (posts.length === 0) {
        return (
            <p className="glass-muted-text">No posts yet. Create the first one!</p>
        );
    }

    return (
        <div className="glass-group-list">
            {posts.map((post) => (
                <div key={post.id} className="glass-post-card">
                    <div className="glass-post-header">
                        <span className="glass-post-author">Posted by {post.full_name}</span>
                        <span>{formatDate(post.created_at)}</span>
                    </div>
                    <div className="glass-post-title">{post.title}</div>

                    {post.image && (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                            className="glass-post-image"
                            src={`${API_BASE}${post.image}`}
                            alt=""
                        />
                    )}

                    <div className="glass-post-content">{post.content}</div>

                    <div className="glass-reactions">
                        <button
                            className="glass-react-btn"
                            onClick={() => router.push(`/posts/comment/${post.id}`)}
                        >
                            💬 Comment
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
}
