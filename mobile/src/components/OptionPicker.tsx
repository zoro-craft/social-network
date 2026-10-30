import { FlatList, Pressable, StyleSheet, Text } from 'react-native';
import Modal from '@/components/Modal';
import { colors, fonts, radii, spacing } from '@/lib/theme';
import { useState } from 'react';

export interface PickOption {
  value: string;
  label: string;
}

interface Props {
  label: string;
  options: PickOption[];
  value: string;
  onSelect: (value: string) => void;
}

// Dropdown-style selector (replaces <select>) with a native modal list.
export default function OptionPicker({ label, options, value, onSelect }: Props) {
  const [open, setOpen] = useState(false);
  const current = options.find((o) => o.value === value);

  return (
    <>
      <Pressable style={styles.field} onPress={() => setOpen(true)}>
        <Text style={styles.label}>{label}</Text>
        <Text style={styles.value}>{current ? current.label : 'Select...'}</Text>
      </Pressable>
      <Modal visible={open} onClose={() => setOpen(false)} title={label}>
        <FlatList
          data={options}
          keyExtractor={(o) => o.value}
          style={styles.list}
          renderItem={({ item }) => {
            const active = item.value === value;
            return (
              <Pressable
                style={[styles.row, active && styles.rowActive]}
                onPress={() => {
                  onSelect(item.value);
                  setOpen(false);
                }}
              >
                <Text style={[styles.rowLabel, active && styles.rowLabelActive]}>
                  {item.label}
                </Text>
              </Pressable>
            );
          }}
        />
      </Modal>
    </>
  );
}

const styles = StyleSheet.create({
  field: {
    backgroundColor: colors.bgInput,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radii.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  label: {
    color: colors.textMuted,
    fontFamily: fonts.sansSemiBold,
    fontSize: 11,
    letterSpacing: 0.08,
    textTransform: 'uppercase',
  },
  value: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 16,
    marginTop: 2,
  },
  list: {
    maxHeight: 260,
    flexGrow: 0,
  },
  row: {
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  rowActive: {
    backgroundColor: colors.accentSoft,
    borderRadius: radii.sm,
    borderBottomWidth: 0,
    paddingHorizontal: spacing.md,
  },
  rowLabel: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
  rowLabelActive: {
    color: colors.accent,
    fontFamily: fonts.sansSemiBold,
  },
});