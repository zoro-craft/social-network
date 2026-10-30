import { useCallback, useEffect, useRef, useState } from 'react';
import {
  FlatList,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import HeaderBar from '@/components/HeaderBar';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useRealtime } from '@/context/RealtimeContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { conversationKey } from '@/lib/chat';
import { EMOJI } from '@/lib/emoji';
import { ChatMessage } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

const BATCH = 10;
const PAGING_COOLDOWN = 2000;

export default function ChatThreadScreen() {
  const params = useLocalSearchParams<{ type?: string; id?: string; name?: string }>();
  const router = useRouter();
  const { auth } = useAuth();
  const { subscribe } = useRealtime();
  const showToast = useToast();

  const type = params.type === 'group' ? 'group' : 'private';
  const id = Number(params.id);
  const name = params.name ? decodeURIComponent(params.name) : 'Chat';

  const [messages, setMessages] = useState<ChatMessage[]>([]); // newest-first
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(true);
  const [canChat, setCanChat] = useState(true);
  const [input, setInput] = useState('');
  const [showEmoji, setShowEmoji] = useState(false);
  const [sending, setSending] = useState(false);

  const loadingMoreRef = useRef(false);
  const lastFetchRef = useRef(0);
  const listRef = useRef<FlatList<ChatMessage>>(null);
  const selectedRef = useRef({ type, id });
  selectedRef.current = { type, id };

  const historyPath = type === 'group' ? `/groups/${id}/messages` : `/messages/${id}`;
  const key = conversationKey(type, id);

  const loadHistory = useCallback(
    async (beforeId?: number, isMore = false) => {
      const query = `?limit=${BATCH}${beforeId ? `&before=${beforeId}` : ''}`;
      try {
        const data = await apiRequest<{ messages: ChatMessage[]; has_more: boolean }>(`${historyPath}${query}`);
        // Server returns newest-first; keep that order for the inverted list.
        const batch = data.messages || [];
        if (isMore) {
          setMessages((prev) => {
            const seen = new Set(prev.map((m) => m.id));
            return [...prev, ...batch.filter((m) => !seen.has(m.id))];
          });
        } else {
          setMessages(batch);
        }
        setHasMore(!!data.has_more);
      } catch (e: any) {
        if (e?.status === 400) {
          setCanChat(false);
        }
        setHasMore(false);
      } finally {
        setLoading(false);
      }
    },
    [historyPath],
  );

  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  // Load older messages when the user stays at the top.
  const loadOlder = useCallback(() => {
    const oldest = messages[messages.length - 1];
    if (!oldest) return;
    if (loadingMoreRef.current || !hasMore) return;
    if (PAGING_COOLDOWN - (Date.now() - lastFetchRef.current) > 0) return;
    loadingMoreRef.current = true;
    lastFetchRef.current = Date.now();
    loadHistory(oldest.id, true).finally(() => {
      loadingMoreRef.current = false;
    });
  }, [hasMore, messages, loadHistory]);

  // Live arrivals append to the open thread.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onMessage = (payload: ChatMessage) => {
      const convKey =
        payload.type === 'private'
          ? `private:${payload.sender_id === auth.user?.id ? payload.ref_id : payload.sender_id}`
          : `group:${payload.ref_id}`;
      if (convKey !== key) return;
      setMessages((prev) => {
        if (prev.some((m) => m.id === payload.id)) return prev;
        return [payload, ...prev];
      });
    };
    const unsub1 = subscribe('chat_new_message', onMessage);
    const unsub2 = subscribe('chat_new_group_message', onMessage);
    return () => {
      unsub1();
      unsub2();
    };
  }, [auth.authenticated, auth.user?.id, key, subscribe]);

  const send = async () => {
    const content = input.trim();
    if (!content) return;
    const sel = selectedRef.current;
    if (!sel) return;

    setSending(true);
    const body = sel.type === 'group' ? { content } : { recipient_id: sel.id, content };
    const path = sel.type === 'group' ? `/groups/${sel.id}/messages` : `/messages`;
    try {
      await apiRequest(path, { method: 'POST', data: body });
      setInput('');
      setShowEmoji(false);
      // Jump to the message we just sent, in case the user had scrolled up.
      listRef.current?.scrollToOffset({ offset: 0, animated: true });
    } catch (e: any) {
      if (e?.status === 400) {
        setCanChat(false);
      } else {
        showToast('Could not send message');
      }
    } finally {
      setSending(false);
    }
  };

  const insertEmoji = (emoji: string) => setInput((prev) => `${prev}${emoji}`);

  if (loading) {
    return (
      <View style={styles.screen}>
        <HeaderBar title={name} onBack={() => router.back()} />
        <Loading />
      </View>
    );
  }

  return (
    <KeyboardAvoidingView
      style={styles.screen}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 0 : 0}
    >
      <HeaderBar title={name} onBack={() => router.back()} />

      {!canChat && (
        <View style={styles.frozen}>
          <Text style={styles.frozenText}>
            You can't chat with this person yet (they have restricted messages).
          </Text>
        </View>
      )}

      <FlatList
        ref={listRef}
        inverted
        data={messages}
        keyExtractor={(item) => String(item.id)}
        contentContainerStyle={styles.list}
        onEndReached={loadOlder}
        onEndReachedThreshold={0.3}
        ListEmptyComponent={<EmptyState message="No messages yet. Say hello!" />}
        renderItem={({ item }) => {
          const own = item.sender_id === auth.user?.id;
          return <Bubble message={item} own={own} showName={type === 'group' && !own} />;
        }}
      />

      {showEmoji && (
        <View style={styles.emojiPanel}>
          {EMOJI.map((emoji, index) => (
            <Pressable key={index} onPress={() => insertEmoji(emoji)} style={styles.emojiCell}>
              <Text style={styles.emojiText}>{emoji}</Text>
            </Pressable>
          ))}
        </View>
      )}

      <View style={styles.composer}>
        <Pressable onPress={() => setShowEmoji((v) => !v)} style={styles.emojiBtn} hitSlop={8}>
          <Ionicons name="happy-outline" size={24} color={colors.textSecondary} />
        </Pressable>
        <TextInput
          value={input}
          onChangeText={setInput}
          placeholder="Type a message..."
          placeholderTextColor={colors.textMuted}
          style={styles.input}
          multiline
          maxLength={5000}
          returnKeyType="send"
          blurOnSubmit={false}
          onSubmitEditing={send}
        />
        <Pressable
          onPress={send}
          disabled={sending || !input.trim()}
          style={[styles.sendBtn, (sending || !input.trim()) && styles.sendBtnDisabled]}
        >
          <Ionicons name="send" size={18} color={colors.onAccent} />
        </Pressable>
      </View>
    </KeyboardAvoidingView>
  );
}

