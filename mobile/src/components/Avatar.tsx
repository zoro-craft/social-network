import { useState } from 'react';
import { Image } from 'expo-image';
import { StyleSheet, Text, View } from 'react-native';
import { mediaSource } from '@/lib/api';
import { colors, fonts } from '@/lib/theme';

interface Props {
  src?: string | null;
  name?: string;
  size?: number;
}

export default function Avatar({ src, name, size = 28 }: Props) {
  const source = mediaSource(src);
  const uri = source?.uri ?? null;
  const radius = size / 2;
  const initial = (name || '').trim().charAt(0).toUpperCase() || '?';
  // Remember *which* url failed, so a new avatar retries instead of staying
  // stuck on the fallback after an edit.
  const [failedUri, setFailedUri] = useState<string | null>(null);
  const showFallback = !uri || failedUri === uri;

  return (
    <View style={[styles.wrap, { width: size, height: size, borderRadius: radius }]}>
      {!showFallback ? (
        <Image
          source={source}
          style={{ width: size, height: size, borderRadius: radius }}
          contentFit="cover"
          onError={() => setFailedUri(uri)}
        />
      ) : (
        <Text style={[styles.initial, { fontSize: size * 0.4 }]}>{initial}</Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    backgroundColor: colors.bgHover,
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
    borderWidth: 1,
    borderColor: colors.border,
  },
  initial: {
    color: colors.textSecondary,
    fontFamily: fonts.sansSemiBold,
  },
});