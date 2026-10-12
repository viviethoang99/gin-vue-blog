import 'package:jaspr/dom.dart';

const primaryColor = Color.variable('--color-primary');
const onPrimaryColor = Color.variable('--color-on-primary');
const backgroundColor = Color.variable('--color-background');
const surfaceColor = Color.variable('--color-surface');
const textPrimaryColor = Color.variable('--color-text-primary');
const textSecondaryColor = Color.variable('--color-text-secondary');
const textMutedColor = Color.variable('--color-text-muted');
const borderColor = Color.variable('--color-border');
const pageAccentColor = Color.variable('--color-page-accent');
const postHoverColor = Color.variable('--color-post-hover');
const noteBackgroundColor = Color.variable('--color-note-background');
const noteTextColor = Color.variable('--color-note-text');
const dangerColor = Color.variable('--color-danger');
const successColor = Color.variable('--color-success');

const _lightTheme = <String, String>{
  'color-scheme': 'light',
  '--color-brand': '#d8ac7e',
  '--color-primary': '#8f5b2a',
  '--color-primary-hover': '#77491f',
  '--color-on-primary': '#ffffff',
  '--color-background': '#f7f7fa',
  '--color-surface': '#ffffff',
  '--color-surface-elevated': '#ffffff',
  '--color-text-primary': '#24232a',
  '--color-text-secondary': '#55535e',
  '--color-text-muted': '#7e7b88',
  '--color-border': '#e5d8ca',
  '--color-page-accent': '#f7eee5',
  '--color-post-hover': 'rgba(216, 172, 126, 0.16)',
  '--color-note-background': 'rgba(180, 119, 20, 0.10)',
  '--color-note-text': '#9a640d',
  '--color-focus': '#8f5b2a',
  '--color-danger': '#c93737',
  '--color-success': '#237a4b',
  '--shadow-elevated': '0 14px 38px rgba(20, 33, 61, 0.14)',
};

const _darkTheme = <String, String>{
  'color-scheme': 'dark',
  '--color-brand': '#d8ac7e',
  '--color-primary': '#d8ac7e',
  '--color-primary-hover': '#f4d6ab',
  '--color-on-primary': '#24232a',
  '--color-background': '#101012',
  '--color-surface': '#18181c',
  '--color-surface-elevated': '#202027',
  '--color-text-primary': '#eeeeF2',
  '--color-text-secondary': '#bebcc5',
  '--color-text-muted': '#85838e',
  '--color-border': '#3a3028',
  '--color-page-accent': '#33271e',
  '--color-post-hover': 'rgba(216, 172, 126, 0.13)',
  '--color-note-background': 'rgba(190, 125, 24, 0.16)',
  '--color-note-text': '#f1c75b',
  '--color-focus': '#f4d6ab',
  '--color-danger': '#ff7b7b',
  '--color-success': '#65d49a',
  '--shadow-elevated': '0 14px 38px rgba(0, 0, 0, 0.34)',
};

@css
List<StyleRule> get styles => [
  css(':root').styles(raw: _lightTheme),
  css(':root[data-theme="light"]').styles(raw: _lightTheme),
  css(':root[data-theme="dark"]').styles(raw: _darkTheme),
  css.media(MediaQuery.screen(prefersColorScheme: ColorScheme.dark), [
    css(':root:not([data-theme])').styles(raw: _darkTheme),
  ]),
  css('html, body').styles(
    width: 100.percent,
    minHeight: 100.vh,
    padding: .zero,
    margin: .zero,
    color: textPrimaryColor,
    fontFamily: const .list([FontFamily('Inter'), FontFamilies.sansSerif]),
    backgroundColor: backgroundColor,
  ),
  css('body').styles(
    transition: Transition.combine([
      Transition('color', duration: 180.ms, curve: .easeOut),
      Transition('background-color', duration: 180.ms, curve: .easeOut),
    ]),
  ),
  css('h1, h2, h3, p').styles(margin: .zero),
  css('.post-hover-card').styles(
    padding: .symmetric(vertical: 12.px, horizontal: 18.px),
    margin: .symmetric(vertical: (-12).px, horizontal: (-18).px),
    radius: .all(.circular(6.px)),
    transition: Transition('background-color', duration: 180.ms, curve: .easeOut),
    backgroundColor: Colors.transparent,
  ),
  css('.post-hover-card:focus-within').styles(backgroundColor: postHoverColor),
  css.media(const MediaQuery.raw('(hover: hover)'), [
    css('.post-hover-card:hover').styles(backgroundColor: postHoverColor),
  ]),
  css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
    css('body').styles(raw: const {'transition': 'none'}),
    css('.post-hover-card').styles(raw: const {'transition': 'none'}),
  ]),
];