function Bubble({ message, own, showName }: { message: ChatMessage; own: boolean; showName: boolean }) {
  return (
    <View style={[styles.bubbleRow, own && styles.bubbleRowOwn]}>
      <View style={[styles.bubble, own ? styles.bubbleOwn : styles.bubbleOther]}>
        {showName ? (
          <Text style={styles.senderName}>{message.sender_name}</Text>
        ) : null}
        <Text style={[styles.bubbleText, own && styles.bubbleTextOwn]}>
          {message.content}
        </Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  frozen: {
    backgroundColor: colors.dangerSoft,
    borderBottomWidth: 1,
    borderBottomColor: colors.dangerBorder,
    padding: spacing.md,
  },
  frozenText: {
    color: colors.danger,
    fontFamily: fonts.sans,
    fontSize: 13,
    textAlign: 'center',
  },
  list: {
    padding: spacing.lg,
    paddingBottom: spacing.xl,
  },
  bubbleRow: {
    marginVertical: spacing.xs,
    alignItems: 'flex-start',
  },
  bubbleRowOwn: {
    alignItems: 'flex-end',
  },
  bubble: {
    maxWidth: '80%',
    borderRadius: radii.lg,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  bubbleOwn: {
    backgroundColor: colors.accent,
  },
  bubbleOther: {
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
  },
  bubbleText: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
  bubbleTextOwn: {
    color: colors.onAccent,
  },
  senderName: {
    color: colors.accent,
    fontFamily: fonts.mono,
    fontSize: 11,
    letterSpacing: 0.04,
    textTransform: 'uppercase',
    marginBottom: 2,
  },
  emojiPanel: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
    backgroundColor: colors.bgCard,
  },
  emojiCell: {
    padding: spacing.sm,
  },
  emojiText: {
    fontSize: 22,
  },
  composer: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: spacing.sm,
    padding: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
    backgroundColor: colors.bg,
  },
  emojiBtn: {
    paddingVertical: spacing.sm,
  },
  input: {
    flex: 1,
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radii.pill,
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
    paddingHorizontal: spacing.lg,
    paddingVertical: 10,
    maxHeight: 120,
  },
  sendBtn: {
    width: 38,
    height: 38,
    borderRadius: 19,
    backgroundColor: colors.accent,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sendBtnDisabled: {
    opacity: 0.4,
  },
});