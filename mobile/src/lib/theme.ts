import { TextStyle, ViewStyle } from 'react-native';

// ───────────────────────────────────────────────────────────────
// Design system — mirrors frontend/src/app/styles/variables.css
// Dark Minimal UI: violet accent, Inter + JetBrains Mono.
// Colors, radii, spacing, shadows and typography are kept in sync
// with the web tokens so both clients look identical.
// ───────────────────────────────────────────────────────────────

export const colors = {
  // ── Surfaces ──
  bg: '#0a0a0c', // --bg-page
  bgElevated: '#111114', // --bg-surface
  bgCard: '#15151a', // --bg-card
  bgInput: '#0d0d10', // --bg-input
  bgHover: '#1c1c22', // --bg-hover
  bgNavbar: 'rgba(10,10,12,0.78)', // --bg-navbar
  scrim: 'rgba(4,4,6,0.72)', // --scrim

  // ── Borders ──
  border: '#23232a', // --border
  borderStrong: '#34343f', // --border-strong
  borderHover: '#34343f', // --border-hover

  // ── Brand accent (violet) ──
  accent: '#6e6bf0', // --accent
  accentHover: '#5b58e0', // --accent-hover
  accentDark: '#524fcc', // --accent-dark
  accentStrong: '#8e8cf5', // --accent-strong
  accentSoft: 'rgba(110,107,240,0.12)', // --accent-soft
  accentBorder: 'rgba(110,107,240,0.30)', // --accent-border
  onAccent: '#ffffff', // --on-accent

  // ── Semantic colors ──
  success: '#3ecf8e', // --success
  successSoft: 'rgba(62,207,142,0.12)', // --success-soft
  successBorder: 'rgba(62,207,142,0.30)', // --success-border
  danger: '#f0605f', // --danger
  dangerSoft: 'rgba(240,96,95,0.12)', // --danger-soft
  dangerBorder: 'rgba(240,96,95,0.30)', // --danger-border
  warning: '#f0b232', // --warning
  warningSoft: 'rgba(240,178,50,0.12)', // --warning-soft
  warningBorder: 'rgba(240,178,50,0.30)', // --warning-border

  // Legacy alias — "pending" states map to the warning token.
  pending: '#f0b232',
  pendingSoft: 'rgba(240,178,50,0.12)',

  // ── Typography colors ──
  textPrimary: '#f4f4f6', // --text-primary
  textSecondary: '#a3a3ad', // --text-secondary
  textMuted: '#6b6b75', // --text-muted
  textBody: '#d6d6df', // --text-body
};

export const radii = {
  sm: 8, // --radius-sm
  md: 12, // --radius-md
  base: 12, // --radius
  lg: 18, // --radius-lg
  pill: 999, // --radius-pill
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  xxl: 32,
};

// Font weight strings are kept for backward compatibility, but with
// custom fonts loaded prefer the explicit family names below.
export const fonts = {
  regular: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,

  // Inter (sans) — mirrors --font-sans
  sans: 'Inter_400Regular',
  sansMedium: 'Inter_500Medium',
  sansSemiBold: 'Inter_600SemiBold',
  sansBold: 'Inter_700Bold',
  sansExtraBold: 'Inter_800ExtraBold',

  // JetBrains Mono — mirrors --font-mono
  mono: 'JetBrainsMono_500Medium',
  monoSemiBold: 'JetBrainsMono_600SemiBold',
};

// ── Depth tokens (approximations of the web box-shadows) ──
export const shadowSm: ViewStyle = {
  shadowColor: '#000000',
  shadowOpacity: 0.4,
  shadowRadius: 3,
  shadowOffset: { width: 0, height: 1 },
  elevation: 2,
};

export const shadow: ViewStyle = {
  shadowColor: '#000000',
  shadowOpacity: 0.4,
  shadowRadius: 12,
  shadowOffset: { width: 0, height: 6 },
  elevation: 6,
};

export const shadowLg: ViewStyle = {
  shadowColor: '#000000',
  shadowOpacity: 0.5,
  shadowRadius: 24,
  shadowOffset: { width: 0, height: 12 },
  elevation: 12,
};

export const shadowFocus: ViewStyle = {
  shadowColor: colors.accent,
  shadowOpacity: 0.4,
  shadowRadius: 3,
  shadowOffset: { width: 0, height: 0 },
  elevation: 0,
};

// Kept as the default card shadow for existing call sites.
export const cardShadow: ViewStyle = {
  shadowColor: '#000000',
  shadowOpacity: 0.4,
  shadowRadius: 12,
  shadowOffset: { width: 0, height: 6 },
  elevation: 6,
};

export const screen: ViewStyle = {
  flex: 1,
  backgroundColor: colors.bg,
};

export const card: ViewStyle = {
  backgroundColor: colors.bgCard,
  borderWidth: 1,
  borderColor: colors.border,
  borderRadius: radii.lg,
  padding: spacing.lg,
  ...cardShadow,
};

export const input: TextStyle = {
  backgroundColor: colors.bgInput,
  borderWidth: 1,
  borderColor: colors.border,
  borderRadius: radii.sm,
  color: colors.textPrimary,
  fontFamily: fonts.sans,
  fontSize: 14,
  paddingHorizontal: spacing.md,
  paddingVertical: spacing.md,
};

export const cardTitle: TextStyle = {
  color: colors.textPrimary,
  fontFamily: fonts.sansSemiBold,
  fontSize: 18,
  letterSpacing: -0.01,
  marginBottom: spacing.sm,
};

export const h1: TextStyle = {
  color: colors.textPrimary,
  fontFamily: fonts.sansBold,
  fontSize: 26,
  letterSpacing: -0.02,
};

export const h2: TextStyle = {
  color: colors.textPrimary,
  fontFamily: fonts.sansSemiBold,
  fontSize: 18,
  letterSpacing: -0.01,
};

// Small uppercase section label — mirrors the web form labels.
export const label: TextStyle = {
  color: colors.textMuted,
  fontFamily: fonts.sansSemiBold,
  fontSize: 11,
  letterSpacing: 0.08,
  textTransform: 'uppercase',
};

// Monospaced numerals — mirrors --font-mono usages (stats, timestamps).
export const mono: TextStyle = {
  color: colors.textMuted,
  fontFamily: fonts.mono,
  fontSize: 12,
};

export const primary: TextStyle = {
  color: colors.textPrimary,
  fontFamily: fonts.sans,
  fontSize: 15,
};

export const secondary: TextStyle = {
  color: colors.textSecondary,
  fontFamily: fonts.sans,
  fontSize: 14,
};

export const muted: TextStyle = {
  color: colors.textMuted,
  fontFamily: fonts.sans,
  fontSize: 13,
};

export const small: TextStyle = {
  color: colors.textMuted,
  fontFamily: fonts.sans,
  fontSize: 12,
};

export const theme = {
  colors,
  radii,
  spacing,
  fonts,
};
