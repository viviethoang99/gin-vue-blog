import 'dart:async';

import 'package:jaspr/dom.dart';
import 'package:jaspr/jaspr.dart';

enum ToastType { info, success, warning, error }

class ToastMessage {
  const ToastMessage({required this.id, required this.message, required this.type, required this.duration});

  final String id;
  final String message;
  final ToastType type;
  final Duration duration;
}

sealed class ToastEvent {
  const ToastEvent();
}

final class ToastShown extends ToastEvent {
  const ToastShown(this.toast);
  final ToastMessage toast;
}

final class ToastDismissed extends ToastEvent {
  const ToastDismissed(this.id);
  final String id;
}

class ToastController {
  ToastController._();

  static final ToastController instance = ToastController._();

  final StreamController<ToastEvent> _events = StreamController<ToastEvent>.broadcast();
  int _nextId = 0;

  Stream<ToastEvent> get events => _events.stream;

  String show(String message, {ToastType type = ToastType.info, Duration duration = const Duration(seconds: 4)}) {
    final id = 'toast-${_nextId++}';
    _events.add(ToastShown(ToastMessage(id: id, message: message, type: type, duration: duration)));
    return id;
  }

  String success(String message, {Duration duration = const Duration(seconds: 4)}) =>
      show(message, type: ToastType.success, duration: duration);

  String warning(String message, {Duration duration = const Duration(seconds: 4)}) =>
      show(message, type: ToastType.warning, duration: duration);

  String error(String message, {Duration duration = const Duration(seconds: 5)}) =>
      show(message, type: ToastType.error, duration: duration);

  void dismiss(String id) => _events.add(ToastDismissed(id));
}

@client
class ToastViewport extends StatefulComponent {
  const ToastViewport({this.maxVisible = 4, super.key});

  final int maxVisible;

  @override
  State<ToastViewport> createState() => _ToastViewportState();

  @css
  static List<StyleRule> get styles => [
    css('.toast-viewport').styles(
      display: .flex,
      position: .fixed(top: 20.px, right: 20.px),
      zIndex: ZIndex(1000),
      width: 360.px,
      maxWidth: .expression('calc(100vw - 32px)'),
      pointerEvents: .none,
      flexDirection: .column,
      gap: .all(10.px),
    ),
    css('.app-toast').styles(
      display: .flex,
      padding: .symmetric(vertical: 12.px, horizontal: 14.px),
      border: .all(style: .solid, color: const Color.variable('--color-border'), width: 1.px),
      radius: .all(.circular(12.px)),
      pointerEvents: .auto,
      animation: Animation(name: 'toast-enter', duration: 180.ms, curve: .easeOut),
      alignItems: .center,
      gap: .all(10.px),
      color: const Color.variable('--color-text-primary'),
      backgroundColor: const Color.variable('--color-surface-elevated'),
      raw: const {'box-shadow': 'var(--shadow-elevated)'},
    ),
    css('.app-toast--success').styles(
      border: .only(
        left: .solid(color: const Color.variable('--color-success'), width: 4.px),
      ),
    ),
    css('.app-toast--warning').styles(
      border: .only(
        left: .solid(color: const Color('#c78318'), width: 4.px),
      ),
    ),
    css('.app-toast--error').styles(
      border: .only(
        left: .solid(color: const Color.variable('--color-danger'), width: 4.px),
      ),
    ),
    css('.app-toast__message').styles(flex: Flex(grow: 1), fontSize: 0.9.rem, lineHeight: 1.45.em),
    css('.app-toast__close').styles(
      padding: .all(4.px),
      border: .none,
      cursor: .pointer,
      color: const Color.variable('--color-text-secondary'),
      fontSize: 1.1.rem,
      backgroundColor: Colors.transparent,
    ),
    css.keyframes('toast-enter', {
      '0%': Styles(opacity: 0, transform: .translate(y: (-8).px)),
      '100%': Styles(opacity: 1, transform: .translate(y: 0.px)),
    }),
    css.media(MediaQuery.screen(maxWidth: 600.px), [
      css('.toast-viewport').styles(
        position: .fixed(top: 12.px, right: 16.px, left: 16.px),
        width: .auto,
      ),
    ]),
    css.media(const MediaQuery.raw('(prefers-reduced-motion: reduce)'), [
      css('.app-toast').styles(raw: const {'animation': 'none'}),
    ]),
  ];
}

class _ToastViewportState extends State<ToastViewport> {
  final List<ToastMessage> _toasts = [];
  final Map<String, Timer> _timers = {};
  StreamSubscription<ToastEvent>? _subscription;

  @override
  void initState() {
    super.initState();
    _subscription = ToastController.instance.events.listen(_handleEvent);
  }

  void _handleEvent(ToastEvent event) {
    switch (event) {
      case ToastShown(:final toast):
        _show(toast);
      case ToastDismissed(:final id):
        _dismiss(id);
    }
  }

  void _show(ToastMessage toast) {
    if (!mounted) return;
    setState(() {
      _toasts.add(toast);
      while (_toasts.length > component.maxVisible) {
        final removed = _toasts.removeAt(0);
        _timers.remove(removed.id)?.cancel();
      }
    });
    if (toast.duration > Duration.zero) {
      _timers[toast.id] = Timer(toast.duration, () => _dismiss(toast.id));
    }
  }

  void _dismiss(String id) {
    _timers.remove(id)?.cancel();
    if (!mounted || !_toasts.any((toast) => toast.id == id)) return;
    setState(() => _toasts.removeWhere((toast) => toast.id == id));
  }

  @override
  void dispose() {
    _subscription?.cancel();
    for (final timer in _timers.values) {
      timer.cancel();
    }
    _timers.clear();
    super.dispose();
  }

  @override
  Component build(BuildContext context) {
    return div(
      classes: 'toast-viewport',
      attributes: const {'aria-live': 'polite', 'aria-atomic': 'false'},
      [
        for (final toast in _toasts)
          div(
            key: ValueKey(toast.id),
            classes: 'app-toast app-toast--${toast.type.name}',
            attributes: {'role': toast.type == ToastType.error ? 'alert' : 'status'},
            [
              span(classes: 'app-toast__message', [.text(toast.message)]),
              button(
                type: .button,
                classes: 'app-toast__close',
                attributes: const {'aria-label': 'Đóng thông báo'},
                onClick: () => ToastController.instance.dismiss(toast.id),
                [.text('×')],
              ),
            ],
          ),
      ],
    );
  }
}
