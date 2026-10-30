import AppButton, { ButtonVariant } from '@/components/AppButton';

interface Props {
  relation: 'self' | 'accepted' | 'pending' | 'none';
  isPrivate?: boolean;
  onToggle: () => void;
  disabled?: boolean;
}

// Mirrors FollowButton.jsx relation-aware states.
const CONFIG: Record<
  string,
  { text: string; variant: ButtonVariant }
> = {
  self: { text: 'Edit Profile', variant: 'ghost' },
  accepted: { text: 'Unfollow', variant: 'danger' },
  none: { text: 'Follow', variant: 'primary' },
};

export default function FollowButton({ relation, isPrivate, onToggle, disabled }: Props) {
  const config =
    relation === 'pending'
      ? { text: isPrivate ? 'Cancel Request' : 'Following', variant: 'pending' as ButtonVariant }
      : CONFIG[relation] || CONFIG.none;

  return (
    <AppButton
      title={config.text}
      onPress={onToggle}
      variant={config.variant}
      disabled={disabled}
    />
  );
}