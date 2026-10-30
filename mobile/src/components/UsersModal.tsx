import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import Modal from '@/components/Modal';
import Avatar from '@/components/Avatar';
import { UserSearchResult } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

interface Props {
  type: 'followers' | 'following';
  users: UserSearchResult[];
  onClose: () => void;
  onOpenUser: (id: number) => void;
}

export default function UsersModal({ type, users, onClose, onOpenUser }: Props) {
  return (
    <Modal visible onClose={onClose} title={type}>
      {users.length === 0 ? (
        <Text style={styles.empty}>No {type} found.</Text>
      ) : (
        <FlatList
          data={users}
          keyExtractor={(item) => String(item.id)}
          style={styles.list}
          renderItem={({ item }) => (
            <Pressable
              style={styles.row}
              onPress={() => {
                onClose();
                onOpenUser(item.id);
              }}
            >
              <Avatar src={item.avatar} name={item.first_name} size={40} />
              <View>
                <Text style={styles.name}>
                  {item.first_name} {item.last_name}
                </Text>
                <Text style={styles.handle}>@{item.nickname || item.first_name}</Text>
              </View>
            </Pressable>
          )}
        />
      )}
    </Modal>
  );
}

const styles = StyleSheet.create({
  list: {
    maxHeight: 320,
    flexGrow: 0,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  name: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
  handle: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 13,
  },
  empty: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    textAlign: 'center',
    paddingVertical: spacing.lg,
  },
});