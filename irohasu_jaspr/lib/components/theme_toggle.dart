import 'dart:async';

import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';
import 'package:universal_web/web.dart' as web;

enum AppThemeMode { system, light, dark }

class ThemeController {
  ThemeController._();

  static final ThemeController instance = ThemeController._();
  static const storageKey = 'irohasu-theme';

  final StreamController<AppThemeMode> _changes = StreamController<AppThemeMode>.broadcast();
  StreamSubscription<web.Event>? _systemThemeSubscription;
  web.MediaQueryList? _darkModeQuery;
  AppThemeMode _mode = AppThemeMode.system;
  bool _initialized = false;

  Stream<AppThemeMode> get changes => _changes.stream;
  AppThemeMode get mode => _mode;

  AppThemeMode get resolvedMode {
    if (_mode != AppThemeMode.system) return _mode;
    return _darkModeQuery?.matches == true ? AppThemeMode.dark : AppThemeMode.light;
  }

  void initialize() {
    if (!kIsWeb || _initialized) return;
    _initialized = true;
    _darkModeQuery = web.window.matchMedia('(prefers-color-scheme: dark)');

    try {
      _mode = AppThemeMode.values.firstWhere(
        (mode) => mode.name == web.window.localStorage.getItem(storageKey),
        orElse: () => AppThemeMode.system,
      );
    } catch (_) {
      _mode = AppThemeMode.system;
    }

    _applyTheme();
    _systemThemeSubscription = web.EventStreamProviders.changeEvent.forTarget(_darkModeQuery).listen((_) {
      if (_mode == AppThemeMode.system) _changes.add(_mode);
    });
  }

  void setMode(AppThemeMode mode) {
    if (!kIsWeb) return;
    initialize();
    if (_mode == mode) return;
    _mode = mode;

    try {
      if (mode == AppThemeMode.system) {
        web.window.localStorage.removeItem(storageKey);
      } else {
        web.window.localStorage.setItem(storageKey, mode.name);
      }
    } catch (_) {
      // The theme still works when browser storage is unavailable.
    }

    _applyTheme();
    _changes.add(_mode);
  }

  void toggle() => setMode(resolvedMode == AppThemeMode.dark ? AppThemeMode.light : AppThemeMode.dark);

  void _applyTheme() {
    final root = web.document.documentElement;
    if (root == null) return;
    if (_mode == AppThemeMode.system) {
      root.removeAttribute('data-theme');
    } else {
      root.setAttribute('data-theme', _mode.name);
    }
  }

  void dispose() {
    _systemThemeSubscription?.cancel();
    _changes.close();
  }
}

@client
class ThemeToggle extends StatefulComponent {
  const ThemeToggle({super.key});

  @override
  State<ThemeToggle> createState() => _ThemeToggleState();

  @css
  static List<StyleRule> get styles => [
    css('.theme-toggle', [
      css('&').styles(
        display: .flex,
        width: 3.em,
        height: 3.em,
        padding: .zero,
        border: .none,
        cursor: .pointer,
        justifyContent: .center,
        alignItems: .center,
        color: Colors.white,
        fontSize: 1.rem,
        backgroundColor: Colors.transparent,
      ),
      css('&:hover').styles(backgroundColor: const Color('#0005')),
      css('&:focus-visible').styles(
        outline: Outline(color: Colors.white, style: .solid, width: OutlineWidth(2.px), offset: (-4).px),
      ),
    ]),
  ];
}

class _ThemeToggleState extends State<ThemeToggle> {
  StreamSubscription<AppThemeMode>? _subscription;
  AppThemeMode _resolvedMode = AppThemeMode.light;

  @override
  void initState() {
    super.initState();
    if (!kIsWeb) return;
    ThemeController.instance.initialize();
    _resolvedMode = ThemeController.instance.resolvedMode;
    _subscription = ThemeController.instance.changes.listen((_) {
      if (!mounted) return;
      setState(() => _resolvedMode = ThemeController.instance.resolvedMode);
    });
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }

  @override
  Component build(BuildContext context) {
    final isDark = _resolvedMode == AppThemeMode.dark;
    return button(
      type: .button,
      classes: 'theme-toggle',
      attributes: {'aria-label': isDark ? 'Chuyển sang giao diện sáng' : 'Chuyển sang giao diện tối'},
      onClick: ThemeController.instance.toggle,
      [
        span(
          classes: 'material-symbols-outlined',
          attributes: const {'aria-hidden': 'true'},
          [.text(isDark ? 'light_mode' : 'dark_mode')],
        ),
      ],
    );
  }
}
