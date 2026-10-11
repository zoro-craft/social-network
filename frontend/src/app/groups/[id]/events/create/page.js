'use client';

import { useToast } from '@/app/context/ToastContext';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';
import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';

export default function CreatEvent() {
  const { auth, loading: authLoading } = useRequireAuth();
  const router = useRouter();
  const showToast = useToast();
  const { id } = useParams();
  const groupId = Number(id);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [eventDate, setEventDate] = useState('');
  const [eventTime, setEventTime] = useState('');
  const [fitsOptions, setFitsOptions] = useState('');
  const [secondOptions, setSecondOptions] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!title.trim()) {
      showToast('Title is required');
      return;
    }
    if (!eventDate) {
      showToast('Event date is required');
      return;
    }
    if (!eventTime) {
      showToast('Event time is required');
      return;
    }
    if (!fitsOptions.trim() || !secondOptions.trim()) {
      showToast('Both participation options are required');
      return;
    }
    if (fitsOptions.trim() === secondOptions.trim()) {
      showToast('The two options must be different');
      return;
    }

    setLoading(true);

    try {
      const res = await apiRequest('/events', 'POST', {
        group_id: groupId,
        title: title.trim(),
        description: description.trim(),
        eventDate,
        eventTime,
        fitsOptions: fitsOptions.trim(),
        secondOptions: secondOptions.trim(),
      });

      showToast(res.message || 'Event created successfully!', false);
      router.push(`/groups/${id}`);
    } catch (err) {
      showToast(err.message || 'Failed to create event');
      setLoading(false);
    }
  };

  if (authLoading || !auth.authenticated) return null;

  return (
    <div className="glass-page">
      <div className="glass-container glass-container-lg">
        <h2>Create Event</h2>
        <p className="glass-subtitle">Organize an activity for your group</p>

        <form onSubmit={handleSubmit} className="glass-form">
          <div className="glass-input-group">
            <input
              className="glass-input"
              type="text"
              placeholder="Event title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <textarea
              className="glass-input"
              placeholder="Event description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>

          <div className="glass-input-group">
            <label>Event date</label>
            <input
              className="glass-input"
              type="date"
              value={eventDate}
              onChange={(e) => setEventDate(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <label>Event time</label>
            <input
              className="glass-input"
              type="time"
              value={eventTime}
              onChange={(e) => setEventTime(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <label>Participation option 1</label>
            <input
              className="glass-input"
              type="text"
              placeholder="Going"
              value={fitsOptions}
              onChange={(e) => setFitsOptions(e.target.value)}
              required
            />
          </div>

          <div className="glass-input-group">
            <label>Participation option 2</label>
            <input
              className="glass-input"
              type="text"
              placeholder="Not going"
              value={secondOptions}
              onChange={(e) => setSecondOptions(e.target.value)}
              required
            />
          </div>

          <button
            className="glass-button"
            type="submit"
            disabled={loading}
          >
            {loading ? 'Creating...' : 'Create Event'}
          </button>
        </form>
      </div>
    </div>
  );
}
