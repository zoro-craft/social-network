'use client';

import { useEffect, useState } from 'react';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import { useToast } from '@/app/context/ToastContext';

export default function GroupEvent({ group_id }) {
  const showToast = useToast();
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshKey, setRefreshKey] = useState(0);
  const [votingId, setVotingId] = useState(null);

  useEffect(() => {
    if (!group_id) return;

    const fetchEvents = async () => {
      try {
        const res = await fetch(`${API_BASE}/events?group_id=${group_id}`, {
          credentials: 'include',
        });
        if (res.ok) {
          const data = await res.json();
          setEvents(data || []);
        } else {
          setEvents([]);
        }
      } catch {
        showToast('Failed to fetch events');
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, [group_id, refreshKey]);

  const handleVote = async (eventId, option) => {
    setVotingId(eventId);
    try {
      const res = await fetch(`${API_BASE}/events/vote`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ event_id: eventId, vote: option }),
      });
      if (res.ok) {
        setRefreshKey((k) => k + 1);
      } else {
        const err = await res.json();
        showToast(err.error || 'Failed to save vote');
      }
    } catch {
      showToast('Failed to save vote');
    } finally {
      setVotingId(null);
    }
  };

  const countFor = (event, option) =>
    (event.vote_counts || []).find((vc) => vc.option === option)?.count || 0;

  const totalVotes = (event) =>
    (event.vote_counts || []).reduce((sum, vc) => sum + vc.count, 0);

  const percentOf = (count, event) => {
    const total = totalVotes(event);
    if (total === 0) return 0;
    return Math.round((count / total) * 100);
  };

  if (loading) {
    return <p className="glass-muted-text">Loading events...</p>;
  }

  if (events.length === 0) {
    return (
      <p className="glass-muted-text">No events yet. Create the first one!</p>
    );
  }

  return (
    <div>
      {events.map((event) => {
        const total = totalVotes(event);
        return (
          <div
            key={event.id}
            className="glass-post-card"
            style={{ marginBottom: 16 }}
          >
            <div className="glass-post-header">
              <span className="glass-post-author">{event.title}</span>
              <span>{formatDate(event.created_at)}</span>
            </div>

            <div className="glass-vote-options">
              {(event.options || []).map((option) => {
                const count = countFor(event, option);
                const pct = percentOf(count, event);
                const isSelected = event.user_response === option;
                return (
                  <button
                    key={option}
                    className={`glass-vote-option${
                      isSelected ? ' glass-vote-selected' : ''
                    }`}
                    onClick={() => handleVote(event.id, option)}
                    disabled={votingId === event.id}
                  >
                    <span className="glass-vote-label">{option}</span>
                    <span className="glass-vote-bar">
                      <span
                        className="glass-vote-fill"
                        style={{ width: `${pct}%` }}
                      />
                    </span>
                    <span className="glass-vote-meta">
                      {count} vote{count === 1 ? '' : 's'}
                      {pct > 0 ? ` · ${pct}%` : ''}
                      {isSelected ? ' · Your vote' : ''}
                    </span>
                  </button>
                );
              })}
            </div>

            <p className="glass-vote-total">
              {total} total vote{total === 1 ? '' : 's'}
            </p>
          </div>
        );
      })}
    </div>
  );
}
